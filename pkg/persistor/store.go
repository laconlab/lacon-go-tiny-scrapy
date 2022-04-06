package persistor

import (
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"sync"
)

func NewStore[T Data](cfg *PersistorConfig, items <-chan T) {
	wg := &sync.WaitGroup{}
	id := 0
	fullPathPattern := cfg.getPath() + cfg.getNamePattern() + ".gz"
	for item := range items {
		wg.Add(1)
		go store(wg, item.GetRawWebsiteAsJSON(), fmt.Sprintf(fullPathPattern, id))
		id++
	}

	wg.Wait()
}

func store(wg *sync.WaitGroup, page []byte, path string) {
	defer wg.Done()

	f, err := os.Create(path)
	if err != nil {
		log.Println("Cannot open file", path, err)
		return
	}

	w := gzip.NewWriter(f)
	w.Write(page)
	w.Close()
	f.Close()
}
