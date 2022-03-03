package selector

import "fmt"

type Websites struct {
	Sites []*Website `yaml:"websites"`
}

type Website struct {
	Name        string `yaml:"name"`
	UrlTemplate string `yaml:"urlTemplate"`
	StartId     int    `yaml:"startIndex"`
	EndId       int    `yaml:"endIndex"`
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
