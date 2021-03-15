package data_storing

import "lacon-go-tiny-scrapy/utils"

func newDatabaseConfiguration() configuration {
	var config configuration
	utils.LoadConfig(utils.AppConfigurationFileName, &config)
	return config
}

type configuration struct {
	Database struct {
		WorkerPoolSize int `yaml:"number-of-workers"`
	}
}

func (c *configuration) getWorkerPoolSize() int {
	return c.Database.WorkerPoolSize
}
