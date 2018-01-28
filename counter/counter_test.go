package counter_test

import (
	"fmt"
	"github.com/tamura2004/httperf/counter"
)

func ExampleNewCounter() {
	fmt.Printf("%#v", counter.New())
	// Output:
	// counter.Counter{tr:[]map[string]int{map[string]int{}, map[string]int{}}, Count:0, Multi:0, Total:0}
}

func ExampleTPSString() {
	fmt.Println(counter.TPS)
	// Output:
	// TPS
}

func ExampleTPMString() {
	fmt.Println(counter.TPM)
	// Output:
	// TPM
}

func ExampleMultiUp() {
	c := counter.New()
	c.MultiUp()
	fmt.Printf("%#v", c.Multi)
	// Output:
	// 1
}

func ExampleMultiDown() {
	c := counter.New()
	c.MultiDown()
	fmt.Printf("%#v", c.Multi)
	// Output:
	// -1
}

func ExampleCountUp() {
	c := counter.New()
	c.CountUp()
	fmt.Printf("%#v", c.Count)
	// Output:
	// 1
}

func ExampleCountDown() {
	c := counter.New()
	c.CountDown()
	fmt.Printf("%#v", c.Count)
	// Output:
	// -1
}
