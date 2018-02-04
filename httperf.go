package main

import (
	"crypto/tls"
	"github.com/tamura2004/httperf/config"
	"github.com/tamura2004/httperf/counter"
	"github.com/tamura2004/httperf/netstat"
	"github.com/tamura2004/httperf/slow"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type result struct {
	start    bool // true = start, false = end
	duration time.Duration
	status   string
	user     int
	ix       int
}

var p config.Config

var ch struct {
	res   chan result
	multi chan int
}

var (
	logfile *os.File
	wg      *sync.WaitGroup = &sync.WaitGroup{}
)

var client *http.Client

func main() {
	defer logfile.Close()

	ch.res = make(chan result, 2048)

	for i := 0; i < p.User; i++ {
		time.Sleep(100 * time.Millisecond)
		wg.Add(1)
		go target(i)
	}

	go netstat.Start()
	monitor()
}

func target(userID int) {
	defer wg.Done()
	for i := 0; i < p.Count; i++ {
		ch.res <- result{
			start: true,
			user:  userID,
			ix:    i,
		}

		sleep(p.Duration.Duration)

		status, duration := get(userID)

		ch.res <- result{
			start:    false,
			duration: duration,
			status:   status,
			user:     userID,
			ix:       i,
		}
	}
}

func sleep(d time.Duration) {
	time.Sleep(d * time.Duration(rand.ExpFloat64()))
}

func monitor() {
	c := counter.New()

	go func() {
		wg.Wait()
		close(ch.res)
	}()
	for r := range ch.res {
		if r.start {
			c.MultiUp()
			c.TPSUp()
			c.TPMUp()
		} else {
			log.Println(r.duration, c.Multi, r.user, r.ix, r.status)
			c.CountUp()
			c.AddDuration(r.duration)
		}
	}
	c.EachTr(func(tr, time string, tp int) {
		log.Printf("%s,%#v,%d", tr, time, tp)
	})
	log.Printf("average = %s", c.Total / time.Duration(c.Count))
}

func get(userID int) (status string, duration time.Duration) {
	req, err := http.NewRequest("GET", p.Url, nil)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	duration = time.Since(start)
	status = res.Status

	bodyHandler(res.Body, userID)

	return
}

func bodyHandler(body io.Reader, userID int) {
	if p.BPS != 0 {
		body = slow.NewReader(body, p.BPS)
	}

	var out io.Writer = ioutil.Discard
	if userID == 0 && p.Varbose {
		name := time.Now().Format("index20060102150405.html")
		file, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			panic("cannot open log logfile:" + err.Error())
		}
		out = io.MultiWriter(file, os.Stdout)
	}

	io.Copy(out, body)
}

func init() {
	rand.Seed(time.Now().UnixNano())
	initLog()
	initOption()     // use log
	initHttpClient() // use option
}

// initialize log
func initLog() {
	name := time.Now().Format("log20060102.log")
	logfile, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannot open log logfile:" + err.Error())
	}
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)
}

// initialize command line option
func initOption() {
	p = config.New("config.toml")
	log.Printf("%#v", p)
}

// initialize http client
func initHttpClient() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConnsPerHost: 2048,
	}
	if p.Proxy != "" {
		proxyURL, err := url.Parse(p.Proxy)
		if err != nil {
			log.Fatal(err)
		}
		tr.Proxy = http.ProxyURL(proxyURL)
	}
	client = &http.Client{Transport: tr}
}
