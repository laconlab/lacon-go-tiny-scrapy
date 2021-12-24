package crawler

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"
    "github.com/laconlab/lacon-go-tiny-scrapy/src/selector"
)


type Config struct {
    MainWorkerCount     int
    RetryWorkerCount    int
    DownloadTimeout     time.Duration
}

type CrawledPage struct {
    Page        selector.HttpRequest
    Content     []byte
}

type downloadStatus struct {
    content []byte
    retry   bool
}

func New(reqs chan selector.HttpRequest, cfg Config) <-chan CrawledPage{
    responses := make(chan CrawledPage)
    dead_letter_queue := make(chan selector.HttpRequest)

    for i := 0; i < cfg.MainWorkerCount; i++ {
        go worker(cfg, responses, reqs, dead_letter_queue)
    }

    for i := 0; i < cfg.RetryWorkerCount; i++ {
        go worker(cfg, responses, dead_letter_queue, dead_letter_queue)
    }

    return responses
}

func worker(
    cfg Config,
    resps chan CrawledPage,
    reqs chan selector.HttpRequest,
    dlq chan selector.HttpRequest) {

    for req := range reqs {
        if resp := download(req.Url, cfg.DownloadTimeout); resp.retry {
            dlq<-req
        } else if resp.content != nil {
            resps<-CrawledPage{
                Page:       req,
                Content:    resp.content,
            }
        }
    }
}

func download(url string, timeout time.Duration) downloadStatus {
    callback := make(chan downloadStatus, 1)
    ctx, _ := context.WithTimeout(context.Background(), timeout)

    go asyncDownload(url, callback, ctx)

    select {
    case resp := <-callback:
        return resp
    case <- ctx.Done():
        log.Printf("Timeout url %s\n", url)
        return downloadStatus{retry: true}
    }
}

func asyncDownload(url string, clb chan downloadStatus, ctx context.Context) {
    resp, err := http.Get(url)

    if err != nil {
        log.Printf("Failed to download: ", err)
        clb<-downloadStatus{retry: true}
        return
    }

    defer resp.Body.Close()
    clb<-readResponse(resp)
}

func readResponse(resp *http.Response) downloadStatus {
    if resp.StatusCode >= 500 {
        return downloadStatus{retry: true}
    } else if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
        return downloadStatus{retry: false}
    }

	cnt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while reading http response %v\n", err)
        return downloadStatus{retry: true}
	}

    return downloadStatus{content: cnt}
}

