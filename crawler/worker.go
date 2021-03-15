package crawler

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"lacon-go-tiny-scrapy/logger"
	"lacon-go-tiny-scrapy/selector"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	maxRetryCount int
	timeout       time.Duration
	tempStopTime  time.Duration
	poolSize      int

	workersGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "crawler_worker_count",
		Help: "Number of active crawlers",
	})

	stoppedWorkerGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "number_of_stopped_workers",
		Help: "Number of tmp stopped workers",
	})

	timeoutsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "timeouts_count",
		Help: "Number of timeouts per website",
	})

	downloadErrorCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "download_error_count",
		Help: "Number of errors while downloading",
	})

	downloadSummary = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "crawler_worker_summary",
		Help: "Performances per worker",
	},
	[]string{
		"website",
		"status",
	})

	responseCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "successful_download_count",
		Help: "Number of successful downloads",
	},
	[]string{
		"website",
	})
)

func init() {
	config := newCrawlerConfiguration()
	maxRetryCount = config.getRetryCount()
	timeout = config.getTimeout()
	poolSize = config.getWorkerPoolSize()
	tempStopTime = config.getTempStopTime()


	logger.INFO.Printf("Downloading with poolSize: %d, timeout: %v, " +
		"retryCount: %d, tempLock %v\n", poolSize, timeout, maxRetryCount, tempStopTime)
}


func NewHTMLProducer(waitingGroup *sync.WaitGroup, inputRequests <- chan selector.HTTPRequest) <- chan HTMLPage {
	waitingGroup.Add(1)
	outputHtml := newHTMLOutputChannel()

	innerWaitingGroup := sync.WaitGroup{}

	innerWaitingGroup.Add(poolSize)
	logger.INFO.Printf("Starting %d download workers\n", poolSize)
	workersGauge.Set(float64(poolSize))
	for i := 0; i < poolSize; i++ {
		go startDownloadWorker(&innerWaitingGroup, inputRequests, outputHtml)
	}

	go closeOutput(&innerWaitingGroup, waitingGroup, outputHtml)

	return outputHtml
}

func newHTMLOutputChannel() chan HTMLPage {
	return make(chan HTMLPage)
}

func startDownloadWorker(waitingGroup *sync.WaitGroup, input <- chan selector.HTTPRequest, output chan <- HTMLPage) {
	for req := range input {
		startTime := time.Now().UnixNano()
		response := downloadWithRetry(req.GetURL())

		if response.isSuccess() {
			responseCounter.WithLabelValues(req.GetWebsite()).Inc()
			output <- NewHTMLPage(response, req.GetId(), req.GetWebsite())
		}

		endTime := time.Now().UnixNano()
		downloadSummary.WithLabelValues(req.GetWebsite(), strconv.Itoa(int(response.status))).Observe(float64(endTime - startTime))
	}
	workersGauge.Dec()
	waitingGroup.Done()
}

func downloadWithRetry(url string) httpResponse {
	var response httpResponse

	for tryCount := 0; tryCount < maxRetryCount; tryCount++ {
		response = download(url)

		if response.isForRetry() {
			logger.WARNING.Printf("Retry of %s, retry count %d\n", url, tryCount)
			stoppedWorkerGauge.Inc()
			time.Sleep(tempStopTime)
			stoppedWorkerGauge.Dec()
			continue
		}
		return response
	}
	return createFailedHTTPResponse()
}

func download(url string) httpResponse {
	callback := make(chan httpResponse, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go asyncDownload(url, ctx, callback)

	select {
	case ret := <- callback:
		return ret
	case <-time.After(timeout):
		logger.WARNING.Printf("Timeout on url: %s\n", url)
		timeoutsCounter.Inc()
		return createFailedHTTPResponse()
	}
}

func asyncDownload(url string, ctx context.Context, callback chan httpResponse) {
	select {
	default:
		response, err := http.Get(url)
		if err == nil {
			callback <- newHTTPResponse(response)
		} else {
			logger.ERROR.Printf("Error while downloading %s, %v", url, err)
			downloadErrorCounter.Inc()
		}
	case <-ctx.Done():
		return
	}
}

func closeOutput(innerWaitingGroup *sync.WaitGroup, waitingGroup *sync.WaitGroup, outputChannel chan HTMLPage) {
	logger.INFO.Println("Waiting to download all pages...")
	innerWaitingGroup.Wait()
	logger.INFO.Println("Download completed, closing the channel")
	close(outputChannel)
	logger.INFO.Println("Download completed, channel closed")
	waitingGroup.Done()
}
