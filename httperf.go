package main

import (
	"crypto/tls"
	"flag"
	"github.com/tamura2004/httperf/ctr"
	"github.com/tamura2004/httperf/ns"
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

type parm struct {
	url      string
	proxy    string
	count    int
	user     int
	bps      int
	duration time.Duration
	varbose  bool
	agent    string
}

var p parm

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

	for i := 0; i < p.user; i++ {
		time.Sleep(100 * time.Millisecond)
		wg.Add(1)
		go target(i)
	}

	go ns.Log()
	monitor()
}

func target(userID int) {
	defer wg.Done()
	for i := 0; i < p.count; i++ {
		ch.res <- result{
			start: true,
			user:  userID,
			ix:    i,
		}

		sleep(p.duration)

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
	c := ctr.New()

	go func() {
		wg.Wait()
		close(ch.res)
	}()
	for r := range ch.res {
		if r.start {
			c.MultiUp()
			c.TrUp(ctr.TPS)
			c.TrUp(ctr.TPM)
		} else {
			log.Println(r.duration, c.Multi, r.user, r.ix, r.status)
			c.CountUp()
			c.AddDuration(r.duration)
		}
	}
	for _, tr := range []ctr.TransactionPerTime{ctr.TPM, ctr.TPS} {
		c.Each(tr, func(time string, tps int) {
			log.Printf("%s,%#v,%d", tr, time, tps)
		})
	}
}

func get(userID int) (status string, duration time.Duration) {
	req, err := http.NewRequest("GET", p.url, nil)
	if err != nil {
		log.Fatal(err)
	}

	if p.agent != "" {
		req.Header.Set("User-Agent", p.agent)
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
	if p.bps != 0 {
		body = slow.NewReader(body, p.bps)
	}

	var out io.Writer = ioutil.Discard
	if userID == 0 && p.varbose {
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
	flag.StringVar(&p.url, "url", "http://192.168.10.32/hello/world.txt", "url")
	flag.StringVar(&p.proxy, "proxy", "", "proxy")
	flag.IntVar(&p.count, "count", 3, "num of measure per user")
	flag.IntVar(&p.user, "user", 3, "num of user")
	flag.DurationVar(&p.duration, "duration", 3*time.Second, "average duration between measure by user")
	flag.IntVar(&p.bps, "bps", 0, "bytes par sec to read for slow reader, if bps is 0 then not use slow reader")
	flag.BoolVar(&p.varbose, "varbose", false, "display stdout and save file string read from body")
	flag.StringVar(&p.agent, "agent", "", "user agent")
	flag.Parse()

	log.Printf("%#v\n", p)
	log.Println("start")
	log.Printf("url=%s, proxy=%s, count=%d, user=%d, bps=%d, duration=%s, varbose=%v",
		p.url,
		p.proxy,
		p.count,
		p.user,
		p.bps,
		p.duration,
		p.varbose,
	)
}

// initialize http client
func initHttpClient() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConnsPerHost: 2048,
	}
	if p.proxy != "" {
		proxyURL, err := url.Parse(p.proxy)
		if err != nil {
			log.Fatal(err)
		}
		tr.Proxy = http.ProxyURL(proxyURL)
	}
	client = &http.Client{Transport: tr}
}
