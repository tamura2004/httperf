package client

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
)

func New(proxy string) *http.Client {
	return &http.Client{Transport: newTR(proxy)}
}

func newTR(proxy string) *http.Transport {
	t := defaultTR()
	if proxy == "" {
		return t
	}
	return addProxy(t, proxy)
}

func defaultTR() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConnsPerHost: 2048,
	}
}

func addProxy(t *http.Transport, proxy string) *http.Transport {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		log.Fatal(err)
	}
	t.Proxy = http.ProxyURL(proxyURL)
	return t
}
