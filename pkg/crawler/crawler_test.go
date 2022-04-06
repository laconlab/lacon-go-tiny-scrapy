package crawler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/result"
)

func Test200Download(t *testing.T) {
	cfg := &CrawlerConfig{}
	cfg.Config.Timeout = time.Second
	cfg.Config.BufferSize = 10
	cfg.Config.WorkerPoolSize = 1

	agents := &HttpAgents{
		Agents: []string{"test"},
	}

	expected := "test response"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, expected)
	}))
	defer ts.Close()

	logs := make(chan *result.FullWebsiteResult, 1)
	logs <- &result.FullWebsiteResult{
		Id:            0,
		Website:       "",
		Url:           ts.URL,
		DownloadDate:  "",
		RawContent:    []byte{},
		ParsedContent: map[string]string{},
	}
	close(logs)
	crw := NewCrawler(logs, agents, cfg)

	resp := <-crw

	actual := string(resp.GetRawContent())
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

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, notExpected)
	}))
	defer ts.Close()

	logs := make(chan *result.FullWebsiteResult, 1)
	logs <- &result.FullWebsiteResult{Url: ts.URL}
	close(logs)
	crw := NewCrawler(logs, agents, cfg)

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
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if callCount == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			callCount++
			return
		}
		callCount++

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, expected)
	}))
	defer ts.Close()

	logs := make(chan *result.FullWebsiteResult, 1)
	logs <- &result.FullWebsiteResult{Url: ts.URL}
	close(logs)
	crw := NewCrawler(logs, agents, cfg)

	resp := <-crw

	if callCount != 2 {
		t.Errorf("Expected to download page 2 times, actual %d times\n", callCount)
	}

	actual := string(resp.GetRawContent())
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
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if callCount < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			callCount++
			return
		}

		callCount++
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, expected)
	}))
	defer ts.Close()

	logs := make(chan *result.FullWebsiteResult, 1)
	logs <- &result.FullWebsiteResult{Url: ts.URL}
	close(logs)
	crw := NewCrawler(logs, agents, cfg)

	resp := <-crw
	// filter call -> actual call -> retry filter call -> retry actual call
	if callCount != 4 {
		t.Errorf("Expected to download page 2 times, actual %d times\n", callCount)
	}

	actual := string(resp.GetRawContent())
	if actual != expected {
		t.Error("Got unexpcted result")
	}

}
