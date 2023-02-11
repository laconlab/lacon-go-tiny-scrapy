package crawler

import (
	"sync"
	"time"
)

type CrawlerResult interface {
	GetUrl() string
	SetRawContent([]byte)
	SetDownloadDate(string)
}

type SiteResponse struct {
	Id           int    `json:"id"`
	Name         string `json:"site"`
	Cnt          string `json:"content"`
	Url          string `json:"url"`
	DownloadDate string `json:"download_date"`
}

type crawlerWorker[T CrawlerResult] struct {
	agents  *HttpAgents
	timeout time.Duration
	wg      *sync.WaitGroup
	in      chan T
	out     chan T
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
