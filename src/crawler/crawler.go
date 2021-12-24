package crawler

import (
	"context"
	"lacon-go-tiny-scrapy/selector"
	"time"
	"io/ioutil"
	"net/http"
    "log"
)


type Config struct {

}


func NewCrawlerConfig() {

}


type CrawledPage struct {
    Id      int
    Name    string
    Page    []byte
}


func New(reqs chan selector.HttpRequests) <-chan CrawledPage{
    responses := make(chan CrawledPage)
    dead_letter_queue := make(chan selector.HttpRequests)

    for i := 0; i < 10; i++ {
        go download(responses, reqs, dead_letter_queue)
    }

    for i := 0; i < 5; i++ {
        go download(responses, dead_letter_queue, dead_letter_queue)
    }
    return responses
}

const retryPauseTime = 1 * time.Minute
func download(
    resps chan CrawledPage,
    reqs chan selector.HttpRequests,
    dlq chan selector.HttpRequests) {

    for req := range reqs {
        resp, retry := asyncDownload(req.Url)
        if retry {
            dlq<-req
        } else if resp != nil {
            resps<-CrawledPage{
                Id: req.id,
                Name: req.Name,
                Page: resp,
            }
        }
    }
}

const asyncTimeout = 5 * time.Second
func asyncDownload(url string) (content []byte, timeout bool) {
    callback := make(chan []byte, 1)
    ctx, _ := context.WithTimeout(context.Background(), asyncTimeout)

    go func() {
        if resp, err := http.Get(url); err == nil {
            defer resp.Body.Close()
            if resp.StatusCode >= 200 && resp.StatusCode < 300 {
                if content, err := ioutil.ReadAll(resp.Body); err != nil {
                    log.Printf("Error while reading http response %v\n", err)
                } else {
                    callback <- content
                }
            } else {
                callback <- nil
            }
        }
    }()

    select {
    case ret := <-callback:
        return ret, false
    case <- ctx.Done():
        return nil, true
    }
}

