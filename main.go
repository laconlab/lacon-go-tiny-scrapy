package main

/*

import (
	"lacon-go-tiny-scrapy/crawler"
	"lacon-go-tiny-scrapy/data_storing"
	"lacon-go-tiny-scrapy/logger"
	"lacon-go-tiny-scrapy/parser"
	"lacon-go-tiny-scrapy/selector"
	"sync"
)
*/
import (
	"lacon-go-tiny-scrapy/data_store"
	"sync"
)

/*
type User struct {
	Name   string
	Number int
}

type JsonlinesStore struct {
	waitGroup          *sync.WaitGroup
	max_file_size_byte int64
	current_index      int
}

func StoreToJsonlines(wg *sync.WaitGroup, users chan User) {
	wg.Add(1)

	worker := JsonlinesStore{
		waitGroup:          wg,
		max_file_size_byte: 30000,
		current_index:      0,
	}

	go worker.start(wg, users)
}

func (js JsonlinesStore) start(wg *sync.WaitGroup, users chan User) {
	f, _ := os.Create("out/data.txt")

	for user := range users {
		b, _ := json.Marshal(user)
		f.WriteString(string(b) + "\n")
		fi, _ := f.Stat()
		if fi.Size() >= js.max_file_size_byte {
			fmt.Printf("The file is %d bytes long\n", fi.Size())
			f.Close()
			f, _ = os.Create(fmt.Sprintf("out/data_%d.txt", js.current_index))
			js.current_index += 1
		}
	}

	f.Close()
	wg.Done()
}
*/

func main() {

	users := make(chan data_store.User)
	waitGroup := sync.WaitGroup{}
	StoreToJsonlines(&waitGroup, users)
	for i := 0; i < 100000; i++ {
		user := User{Name: "Mario\ntsedt", Number: 30}
		users <- user
	}
	close(users)

	waitGroup.Wait()
	/*
		logger.INFO.Println("Starting application...")

		httpRequests := selector.NewHTTPRequestProducer(&waitGroup)
		htmlPages := crawler.NewHTMLProducer(&waitGroup, httpRequests)
		parsedFiles := parser.NewParser(&waitGroup, htmlPages)
		data_storing.NewDatabaseInsertWorker(&waitGroup, parsedFiles)


	*/
	//waitGroup.Wait()
	//logger.INFO.Println("Stopping application...")
}
