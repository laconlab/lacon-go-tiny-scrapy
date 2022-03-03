package persistor

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/result"
)

func NewLoader(cfg *PersistorConfig) chan *result.FullWebsiteResult {
	wg := sync.WaitGroup{}
	output := make(chan *result.FullWebsiteResult, cfg.getBufferSize())

	files, err := ioutil.ReadDir(cfg.getPath())
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		wg.Add(1)
		go load(&wg, cfg.getPath()+file.Name(), output)
	}

	go func(wg *sync.WaitGroup, out chan *result.FullWebsiteResult) {
		wg.Wait()
		close(output)
		log.Println("Loading completed, channel closed")
	}(&wg, output)

	wg.Wait()

	return output
}

func load(wg *sync.WaitGroup, path string, out chan *result.FullWebsiteResult) {
	defer wg.Done()

	f, err := os.Open(path)
	if err != nil {
		log.Println("Error opening a file", path, err)
		return
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		log.Println("Error creating gzip reader", path, err)
		return
	}
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println("Error while trying to read whole file", path, err)
		return
	}

	page := &result.FullWebsiteResult{}
	if err = json.Unmarshal(b, page); err != nil {
		log.Println("Failed to convert to page", path, err)
		return
	}

	out <- page
}
