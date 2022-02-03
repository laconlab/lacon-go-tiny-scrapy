package persistor

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

func NewStore(cfg *StoreConfig, items <-chan interface{}) {
	wg := &sync.WaitGroup{}
	id := 0
	fullPathPattern := cfg.getSavePath() + cfg.getNamePattern() + ".gz"
	for item := range items {
		item, ok := item.(Page)
		if !ok {
			log.Println("Item does not implement Page")
			continue
		}

		wg.Add(1)
		go store(wg, item, fmt.Sprintf(fullPathPattern, id))
		id++
	}

	wg.Wait()
}

func store(wg *sync.WaitGroup, page Page, path string) {
	defer wg.Done()

	b, err := json.Marshal(page)
	if err != nil {
		log.Println("Cannot convert to json", err)
		return
	}

	f, err := os.Create(path)
	if err != nil {
		log.Println("Cannot open file", path, err)
		return
	}

	w := gzip.NewWriter(f)
	w.Write(b)
	w.Close()
	f.Close()
}
