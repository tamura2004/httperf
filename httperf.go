package main

import (
	"flag"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

// options
var url *string = flag.String("url", "https://ogisui.azurewebsites.net/", "url")
var count *int = flag.Int("count", 3, "num of measure per user")
var user *int = flag.Int("user", 3, "num of user")
var duration *int64 = flag.Int64("duration", 3, "average sec between measure")

var logfile *os.File
var wg *sync.WaitGroup = &sync.WaitGroup{}

type result struct {
	duration time.Duration
	status   string
	user     int //user number
}

func main() {
	defer logfile.Close()
	ch := make(chan result)

	go monitor(ch)
	wg.Add(*user * *count)
	for i := 0; i < *user; i++ {
		go target(i, ch)
	}
	wg.Wait()
}

func monitor(ch chan result) {
	for {
		select {
		case r := <-ch:
			log.Println(r.user, *url, r.status, r.duration)
			wg.Done()
		}
	}
}

func sleep() {
	randMillisec := int64(rand.ExpFloat64() * 1000000) //average 1sec = 1000millisec
	d := *duration                                     // sec
	waitMillisec := time.Duration(randMillisec * d)
	time.Sleep(waitMillisec)
}

func target(n int, ch chan result) {
	for i := 0; i < *count; i++ {
		sleep()

		start := time.Now()
		res, err := http.Get(*url)
		defer res.Body.Close()

		if err != nil {
			log.Fatal(err)
		}

		ch <- result{
			duration: time.Since(start),
			status:   res.Status,
			user:     n,
		}
	}
}

func init() {
	// initialize log file
	name := time.Now().Format("200601021504.log")
	logfile, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannot open log logfile:" + err.Error())
	}
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))
	log.SetFlags(log.Ldate | log.Lmicroseconds)

	// get command line option
	flag.Parse()
}
