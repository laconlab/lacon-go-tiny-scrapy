package main

/*
import (
	"lacon-go-tiny-scrapy/crawler"
	"lacon-go-tiny-scrapy/data_storing"
	"lacon-go-tiny-scrapy/logger"
	"lacon-go-tiny-scrapy/parser"
	"sync"
)
*/
//import (
//"flag"
//"fmt"
//"lacon-go-tiny-scrapy/data_store"
//	"lacon-go-tiny-scrapy/selector"
//"sync"
//)

func main() {

	/*
		var dir_out string
		flag.StringVar(&dir_out, "dir_out", "output_data", "directory to store JSON files")
		var file_name string
		flag.StringVar(&file_name, "json_name", "parsed_sites", "name of JSONs generated ({json_name}.jons)")
		flag.Parse()
		waitGroup := sync.WaitGroup{}

		httpRequests := selector.NewHTTPRequestChannel(&waitGroup, 30)
		for request := range httpRequests {
			fmt.Printf("Test%v\n", request)
		}

		users := make(chan data_store.User)
		data_store.StoreJsonLines(dir_out, file_name, &waitGroup, users)

		for i := 0; i < 10; i++ {
			user := data_store.User{Name: "Mario\ntsedt", Number: "jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj"}
			users <- user
		}
		close(users)

		waitGroup.Wait()
	*/
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
