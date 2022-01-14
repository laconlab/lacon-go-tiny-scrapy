package crawler

import (
	"log"
	"net/http"
	"time"
)

type agents interface {
	next() string
}

type headerFilter struct {
	ags     agents
	timeout time.Duration
}

func newHeaderFilter(ags agents, timeout time.Duration) *headerFilter {
	return &headerFilter{
		ags:     ags,
		timeout: timeout,
	}
}

func (h *headerFilter) getAgent() string {
	return h.ags.next()
}

// return true if client does not support header download
// or if http status is not between [300, 500)
// otherwise false
func (h *headerFilter) filter(url string) bool {
	client := &http.Client{
		Timeout: h.timeout,
	}

	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		log.Println(err)
		return true
	}

	req.Header.Set("User-Agent", h.getAgent())

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return true
	}

	return !(resp.StatusCode >= 300 && resp.StatusCode < 500)
}
