package port

import (
	"io"
)

type client interface {
	Get(url string) io.Reader
}

var Client client
