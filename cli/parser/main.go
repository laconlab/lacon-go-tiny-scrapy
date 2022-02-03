package main

import (
	"io/ioutil"
	"log"

	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/persistor"
	"gopkg.in/yaml.v2"
)

func main() {
	cfg, err := ioutil.ReadFile("resources/application.yml")

	if err != nil {
		log.Fatal(err)
	}

	loadConfig := &persistor.LoadConfig{}
	if err := yaml.Unmarshal(cfg, loadConfig); err != nil {
		log.Fatal(err)
	}

	persistor.NewLoader(loadConfig)
}
