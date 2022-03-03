package crawler

import (
	"time"
)

type crawlerWorker struct {
	agents  *HttpAgents
	timeout time.Duration
}

type CrawlerConfig struct {
	Config struct {
		Timeout        time.Duration `yaml:"timeout"`
		BufferSize     int           `yaml:"bufferSize"`
		WorkerPoolSize int           `yaml:"workerPoolSize"`
	} `yaml:"crawler"`
}

func (c *CrawlerConfig) getTimeout() time.Duration {
	return c.Config.Timeout
}

func (c *CrawlerConfig) getBufferSize() int {
	return c.Config.BufferSize
}

func (c *CrawlerConfig) getWorkerPoolSize() int {
	return c.Config.WorkerPoolSize
}
