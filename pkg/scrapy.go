package lacon

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

func RunScrapy(cfgPath string) error {
	cfg, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return err
	}

	websites := &Websites{}
	if err := yaml.Unmarshal(cfg, websites); err != nil {
		return err
	}

	agents := &HttpAgents{}
	if err := yaml.Unmarshal(cfg, agents); err != nil {
		return err
	}

	crawlerCfg := &CrawlerConfig{}
	if err := yaml.Unmarshal(cfg, crawlerCfg); err != nil {
		return err
	}

	presisterCfg := &PersistorConfig{}
	if err := yaml.Unmarshal(cfg, presisterCfg); err != nil {
		return err
	}

	siteReqs := make(chan *SiteRequest)
	go provideSiteRequests(websites, agents, siteReqs)

	// download stats
	stats := NewStats()
	go func() {
		for {
			time.Sleep(time.Minute)
			log.Println(stats.String())
		}
	}()

	// download
	siteResps := make(chan *SiteResponse)
	downWg := &sync.WaitGroup{}
	for i := 0; i < crawlerCfg.Config.WorkerPoolSize; i++ {
		downWg.Add(1)
		go startDownloadWorker(downWg, siteReqs, siteResps, stats)
	}

	// wait to download all and close channel
	go func() {
		log.Println("wait")
		downWg.Wait()
		close(siteResps)
		log.Println("done")
	}()

	// store
	storeWg := &sync.WaitGroup{}
	root := presisterCfg.Config.Path
	for resp := range siteResps {
		storeWg.Add(1)
		go saveResponse(storeWg, root, resp)
	}
	storeWg.Wait()

	return nil
}

// round robin around websites and create request for each
func provideSiteRequests(ws *Websites, agents *HttpAgents, out chan *SiteRequest) {
	defer close(out)

	roundRobinId := -1
	for len(ws.Sites) > 0 {
		roundRobinId = (roundRobinId + 1) % len(ws.Sites)

		if req := ws.Sites[roundRobinId].Next(); req != nil {
			req.Agent = agents.Next()
			out <- req
		} else {
			// remove webiste from the list
			ws.Sites = append(ws.Sites[:roundRobinId], ws.Sites[roundRobinId+1:]...)
		}
	}
}

func startDownloadWorker(
	wg *sync.WaitGroup,
	in chan *SiteRequest,
	out chan *SiteResponse,
	stats *DownloadStats) {
	defer wg.Done()

	for req := range in {
		if resp := downloadWithRetry(req, stats); resp != nil {
			out <- resp
		}
	}
}

func downloadWithRetry(req *SiteRequest, stats *DownloadStats) *SiteResponse {
	for {
		if cnt, done := download(req, stats); cnt != nil && done {
			return &SiteResponse{
				Id:           req.Id,
				Name:         req.Name,
				Cnt:          string(cnt),
				Url:          req.Url,
				DownloadDate: time.Now().GoString(),
			}
		} else if done {
			return nil
		} else {
			time.Sleep(time.Second)
		}
	}
}

func download(req *SiteRequest, stats *DownloadStats) ([]byte, bool) {
	start := time.Now()
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	httpReq, err := http.NewRequest(http.MethodGet, req.Url, nil)
	if err != nil {
		log.Println("Error while creating new http request", err)
		return nil, false
	}

	httpReq.Header.Set("User-Agent", req.Agent)

	resp, err := client.Do(httpReq)
	if err != nil {
		log.Println("Error while getting http response", err)
		return nil, false
	}
	defer resp.Body.Close()

	stats.Record(req.Name, resp.StatusCode, time.Since(start).Milliseconds())

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		if resp.StatusCode != 404 {
			log.Printf("Status code %d at id %d webiste %s\n", resp.StatusCode, req.Id, req.Name)
		}
		return nil, true
	}

	if resp.StatusCode >= 500 {
		log.Printf("Status code %d website %s id %d\n", resp.StatusCode, req.Name, req.Id)
		return nil, false
	}

	cnt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Recived error while reanind a website %s id %d %s\n", req.Name, req.Id, err)
		return nil, false
	}

	return cnt, true
}

func saveResponse(wg *sync.WaitGroup, root string, resp *SiteResponse) {
	defer wg.Done()

	cnt, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error: %s %d %s\n", resp.Name, resp.Id, err)
		return
	}

	path := path.Join(root, resp.Name, fmt.Sprintf("%d.gz", resp.Id))
	f, err := create(path)
	if err != nil {
		log.Printf("error: %s %d %s\n", resp.Name, resp.Id, err)
		return
	}
	defer f.Close()

	w := gzip.NewWriter(f)
	defer w.Close()
	w.Write(cnt)
}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}
