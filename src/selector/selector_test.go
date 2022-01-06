package selector

import (
	"testing"
    "fmt"
)

func TestOneStateWebsite(t *testing.T) {
	config := `
    websites:
        -   name: "test-example-1"
            url-template: "example1-%d"
            start-index: 0
            end-index: 5
    `
    reqs := NewSelector(config)

    for i := 0; i < 6; i++ {
        req := <-reqs
        if req.Id != i || req.Name != "test-example-1" || req.Url != fmt.Sprintf("example1-%d", i) {
            t.Error("Failed")
        }

    }
}

func TestRoundRobin(t *testing.T) {
	config := `
    websites:
        -   name: "test-example-1"
            url-template: "example1-%d"
            start-index: 0
            end-index: 2

        -   name: "test-example-2"
            url-template: "example2-%d"
            start-index: 10
            end-index: 12
    `
    reqs := NewSelector(config)

    req := <-reqs
    if req.Id != 0 || req.Url != "example1-0" || req.Name != "test-example-1" {
        t.Error("Failed")
    }

    req = <-reqs
    if req.Id != 10 || req.Url != "example2-10" || req.Name != "test-example-2" {
        t.Error("Failed")
    }

    req = <-reqs
    if req.Id != 1 || req.Url != "example1-1" || req.Name != "test-example-1" {
        t.Error("Failed")
    }

    req = <-reqs
    if req.Id != 11 || req.Url != "example2-11" || req.Name != "test-example-2" {
        t.Error("Failed")
    }

    req = <-reqs
    if req.Id != 2 || req.Url != "example1-2" || req.Name != "test-example-1" {
        t.Error("Failed")
    }

    req = <-reqs
    if req.Id != 12 || req.Url != "example2-12" || req.Name != "test-example-2" {
        t.Error("Failed")
    }
}
