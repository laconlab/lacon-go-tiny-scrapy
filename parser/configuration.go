package parser

import (
	"lacon-go-tiny-scrapy/utils"
)

func newParserConfiguration() configuration {
	var config configuration
	utils.LoadConfig(utils.AppConfigurationFileName, &config)
	return config
}

type configuration struct {
	Parser struct {
		WorkerPoolSize int `yaml:"number-of-workers"`
	}
}

func (c *configuration) getWorkerPoolSize() int {
	return c.Parser.WorkerPoolSize
}
