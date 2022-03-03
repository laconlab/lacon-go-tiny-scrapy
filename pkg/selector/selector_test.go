package selector

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/result"
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

	sites := &Websites{}
	if err := yaml.Unmarshal([]byte(cfg), sites); err != nil {
		t.Error(err)
	}

	i := 0
	for req := range NewHttpReqChan(sites) {
		exp := &result.FullWebsiteResult{}
		exp.SetId(i)
		exp.SetWebsite("test-example-1")
		exp.SetUrl(fmt.Sprintf("example1-%d", i))
		i++

		if !reflect.DeepEqual(exp, req) {
			t.Error("Expected: ", exp, " got: ", req)
		}
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

	sites := &Websites{}
	if err := yaml.Unmarshal([]byte(cfg), sites); err != nil {
		t.Error(err)
	}

	reqs := NewHttpReqChan(sites)

	req := <-reqs
	exp := &result.FullWebsiteResult{}
	exp.SetId(0)
	exp.SetWebsite("test-example-1")
	exp.SetUrl("example1-0")

	if !reflect.DeepEqual(exp, req) {
		t.Error("Expected: ", exp, " got: ", req)
	}

	req = <-reqs
	exp.SetId(10)
	exp.SetWebsite("test-example-2")
	exp.SetUrl("example2-10")

	if !reflect.DeepEqual(exp, req) {
		t.Error("Expected: ", exp, " got: ", req)
	}

	req = <-reqs
	exp.SetId(1)
	exp.SetWebsite("test-example-1")
	exp.SetUrl("example1-1")

	if !reflect.DeepEqual(exp, req) {
		t.Error("Expected: ", exp, " got: ", req)
	}

	req = <-reqs
	exp.SetId(11)
	exp.SetWebsite("test-example-2")
	exp.SetUrl("example2-11")

	if !reflect.DeepEqual(exp, req) {
		t.Error("Expected: ", exp, " got: ", req)
	}

	req = <-reqs
	exp.SetId(2)
	exp.SetWebsite("test-example-1")
	exp.SetUrl("example1-2")

	if !reflect.DeepEqual(exp, req) {
		t.Error("Expected: ", exp, " got: ", req)
	}

	select {
	case <-time.After(time.Millisecond):
	case _, ok := <-reqs:
		if ok {
			t.Error("should be null")
		}
	}
}
