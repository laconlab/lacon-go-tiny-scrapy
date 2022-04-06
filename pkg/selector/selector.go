package selector

import (
	"log"
)

func NewHttpReqChan[T WebsiteReqest](websites *Websites, emptyReqests chan T) chan T {
	out := make(chan T, 10)

	go func(in, out chan T, sites []*Website) {
		defer close(out)

		idx := 0
		for req := range in {

			var site *Website
			for i := 0; i < len(sites); i++ {
				if !sites[idx].isDone() {
					site = sites[idx]
				}
				idx = (idx + 1) % len(sites)

				if site != nil {
					break
				}
			}

			if site == nil {
				break
			}

			req.SetId(site.Id)
			req.SetWebsite(site.Name)
			req.SetUrl(site.getUrl())
			site.inc()
			out <- req

		}
		log.Println("all requests are created, closing channel")
	}(emptyReqests, out, websites.Sites)

	return out
}
