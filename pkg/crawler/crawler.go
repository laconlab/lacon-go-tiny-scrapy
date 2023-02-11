package crawler

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

func NewCrawler[T CrawlerResult](reqs chan T, agents *HttpAgents, cfg *CrawlerConfig) chan T {

	wg := sync.WaitGroup{}
	out := make(chan T, cfg.getBufferSize())

	for i := 0; i < cfg.getWorkerPoolSize(); i++ {
		worker := &crawlerWorker[T]{
			agents:  agents,
			timeout: cfg.getTimeout(),
			wg:      &wg,
			in:      reqs,
			out:     out,
		}

		wg.Add(1)
		go worker.start()
	}

	go func(wg *sync.WaitGroup, out chan T) {
		wg.Wait()
		close(out)
		log.Println("Stopping all crawler workers")
	}(&wg, out)

	return out
}

func (w *crawlerWorker[T]) start() {
	log.Println("Starting new crawler worker")
	defer w.wg.Done()

	for wlog := range w.in {
		for {
			if cnt, retry := w.download(wlog.GetUrl()); cnt != nil {
				wlog.SetRawContent(cnt)
				ts := time.Now()
				wlog.SetDownloadDate(ts.Format("YYYYMMDD"))
				w.out <- wlog
				break
			} else if !retry {
				break
			}
		}
	}
}

// return contnet and if request should be retried
func (w *crawlerWorker[T]) download(url string) ([]byte, bool) {
	agent := w.agents.Next()

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
