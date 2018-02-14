package infra

import (
	"github.com/tamura2004/httperf/adapter/driver"
)

func Init() {
	config := Load()
	driver.Client = config.NewClient()
}
