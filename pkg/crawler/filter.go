package crawler

import (
	"log"
	"net/http"
	"time"
)

// return true if there is an error (timeout included)
// or if http status is not between [300, 500)
// otherwise false
func headerFilter(url string, agent string, timeout time.Duration) bool {
	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		log.Println(err)
		return true
	}

	req.Header.Set("User-Agent", agent)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return true
	}

	return !(resp.StatusCode >= 300 && resp.StatusCode < 500)
}
