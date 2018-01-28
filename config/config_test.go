package config_test

import (
	"fmt"
	"github.com/tamura2004/httperf/config"
)

func ExampleNewConfig() {
	c := config.New("test.config")
	fmt.Printf("%#v\n", c)
	// Output:
	// config.Config{Url:"http://192.168.10.32/hello/world.txt", Proxy:"http://127.0.0.1:8888/", Count:3, User:3, BPS:0, Duration:config.duration{Duration:3000000}, Varbose:false, Agent:"go-client/1.1"}
}
