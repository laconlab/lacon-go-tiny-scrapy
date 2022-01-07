package crawler

import (
	"fmt"
	"testing"
	"time"

	"github.com/laconlab/lacon-go-tiny-scrapy/src/selector"
)


func Test200Download(t *testing.T) {
    cfg := Config{
        MainWorkerCount: 1,
        Timeout: time.Second,
    }

    req := selector.HttpRequest{
        Id: 10,
        Name: "testname",
        Url: "example.com",
    }

    reqs := make(chan request, 1)
    reqs <- req
    close(reqs)

    pageResps := New(reqs, cfg)

    for resp := range pageResps {
        fmt.Println(resp)
    }

}
