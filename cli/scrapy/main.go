package main

import (
	"io/ioutil"
	"log"

	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/crawler"
	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/persistor"
	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/result"
	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/selector"
	"gopkg.in/yaml.v2"
)

func main() {
	cfg, err := ioutil.ReadFile("resources/application.yml")

	if err != nil {
		log.Fatal(err)
	}

	websites := &selector.Websites{}
	if err := yaml.Unmarshal(cfg, websites); err != nil {
		log.Fatal(err)
	}

	agents := &crawler.HttpAgents{}
	if err := yaml.Unmarshal(cfg, agents); err != nil {
		log.Fatal(err)
	}

	crawlerCfg := &crawler.CrawlerConfig{}
	if err := yaml.Unmarshal(cfg, crawlerCfg); err != nil {
		log.Fatal(err)
	}

	persistorCfg := &persistor.PersistorConfig{}
	if err := yaml.Unmarshal(cfg, persistorCfg); err != nil {
		log.Fatal(err)
	}

	httpReqs := selector.NewHttpReqChan(websites)

	httpPages := crawler.NewCrawler(httpReqs, agents, crawlerCfg)

	rawPages := make(chan []byte, 10)

	go func(in chan *result.FullWebsiteResult, out chan []byte) {
		defer close(out)
		for result := range in {
			out <- result.GetRawWebsiteAsJSON()
		}
	}(httpPages, rawPages)

	persistor.NewStore(persistorCfg, rawPages)
}
