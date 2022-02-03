package crawler

import (
	"time"
)

// TODO add download time
type HttpPageResponse struct {
	Id      int    `json:"id"`
	Name    string `json:"website"`
	Url     string `json:"url"`
	Content []byte `json:"content"`
}

func (h *HttpPageResponse) GetUrl() string {
	return h.Url
}

func (h *HttpPageResponse) GetId() int {
	return h.Id
}

func (h *HttpPageResponse) GetName() string {
	return h.Name
}

func (h *HttpPageResponse) GetContent() []byte {
	return h.Content
}

func newHttpPageResponse(req HttpRequest, cnt []byte) *HttpPageResponse {
	return &HttpPageResponse{
		Id:      req.GetId(),
		Url:     req.GetUrl(),
		Name:    req.GetName(),
		Contnet: cnt,
	}
}

type HttpRequest interface {
	GetId() int
	GetUrl() string
	GetName() string
}

type HttpRequestIter interface {
	Next() interface{}
}

type crawlerWorker struct {
	nextAgent func() string
	timeout   time.Duration
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
