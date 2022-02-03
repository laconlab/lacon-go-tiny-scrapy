package crawler

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func NewCrawler(
	reqs chan interface{},
	agents *HttpAgents,
	cfg *CrawlerConfig) <-chan interface{} {

	wg := sync.WaitGroup{}
	out := make(chan interface{}, cfg.getBufferSize())

	for i := 0; i < cfg.getWorkerPoolSize(); i++ {
		worker := &crawlerWorker{
			nextAgent: agents.getIter(),
			timeout:   cfg.getTimeout(),
		}

		wg.Add(1)
		go worker.start(&wg, reqs, out)
	}

	go func(wg *sync.WaitGroup, out chan interface{}) {
		wg.Wait()
		close(out)
		log.Println("Stopping all crawler workers")
	}(&wg, out)

	return out
}

func (w *crawlerWorker) start(wg *sync.WaitGroup, input chan interface{}, out chan interface{}) {
	log.Println("Starting new crawler worker")
	defer wg.Done()

	for req := range input {
		req, ok := req.(HttpRequest)
		if !ok {
			continue
		}

		retry := true
		var cnt []byte

		for retry {
			if cnt, retry = w.download(req.GetUrl()); cnt != nil {
				out <- newHttpPageResponse(req, cnt)
				break
			}
		}
	}
}

// return contnet and if request should be retried
func (w *crawlerWorker) download(url string) ([]byte, bool) {
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
