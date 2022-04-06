package main

import (
	"io/ioutil"
	"log"

	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/parser"
	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/persistor"
	"gopkg.in/yaml.v2"
)

func main() {
	cfg, err := ioutil.ReadFile("resources/application.yml")

	if err != nil {
		log.Fatal(err)
	}

	loadConfig := &persistor.PersistorConfig{}
	if err := yaml.Unmarshal(cfg, loadConfig); err != nil {
		log.Fatal(err)
	}

	rules := &parser.ParserRules{}
	if err := yaml.Unmarshal(cfg, rules); err != nil {
		log.Fatal(err)
	}

	items := make(chan persistor.Data)
	persistor.NewStore(loadConfig, items)

	parser := parser.NewParser(rules)
	for page := range persistor.NewLoader(loadConfig) {
		parser.Parse(page)
		items <- page
	}
}
