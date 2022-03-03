package selector

import (
	"log"

	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/result"
)

func NewHttpReqChan(websites *Websites) chan *result.FullWebsiteResult {
	out := make(chan *result.FullWebsiteResult, 10)

	go func(ch chan *result.FullWebsiteResult, sites []*Website) {
		defer close(ch)

		done := 0
		for done < len(sites) {
			done = 0
			for _, site := range sites {
				if site.isDone() {
					done++
					continue
				}

				req := &result.FullWebsiteResult{}
				req.SetId(site.getId())
				req.SetWebsite(site.getName())
				req.SetUrl(site.getUrl())
				site.inc()
				ch <- req

			}
		}
		log.Println("all requests are created, closing channel")
	}(out, websites.Sites)

	return out
}
