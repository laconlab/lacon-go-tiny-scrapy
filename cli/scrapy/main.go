package main

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

	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/crawler"
	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/persistor"
	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/selector"
	"gopkg.in/yaml.v2"
)

func main() {
	cfg, err := ioutil.ReadFile("resources/application.yml")

	if err != nil {
		log.Fatal(err)
	}

	websites := &selector.Websites{}
	if err := yaml.Unmarshal(cfg, websites); err != nil {
		log.Fatal(err)
	}

	agents := &crawler.HttpAgents{}
	if err := yaml.Unmarshal(cfg, agents); err != nil {
		log.Fatal(err)
	}

	crawlerCfg := &crawler.CrawlerConfig{}
	if err := yaml.Unmarshal(cfg, crawlerCfg); err != nil {
		log.Fatal(err)
	}

	presisterCfg := &persistor.PersistorConfig{}
	if err := yaml.Unmarshal(cfg, presisterCfg); err != nil {
		log.Fatal(err)
	}

	// create requests
	siteReqs := make(chan *selector.SiteRequest)
	go func() {
		defer close(siteReqs)
		roundRobinId := -1
		for len(websites.Sites) > 0 {
			roundRobinId = (roundRobinId + 1) % len(websites.Sites)
			site := websites.Sites[roundRobinId]
			req := site.Next()
			if req == nil {
				websites.Sites = append(websites.Sites[:roundRobinId], websites.Sites[roundRobinId+1:]...)
				continue
			}
			siteReqs <- req
		}
	}()

	// download requests
	siteResps := make(chan *crawler.SiteResponse)
	wg := &sync.WaitGroup{}
	wg.Add(crawlerCfg.Config.WorkerPoolSize)
	go func() {
		log.Println("wait")
		wg.Wait()
		close(siteResps)
		log.Println("done")
	}()

	for i := 0; i < crawlerCfg.Config.WorkerPoolSize; i++ {
		go func() {
			defer wg.Done()
			for req := range siteReqs {
				retry := true
				for retry {
					agent := agents.Next()
					var cnt []byte
					if cnt, retry = download(agent, req); retry {
						time.Sleep(time.Second)
					} else if cnt != nil {
						siteResps <- &crawler.SiteResponse{
							Id:           req.Id,
							Name:         req.Name,
							Cnt:          string(cnt),
							Url:          req.Url,
							DownloadDate: time.Now().GoString(),
						}
					}
				}
			}
		}()
	}

	// save requests
	wgStore := &sync.WaitGroup{}
	root := presisterCfg.Config.Path
	for resp := range siteResps {
		wgStore.Add(1)
		go func(r *crawler.SiteResponse) {
			defer wgStore.Done()

			cnt, err := json.Marshal(r)
			if err != nil {
				log.Printf("error: %s %d %s\n", r.Name, r.Id, err)
				return
			}

			path := path.Join(root, r.Name, fmt.Sprintf("%d.gz", r.Id))
			f, err := create(path)
			if err != nil {
				log.Printf("error: %s %d %s\n", r.Name, r.Id, err)
				return
			}
			defer f.Close()

			w := gzip.NewWriter(f)
			defer w.Close()
			w.Write(cnt)
		}(resp)
	}
	wgStore.Wait()
}

// return response and should retry or not
func download(agent string, req *selector.SiteRequest) ([]byte, bool) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	httpReq, err := http.NewRequest(http.MethodGet, req.Url, nil)
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
		log.Printf("Status code %d at id %d webiste %s\n", resp.StatusCode, req.Id, req.Name)
		return nil, false
	}

	if resp.StatusCode >= 500 {
		log.Printf("Status code %d website %s id %d\n", resp.StatusCode, req.Name, req.Id)
		return nil, true
	}

	cnt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Recived error while reanind a website %s id %d %s\n", req.Name, req.Id, err)
		return nil, true
	}

	return cnt, false
}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}
