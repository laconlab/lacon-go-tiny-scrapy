package crawler

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type HttpPageResponse struct {
	request     httpRequest
	pageContnet []byte
}

type httpRequest interface {
	GetId() int
	GetName() string
	GetUrl() string
}

type crawlerWorker struct {
	reqs      <-chan httpRequest
	nextAgent func() string
	timeout   time.Duration
	output    chan HttpPageResponse
}

type CrawlerConfig struct {
	Timeout        time.Duration
	BufferSize     int
	RetryQueueSize int
	WorkerPoolSize int
}

func NewCrawler(
	reqs <-chan httpRequest,
	agents *HttpAgents,
	cfg CrawlerConfig) <-chan HttpPageResponse {

	wg := sync.WaitGroup{}
	output := make(chan HttpPageResponse, cfg.BufferSize)

	for i := 0; i < cfg.WorkerPoolSize; i++ {
		worker := &crawlerWorker{
			reqs:      reqs,
			nextAgent: agents.getIter(),
			timeout:   cfg.Timeout,
			output:    output,
		}

		wg.Add(1)
		go worker.start(&wg, reqs)
	}

	go cleanup(&wg, output)

	return output
}

func cleanup(
	wg *sync.WaitGroup,
	output chan HttpPageResponse) {

	wg.Wait()
	close(output)
}

func (w *crawlerWorker) start(wg *sync.WaitGroup, input <-chan httpRequest) {
	defer wg.Done()

	for req := range input {
		retry := true
		var cnt []byte

		for retry {
			cnt, retry = w.download(req)

			if cnt != nil {
				w.output <- HttpPageResponse{
					request:     req,
					pageContnet: cnt,
				}
			}
		}
	}
}

// return contnet and if request should be retried
func (w *crawlerWorker) download(req httpRequest) ([]byte, bool) {
	url := req.GetUrl()
	agent := w.nextAgent()

	if !headerFilter(url, agent, w.timeout) {
		log.Printf("Url %s filtered out\n", url)
		return nil, false
	}

	client := &http.Client{
		Timeout: w.timeout,
	}

	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Error while creating new http request", err)
		return nil, true
	}

	httpReq.Header.Set("User-Agent", agent)

	resp, err := client.Do(httpReq)
	if err != nil {
		log.Println("Error while getting http response", err)
		return nil, true
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		log.Printf("Status code %d at url %s\n", resp.StatusCode, url)
		return nil, false
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		log.Printf("Recived status %d at url %s\n", resp.StatusCode, url)
		return nil, true
	}

	cnt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Recived error while reanind a response of url: ", url, err)
		return nil, true
	}

	return cnt, false
}
