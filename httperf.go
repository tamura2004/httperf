package main

import (
	"crypto/tls"
	"flag"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"sync"
	"time"
)

type result struct {
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
	duration time.Duration
	wait     time.Duration
}

var p parm

type counter struct {
	tps   map[string]int
	tpm   map[string]int
	count int
	total time.Duration
}

var c counter

var (
	logfile *os.File
	wg      *sync.WaitGroup = &sync.WaitGroup{}
)

var client *http.Client

func main() {
	defer logfile.Close()

	ch := make(chan result)
	go monitor(ch)

	wg.Add(p.user)

	for i := 0; i < p.user; i++ {
		go target(i, ch)
	}
	wg.Wait()

	logMap(c.tps, "tps")
	logMap(c.tpm, "tpm")

	log.Printf("average,%s", time.Duration(c.total/time.Duration(c.count)))
	log.Println("stop")
}

func logMap(m map[string]int, label string) {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, k := range keys {
		log.Printf("%s,%#v,%d", label, k, m[k])
	}
}

func monitor(ch chan result) {
	for {
		select {
		case r := <-ch:
			log.Println(r.duration, r.user, r.ix, r.status)
			key := time.Now().Format("2006/01/02 15:04:05")
			c.tps[key]++
			key = time.Now().Format("2006/01/02 15:04")
			c.tpm[key]++
			c.count++
			c.total += r.duration
		}
	}
}

func sleep(d time.Duration) {
	time.Sleep(d * time.Duration(rand.ExpFloat64()))
}

func target(n int, ch chan result) {
	defer wg.Done()
	for i := 0; i < p.count; i++ {
		sleep(p.duration)

		start := time.Now()

		res := get()

		ch <- result{
			duration: time.Since(start),
			status:   res.Status,
			user:     n,
			ix:       i,
		}
	}
}

func get() *http.Response {
	res, err := client.Get(p.url)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	sleep(p.wait)
	return res
}

func init() {
	// initialize log file
	name := time.Now().Format("20060102.log")
	logfile, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannot open log logfile:" + err.Error())
	}
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)

	c.tps = make(map[string]int)
	c.tpm = make(map[string]int)

	// get command line option
	flag.StringVar(&p.url, "url", "https://www.google.com/teapot", "url")
	flag.StringVar(&p.proxy, "proxy", "", "proxy")
	flag.IntVar(&p.count, "count", 3, "num of measure per user")
	flag.IntVar(&p.user, "user", 3, "num of user")
	flag.DurationVar(&p.duration, "duration", 3*time.Second, "average duration between measure by user")
	flag.DurationVar(&p.wait, "wait", 3*time.Second, "average wait between http session start and close")
	flag.Parse()

	log.Println("start")
	log.Printf("url=%s, proxy=%s, count=%d, user=%d, duration=%s", p.url, p.proxy, p.count, p.user, p.duration)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConnsPerHost: 128,
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
