package crawler

import (
	"lacon-go-tiny-scrapy/utils"
	"time"
)

func newCrawlerConfiguration() configuration {
	var config configuration
	utils.LoadConfig(utils.AppConfigurationFileName, &config)
	return config
}

type configuration struct {
	Crawler struct {
		WorkerPoolSize int `yaml:"number-of-workers"`
		RetryCount int `yaml:"max-retry-attempts"`
		Timeout int `yaml:"timeout"`
		TempStopTime int `yaml:"temp-stop-time"`
	}
}

func (c *configuration) getWorkerPoolSize() int {
	return c.Crawler.WorkerPoolSize
}

func (c *configuration) getRetryCount() int {
	return c.Crawler.RetryCount
}

func (c *configuration) getTimeout() time.Duration {
	return time.Duration(c.Crawler.Timeout) * time.Second
}

func (c *configuration) getTempStopTime() time.Duration {
	return time.Duration(c.Crawler.TempStopTime) * time.Second
}
