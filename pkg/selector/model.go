package selector

import (
	"fmt"
	"log"
)

type WebsiteReqest interface {
	SetId(id int)
	SetWebsite(name string)
	SetUrl(url string)
}

type SiteRequest struct {
	Id   int
	Name string
	Url  string
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
