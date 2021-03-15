package selector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	websiteStateGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "website_left_to_download_gauge",
			Help: "How many http requests remains per website",
		},
		[]string {
			"website",
		},
	)
)

type websiteState struct {
	name        string
	currentId   int
	endId       int
	urlTemplate string
}

func newWebsiteStateFromSettings(name string, setting websiteConfig) *websiteState {
	startId := setting.getStartIndex()
	endId := setting.getEndIndex()
	websiteStateGauge.WithLabelValues(name).Set(float64(endId - startId + 1))
	return &websiteState{
		name:        name,
		currentId:   startId,
		endId:       endId,
		urlTemplate: setting.getURLTemplate(),
	}
}

func (ws *websiteState) createNextHTTPRequest() HTTPRequest {
	url := ws.getURL()
	request := newHTTPRequest(ws.currentId, ws.name, url)
	ws.currentId++
	websiteStateGauge.WithLabelValues(ws.name).Dec()

	return request
}

func (ws *websiteState) isEndState() bool {
	return ws.currentId > ws.endId
}

func (ws *websiteState) getURL() string {
	return fmt.Sprintf(ws.urlTemplate, ws.currentId)
}
