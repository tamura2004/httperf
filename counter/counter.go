package counter

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

func (c *Counter) TPSUp() {
	key := time.Now().Format("2006/01/02 15:04:05")
	c.tr[TPS][key]++
}

func (c *Counter) TPMUp() {
	key := time.Now().Format("2006/01/02 15:04")
	c.tr[TPM][key]++
}

func (c *Counter) AddDuration(d time.Duration) {
	c.Total += d
}

func (c *Counter) EachTr(do func(string, string, int)) {
	for _, tr := range []TransactionPerTime{TPM, TPS} {
		for _, k := range c.keys(tr) {
			timestamp := k
			numOfTr := c.tr[tr][k]
			do(tr.String(), timestamp, numOfTr)
		}
	}
}

func (c *Counter) keys(tr TransactionPerTime) []string {
	keys := []string{}
	for k := range c.tr[tr] {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}
