package main

import (
	"lacon-go-tiny-scrapy/crawler"
	"lacon-go-tiny-scrapy/data_storing"
	"lacon-go-tiny-scrapy/logger"
	"lacon-go-tiny-scrapy/parser"
	"lacon-go-tiny-scrapy/selector"
	"sync"
)

func main() {

	logger.INFO.Println("Starting application...")
	waitGroup := sync.WaitGroup{}
	httpRequests := selector.NewHTTPRequestProducer(&waitGroup)
	htmlPages := crawler.NewHTMLProducer(&waitGroup, httpRequests)
	parsedFiles := parser.NewParser(&waitGroup, htmlPages)
	data_storing.NewDatabaseInsertWorker(&waitGroup, parsedFiles)

	waitGroup.Wait()
	logger.INFO.Println("Stopping application...")
}
