package client_test

import (
	"fmt"
	"github.com/tamura2004/httperf/client"
)

func ExampleNew() {
	cl := client.New("http://test.proxy.net:8888/")
	fmt.Printf("%T", cl)
	// Output:
	// *http.Client
}
