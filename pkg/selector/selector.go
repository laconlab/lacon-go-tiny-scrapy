package selector

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"gopkg.in/yaml.v2"
)

type States struct {
	States []*State `yaml:"websites"`
	mu     sync.Mutex
	id     int
}

func New(filePath string) *States {
	cfg, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	ss := &States{}
	if err := yaml.Unmarshal(cfg, &ss); err != nil {
		log.Fatal(err)
	}

	return ss
}

func (s *States) Next() *HttpRequest {
    s.mu.Lock()
    defer s.mu.Unlock()

	doneCount := 0
	for doneCount < len(s.States) {
		req := s.States[s.id].next()
		s.id = (s.id + 1) % len(s.States)
		if req == nil {
			doneCount++
		} else {
			return req
		}
	}

	return nil
}

type State struct {
	Name        string `yaml:"name"`
	UrlTemplate string `yaml:"url-template"`
	StartId     int `yaml:"start-index"`
	EndId       int `yaml:"end-index"`
}

func (s *State) isDone() bool {
	return s.StartId > s.EndId
}

func (s *State) next() *HttpRequest {
	if s.isDone() {
		return nil
	}

	id := s.StartId
	s.StartId++
	req := &HttpRequest{
		Id:   id,
		Name: s.Name,
		Url:  fmt.Sprintf(s.UrlTemplate, id),
	}
	return req
}

type HttpRequest struct {
	Id   int
	Name string
	Url  string
}

func (r HttpRequest) GetId() int {
	return r.Id
}

func (r HttpRequest) GetName() string {
	return r.Name
}

func (r HttpRequest) GetUrl() string {
	return r.Url
}
