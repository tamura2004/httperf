package client

import (
	"crypto/tls"
	"encoding/base64"
	"github.com/BurntSushi/toml"
	"io"
	"log"
	"net/http"
	"net/url"
)

type _config struct {
	Bps                 int
	InsecureSkipVerify  bool
	MaxIdleConnsPerHost int
	Proxy               string
	UserAgent           string
	UserName            string
	Password            string
	URL                 string
}

func Load() config {
	var c _config
	_, err := toml.DecodeFile("config.toml", &c)
	if err != nil {
		log.Fatal(err)
	}
	return config{
		InsecureSkipVerify: c.InsecureSkipVerify,
		Proxy:              Proxy(c.Proxy),
		Header:             Header(c),
		URL:                string,
	}
}

func (c *_config) Proxy() func(*http.Request) (*url.URL, error) {
	proxyURL, err := url.Parse(c.Proxy)
	if err != nil {
		log.Fatal(err)
	}
	return http.ProxyURL(proxyURL)
}

func (c *_config) Header() map[string]string {
	return map[string]string{
		"User-Agent":          c.UserAgent,
		"Proxy-Authorization": c.BasicAuth(),
	}
}

func (c *_config) BasicAuth() string {
	auth := c.Username + ":" + c.Password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
