package selector

func NewHttpReqChan(websites *Websites) chan interface{} {
	out := make(chan interface{}, 10)

	go func(ch chan interface{}, sites []*Website) {
		defer close(ch)

		done := 0
		for done != len(sites) {
			done = 0
			for _, site := range sites {
				if site.isDone() {
					done++
					continue
				}

				ch <- &HttpRequest{
					id:   site.getId(),
					name: site.getName(),
					url:  site.getUrl(),
				}

				site.inc()
			}
		}
	}(out, websites.Sites)

	return out
}
