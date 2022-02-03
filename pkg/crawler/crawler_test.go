package crawler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type Tmp struct {
	url string
}

func (t *Tmp) GetId() int {
	return 0
}

func (t *Tmp) GetName() string {
	return ""
}

func (t *Tmp) GetUrl() string {
	return t.url
}

func Test200Download(t *testing.T) {
	cfg := &CrawlerConfig{}
	cfg.Config.Timeout = time.Second
	cfg.Config.BufferSize = 10
	cfg.Config.WorkerPoolSize = 1

	agents := &HttpAgents{
		Agents: []string{"test"},
	}

	expected := "test response"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, expected)
	}))
	defer ts.Close()

	reqs := make(chan interface{}, 1)
	reqs <- &Tmp{url: ts.URL}
	close(reqs)
	crw := NewCrawler(reqs, agents, cfg)

	resp := (<-crw).(HttpPageResponse)

	actual := string(resp.GetContnet())
	if actual != expected {
		t.Errorf("Expected %s actual %s\n", expected, actual)
	}

}

func Test400Download(t *testing.T) {
	cfg := &CrawlerConfig{}
	cfg.Config.Timeout = time.Second
	cfg.Config.BufferSize = 10
	cfg.Config.WorkerPoolSize = 1

	agents := &HttpAgents{
		Agents: []string{"test"},
	}

	notExpected := "test response"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, notExpected)
	}))
	defer ts.Close()

	reqs := make(chan interface{}, 1)
	reqs <- &Tmp{url: ts.URL}
	close(reqs)
	crw := NewCrawler(reqs, agents, cfg)

	resp := <-crw

	if resp != nil {
		t.Error("Got unexpcted result")
	}
}

func TestErrorOnFilter(t *testing.T) {
	cfg := &CrawlerConfig{}
	cfg.Config.Timeout = time.Second
	cfg.Config.BufferSize = 10
	cfg.Config.WorkerPoolSize = 1

	agents := &HttpAgents{
		Agents: []string{"test"},
	}

	expected := "test response"

	callCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if callCount == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			callCount++
			return
		}
		callCount++

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, expected)
	}))
	defer ts.Close()

	reqs := make(chan interface{}, 1)
	reqs <- &Tmp{url: ts.URL}
	close(reqs)
	crw := NewCrawler(reqs, agents, cfg)

	resp := (<-crw).(HttpPageResponse)

	if callCount != 2 {
		t.Errorf("Expected to download page 2 times, actual %d times\n", callCount)
	}

	actual := string(resp.GetContnet())
	if actual != expected {
		t.Error("Got unexpcted result")
	}
}

func TestRetry(t *testing.T) {
	cfg := &CrawlerConfig{}
	cfg.Config.Timeout = time.Second
	cfg.Config.BufferSize = 10
	cfg.Config.WorkerPoolSize = 1

	agents := &HttpAgents{
		Agents: []string{"test"},
	}

	expected := "test response"

	callCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if callCount < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			callCount++
			return
		}

		callCount++
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, expected)
	}))
	defer ts.Close()

	reqs := make(chan interface{}, 1)
	reqs <- &Tmp{url: ts.URL}
	close(reqs)
	crw := NewCrawler(reqs, agents, cfg)

	resp := (<-crw).(HttpPageResponse)

	// filter call -> actual call -> retry filter call -> retry actual call
	if callCount != 4 {
		t.Errorf("Expected to download page 2 times, actual %d times\n", callCount)
	}

	actual := string(resp.GetContnet())
	if actual != expected {
		t.Error("Got unexpcted result")
	}

}
