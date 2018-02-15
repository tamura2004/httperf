package infra

import (
	"github.com/BurntSushi/toml"
	"github.com/tamura2004/httperf/domain/entity"
	"log"
)

func Load(filename string) {
	_, err := toml.DecodeFile(filename, &entity.Config)
	if err != nil {
		log.Fatal(err)
	}
}
