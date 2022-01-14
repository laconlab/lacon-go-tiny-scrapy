package crawler

import (
	"log"
	"net/http"
)

type agents interface {
	next() string
}

type headerFilter struct {
	ags agents
}

func newHeaderFilter(ags agents) *headerFilter {
	return &headerFilter{
		ags: ags,
	}
}

func (h *headerFilter) getAgent() string {
	return h.ags.next()
}

// return true if client does not support header download
// or if http status is not between [300, 500)
// otherwise false
func (h *headerFilter) filter(url string) bool {
	client := &http.Client{}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		log.Fatalln(err)
		return true
	}

	req.Header.Set("User-Agent", h.getAgent())

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return true
	}

	return !(resp.StatusCode >= 300 && resp.StatusCode < 500)
}
