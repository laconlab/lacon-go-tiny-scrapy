package crawler

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"io/ioutil"
	"lacon-go-tiny-scrapy/logger"
	"net/http"
	"strconv"
)

var (
	statusCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "downloads_status_counts",
			Help: "downloads per status",
		},
		[]string{
			"response_code",
		})
)

type httpResponse struct {
	status downloadStatus
	content []byte
}

func newHTTPResponse(response *http.Response) httpResponse {
	defer response.Body.Close()
	content := readHttpResponse(response)

	statusCode := response.StatusCode
	status := getDownloadStatusFromStatusCode(statusCode)
	statusCounter.WithLabelValues(strconv.Itoa(statusCode)).Inc()

	return httpResponse{
		status: status,
		content: content,
	}
}

func readHttpResponse(response *http.Response) []byte {
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.ERROR.Printf("Error while reading http response %v\n", err)
	}
	return content
}

func createFailedHTTPResponse() httpResponse {
	return httpResponse{
		status: undefined,
	}
}

func (r httpResponse) isSuccess() bool {
	return r.status == ok
}

func (r httpResponse) isForRetry() bool {
	return r.status == retry || r.status == undefined
}

func (r httpResponse) getContent() []byte {
	return r.content
}
