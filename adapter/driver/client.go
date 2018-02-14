package driver

import (
	"io"
)

type client interface {
	Get() io.Reader
}

var Client client
