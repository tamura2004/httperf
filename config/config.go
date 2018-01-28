package config

import (
	"github.com/BurntSushi/toml"
	"log"
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
	Agent    string
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
	if _, err := toml.DecodeFile(file, &config); err != nil {
		log.Fatal(err)
	}
	return
}
