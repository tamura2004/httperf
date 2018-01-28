package config

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Config struct {
	Url      string
	Proxy    string
	Count    int
	User     int
	BPS      int
	Duration duration
	Varbose  bool
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func New(file string) (config Config) {
	if _, err := os.Stat(file); err != nil {
		content := []byte(
			`Url = "http://www.google.co.jp/"
proxy = "http://127.0.0.1:8888/"
Count = 3
User = 3
Duration = "3s"
Bps = 128000
Varbose = true`)
		ioutil.WriteFile(file, content, os.ModePerm)

	}

	if _, err := toml.DecodeFile(file, &config); err != nil {
		log.Fatal(err)
	}
	return
}
