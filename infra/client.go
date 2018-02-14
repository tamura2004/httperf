package infra

import (
	"io"
	"log"
	"net/http"
)

type Client struct {
	http.Client
	*http.Request
}

func (c *Client) Get() io.Reader {
	res, err := c.Do(c.Request)
	if err != nil {
		log.Fatal(err)
	}
	return res.Body
}
