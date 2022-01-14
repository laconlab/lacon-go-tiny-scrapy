package crawler

import "time"

type request interface {
	GetId() int
	GetName() string
	GetUrl() string
}

type Config struct {
	MainWorkerCount int
	Timeout         time.Duration
}

type PageResp struct {
	PageId  request
	Content []byte
}

type status struct {
	content []byte
	retry   bool
}
