package infra

import (
	"log"
	"net/http"
	"net/url"
)

type config struct {
	InsecureSkipVerify bool
	Proxy              func(*http.Request) (*url.URL, error)
	Header             map[string]string
	URL                string
}

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
				InsecureSkipVerify: c.InsecureSkipVerify,
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
