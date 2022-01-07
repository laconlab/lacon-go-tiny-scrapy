package selector

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

type State struct {
	Name        string `yaml:"name"`
	UrlTemplate string `yaml:"url-template"`
	StartId     int    `yaml:"start-index"`
	EndId       int    `yaml:"end-index"`
}

func (s *State) isDone() bool {
	return s.StartId > s.EndId
}

func (s *State) next() HttpRequest {
	req := HttpRequest{
		Name: s.Name,
		Id:   s.StartId,
		Url:  fmt.Sprintf(s.UrlTemplate, s.StartId),
	}

	s.StartId++
	return req
}

type States struct {
	States []*State `yaml:"websites"`
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

func NewSelector(configFile string) chan HttpRequest {
	ret := make(chan HttpRequest, 1)

	ss := States{}
	if err := yaml.Unmarshal([]byte(configFile), &ss); err != nil {
		log.Fatal(err)
	}

	go worker(&ss, ret)

	return ret
}

func worker(ss *States, ret chan HttpRequest) {
	for !roundRobin(ss, ret) {
	}
	close(ret)
}

func roundRobin(ss *States, ret chan HttpRequest) bool {
	done := true
	for _, state := range ss.States {
		if state.isDone() {
			continue
		}
		done = false
		ret <- state.next()
	}
	return done
}
