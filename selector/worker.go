package selector

import (
	"lacon-go-tiny-scrapy/logger"
	"sync"
)

var bufferSize int

func init() {
	config := newSelectorConfig()
	bufferSize = config.getBufferSize()
}

func NewHTTPRequestProducer(waitingGroup *sync.WaitGroup) chan HTTPRequest {
	waitingGroup.Add(1)

	states := websiteStateFromSettings()
	downloadRequestChannel := newHTTPRequestChannel()
	go startProducingHTTPRequests(waitingGroup, downloadRequestChannel, states)

	return downloadRequestChannel
}

func newHTTPRequestChannel() chan HTTPRequest {
	if bufferSize > 0 {
		logger.INFO.Printf("Creating http request channel with buffer size: %d\n", bufferSize)
		return make(chan HTTPRequest, bufferSize)
	} else {
		logger.INFO.Println("Creating unbuffered http request channel")
		return make(chan HTTPRequest)
	}
}

func startProducingHTTPRequests(
	waitGroup *sync.WaitGroup,
	output chan <- HTTPRequest,
	websitesStates *websitesStates) {

	logger.INFO.Println("Started inserting http requests in channel")

	for !websitesStates.isEndState() {
		request :=	websitesStates.createNextHTTPRequest()
		output <- request
	}

	logger.INFO.Println("All http requests made, closing the channel")
	close(output)
	waitGroup.Done()
}
