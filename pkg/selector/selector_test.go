package selector_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/result"
	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/selector"
	"gopkg.in/yaml.v2"
)

func TestOneStateWebsite(t *testing.T) {

	cfg := `
    websites:
        -   name: "test-example-1"
            urlTemplate: "example1-%d"
            startIndex: 0
            endIndex: 5
    `

	sites := &selector.Websites{}
	if err := yaml.Unmarshal([]byte(cfg), sites); err != nil {
		t.Error(err)
	}

	reqs := result.NewFullWebsiteResultChan()
	i := 0
	for req := range selector.NewHttpReqChan(sites, reqs) {
		if req.Id != i {
			t.Errorf("Expected id %v got %v\n", i, req.Id)
		}
		if req.GetWebsite() != "test-example-1" {
			t.Errorf("Expected test-example-1 got %v\n", req.GetWebsite())
		}
		if req.GetUrl() != fmt.Sprintf("example1-%d", i) {
			t.Errorf("Expected %v got %v\n", fmt.Sprintf("example1-%d", i), req.GetUrl())
		}
		i++

	}
}

func TestRoundRobin(t *testing.T) {
	cfg := `
    websites:
    -   name: "test-example-1"
        urlTemplate: "example1-%d"
        startIndex: 0
        endIndex: 2

    -   name: "test-example-2"
        urlTemplate: "example2-%d"
        startIndex: 10
        endIndex: 11
    `

	sites := &selector.Websites{}
	if err := yaml.Unmarshal([]byte(cfg), sites); err != nil {
		t.Error(err)
	}

	eReqs := result.NewFullWebsiteResultChan()
	reqs := selector.NewHttpReqChan(sites, eReqs)

	req := <-reqs
	if req.Id != 0 {
		t.Errorf("Expected 0 got %v\n", req.Id)
	}
	if req.GetWebsite() != "test-example-1" {
		t.Errorf("Expected test-example-1 got %v\n", req.GetWebsite())
	}
	if req.GetUrl() != "example1-0" {
		t.Errorf("Expected example1-0 got %v\n", req.GetUrl())
	}

	req = <-reqs
	if req.Id != 10 {
		t.Errorf("Expected 10 got %v\n", req.Id)
	}
	if req.GetWebsite() != "test-example-2" {
		t.Errorf("Expected test-example-2 got %v\n", req.GetWebsite())
	}
	if req.GetUrl() != "example2-10" {
		t.Errorf("Expected example2-10 got %v\n", req.GetUrl())
	}

	req = <-reqs
	if req.Id != 1 {
		t.Errorf("Expected 1 got %v\n", req.Id)
	}
	if req.GetWebsite() != "test-example-1" {
		t.Errorf("Expected test-example-1 got %v\n", req.GetWebsite())
	}
	if req.GetUrl() != "example1-1" {
		t.Errorf("Expected example1-1 got %v\n", req.GetUrl())
	}

	req = <-reqs
	if req.Id != 11 {
		t.Errorf("Expected 11 got %v\n", req.Id)
	}
	if req.GetWebsite() != "test-example-2" {
		t.Errorf("Expected test-example-2 got %v\n", req.GetWebsite())
	}
	if req.GetUrl() != "example2-11" {
		t.Errorf("Expected example2-11 got %v\n", req.GetUrl())
	}

	req = <-reqs
	if req.Id != 2 {
		t.Errorf("Expected 2 got %v\n", req.Id)
	}
	if req.GetWebsite() != "test-example-1" {
		t.Errorf("Expected test-example-1 got %v\n", req.GetWebsite())
	}
	if req.GetUrl() != "example1-2" {
		t.Errorf("Expected example1-2 got %v\n", req.GetUrl())
	}

	select {
	case <-time.After(time.Millisecond):
	case _, ok := <-reqs:
		if ok {
			t.Error("should be null")
		}
	}
}
