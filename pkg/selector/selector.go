package selector

import "fmt"

type Websites struct {
	Sites []*Website `yaml:"websites"`
}

type Website struct {
	Name        string `yaml:"name"`
	UrlTemplate string `yaml:"url-template"`
	StartId     int    `yaml:"start-index"`
	EndId       int    `yaml:"end-index"`
}

type HttpRequest struct {
	id   int
	name string
	url  string
}

func (w *Website) getId() int {
	return w.StartId
}

func (w *Website) getName() string {
	return w.Name
}

func (w *Website) getUrl() string {
	return fmt.Sprintf(w.UrlTemplate, w.getId())
}

func (w *Website) isDone() bool {
	return w.StartId > w.EndId
}

func (w *Website) inc() {
	w.StartId++
}

func (r HttpRequest) GetId() int {
	return r.id
}

func (r HttpRequest) GetName() string {
	return r.name
}

func (r HttpRequest) GetUrl() string {
	return r.url
}

func NewHttpReqChan(websites Websites) <-chan *HttpRequest {
	ch := make(chan *HttpRequest, 1)

	go func(ch chan<- *HttpRequest, sites []*Website) {
		defer close(ch)

		done := 0
		for done != len(sites) {
			for _, site := range sites {
				if site.isDone() {
					done++
					continue
				}

				ch <- &HttpRequest{
					id:   site.getId(),
					name: site.getName(),
					url:  site.getUrl(),
				}

				site.inc()
			}
		}
	}(ch, websites.Sites)

	return ch
}
