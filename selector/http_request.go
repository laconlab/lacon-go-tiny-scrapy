package selector

type HTTPRequest struct {
	id  int
	url string
	websiteName string
}

func newHTTPRequest(id int, website string, url string) HTTPRequest {
	return HTTPRequest{
		id: id,
		websiteName: website,
		url: url,
	}
}

func (r HTTPRequest) GetId() int {
	return r.id
}

func (r HTTPRequest) GetWebsite() string {
	return r.websiteName
}

func (r HTTPRequest) GetURL() string {
	return r.url
}
