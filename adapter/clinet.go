package adapter

import (
	"encoding/base64"
	"github.com/tamura2004/httperf/domain/entity"
	"log"
	"net/http"
	"net/url"
)

// type Config struct {
//   Client   ClientConfig
//   Scinario ScinarioConfig
//   Target   TargetConfig
// }

// type ClientConfig struct {
//   Bps      int
//   Proxy    string
//   Header   map[string]string
//   UserName string
//   Password string
// }

// type ScinarioConfig struct {
//   Count     int
//   Interval  string
//   RampUp    string
//   WorkerNum int
//   Timeout   string
// }

// type TargetConfig struct {
//   Url []string
// }

type ClientConfig struct {
	Proxy  func(*http.Request) (*url.URL, error)
	Header map[string]string
	URL    string
}

func ConvertClientConfig() ClientConfig {
	c := entity.Config.Client
	t := entity.Config.Target
	c.Header["Proxy-Authorization"] = BasicAuth()

	return ClientConfig{
		Proxy:  Proxy(),
		Header: c.Header,
		URL:    t.Url[0],
	}
}

func Proxy() func(*http.Request) (*url.URL, error) {
	c := entity.Config.Client
	proxyURL, err := url.Parse(c.Proxy)
	if err != nil {
		log.Fatal(err)
	}
	return http.ProxyURL(proxyURL)
}

func BasicAuth() string {
	c := entity.Config.Client
	auth := c.UserName + ":" + c.Password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
