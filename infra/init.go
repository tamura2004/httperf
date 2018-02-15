package infra

import (
	"github.com/tamura2004/httperf/adapter/driver"
	"github.com/tamura2004/httperf/infra/client"
)

func Init() {
	Load("config.toml")
	adapter.Config = client.New()
}
