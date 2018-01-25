package ctr

import (
	"sort"
	"time"
)

type Counter struct {
	tr    []map[string]int // TPS,TPM
	Count int
	Multi int // number of concurrent user
	Total time.Duration
}

type TransactionPerTime int

const (
	TPS TransactionPerTime = iota
	TPM
)

func (t TransactionPerTime) String() string {
	if t == TPS {
		return "TPS"
	}
	return "TPM"
}

func New() Counter {
	return Counter{
		tr: []map[string]int{
			make(map[string]int),
			make(map[string]int),
		},
	}
}

func (c *Counter) MultiUp() {
	c.Multi++
}

func (c *Counter) MultiDown() {
	c.Multi--
}

func (c *Counter) CountUp() {
	c.Count++
}

func (c *Counter) CountDown() {
	c.Count--
}

func (c *Counter) TrUp(tr TransactionPerTime) {
	key := time.Now().Format("2006/01/02 15:04:05")
	if tr == TPM {
		key = time.Now().Format("2006/01/02 15:04")
	}
	c.tr[tr][key]++
}

func (c *Counter) AddDuration(d time.Duration) {
	c.Total += d
}

func (c *Counter) Each(tr TransactionPerTime, do func(string, int)) {
	for _, k := range c.Keys(tr) {
		do(k, c.tr[tr][k])
	}
}

func (c *Counter) Keys(tr TransactionPerTime) []string {
	keys := []string{}
	for k := range c.tr[tr] {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}
