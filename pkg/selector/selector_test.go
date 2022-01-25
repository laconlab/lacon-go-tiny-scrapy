package selector

import (
	"fmt"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestOneStateWebsite(t *testing.T) {

	cfg := `
    websites:
        -   name: "test-example-1"
            url-template: "example1-%d"
            start-index: 0
            end-index: 5
    `

	sites := Websites{}
	if err := yaml.Unmarshal([]byte(cfg), &sites); err != nil {
		t.Error(err)
	}

	i := 0
	for req := range NewHttpReqChan(sites) {
		exp := HttpRequest{
			id:   i,
			name: "test-example-1",
			url:  fmt.Sprintf("example1-%d", i),
		}
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
        url-template: "example1-%d"
        start-index: 0
        end-index: 2

    -   name: "test-example-2"
        url-template: "example2-%d"
        start-index: 10
        end-index: 11
    `

	sites := Websites{}
	if err := yaml.Unmarshal([]byte(cfg), &sites); err != nil {
		t.Error(err)
	}

	reqs := NewHttpReqChan(sites)

	req := <-reqs
	exp := HttpRequest{
		id:   0,
		name: "test-example-1",
		url:  "example1-0",
	}

	if !reflect.DeepEqual(exp, req) {
		t.Error("Expected: ", exp, " got: ", req)
	}

	req = <-reqs
	exp = HttpRequest{
		id:   10,
		name: "test-example-2",
		url:  "example2-10",
	}

	if !reflect.DeepEqual(exp, req) {
		t.Error("Expected: ", exp, " got: ", req)
	}

	req = <-reqs
	exp = HttpRequest{
		id:   1,
		name: "test-example-1",
		url:  "example1-1",
	}

	if !reflect.DeepEqual(exp, req) {
		t.Error("Expected: ", exp, " got: ", req)
	}

	req = <-reqs
	exp = HttpRequest{
		id:   11,
		name: "test-example-2",
		url:  "example2-11",
	}

	if !reflect.DeepEqual(exp, req) {
		t.Error("Expected: ", exp, " got: ", req)
	}

	req = <-reqs
	exp = HttpRequest{
		id:   2,
		name: "test-example-1",
		url:  "example1-2",
	}

	if !reflect.DeepEqual(exp, req) {
		t.Error("Expected: ", exp, " got: ", req)
	}
}
