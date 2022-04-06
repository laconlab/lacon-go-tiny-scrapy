package selector

import "fmt"

type WebsiteReqest interface {
	SetId(id int)
	SetWebsite(name string)
	SetUrl(url string)
}

type Websites struct {
	Sites []*Website `yaml:"websites"`
}

type Website struct {
	Name        string `yaml:"name"`
	UrlTemplate string `yaml:"urlTemplate"`
	Id          int    `yaml:"startIndex"`
	EndId       int    `yaml:"endIndex"`
}

func (w *Website) getUrl() string {
	return fmt.Sprintf(w.UrlTemplate, w.Id)
}

func (w *Website) isDone() bool {
	return w.Id > w.EndId
}

func (w *Website) inc() {
	w.Id++
}
