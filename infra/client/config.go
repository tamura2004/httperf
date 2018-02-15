package client

import (
	"log"
	"net/http"
	"net/url"
)

var c 

func (c *config) New() Client {
	return Client{
		Client:  c.NewClient(),
		Request: c.NewRequest(),
	}
}

func (c *config) NewClient() http.Client {
	return http.Client{
		&http.Transport{
			TLSClientConfig: {
				InsecureSkipVerify: true,
			},
			Proxy: c.Proxy,
		},
	}
}

func (c *config) NewRequest() *http.Request {
	req, err := http.NewRequest("GET", c.URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range c.Header {
		req.Header.Set(k, v)
	}
}
