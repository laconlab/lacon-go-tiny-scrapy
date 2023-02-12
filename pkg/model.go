package lacon

import (
	"fmt"
	"log"
)

type SiteRequest struct {
	Id    int
	Name  string
	Url   string
	Agent string
}

type Websites struct {
	Sites []*Website `yaml:"websites"`
}

type Website struct {
	Name        string `yaml:"name"`
	UrlTemplate string `yaml:"urlTemplate"`
	Id          int    `yaml:"startIndex"`
	EndId       int    `yaml:"endIndex"`
	realStart   int
}

func (w *Website) Next() *SiteRequest {
	if w.realStart == 0 {
		w.realStart = w.Id
	}
	if w.Id >= w.EndId {
		log.Printf("done with %s\n", w.Name)
		return nil
	}

	req := &SiteRequest{
		Id:   w.Id,
		Name: w.Name,
		Url:  fmt.Sprintf(w.UrlTemplate, w.Id),
	}
	w.Id++
	if w.Id%1_000 == 0 {
		p := float64(w.Id-w.realStart) / float64(w.EndId-w.realStart)
		log.Printf("%s progress %.3f\n", w.Name, p)
	}
	return req
}

type SiteResponse struct {
	Id           int    `json:"id"`
	Name         string `json:"site"`
	Cnt          string `json:"content"`
	Url          string `json:"url"`
	DownloadDate string `json:"download_date"`
}

type CrawlerConfig struct {
	Config struct {
		WorkerPoolSize int `yaml:"workerPoolSize"`
	} `yaml:"crawler"`
}

type PersistorConfig struct {
	Config struct {
		Path       string `yaml:"savePath"`
		BufferSize int    `yaml:"bufferSize"`
	} `yaml:"persistor"`
}

type HttpAgents struct {
	Agents []string `yaml:"userAgents"`
	id     int
}

func (a *HttpAgents) Next() string {
	if len(a.Agents) == 0 {
		return ""
	}
	a.id = (a.id + 1) % len(a.Agents)
	return a.Agents[a.id]
}
