package crawler

type HTMLPage struct {
	id 		int
	website string
	content []byte
}

func NewHTMLPage(response httpResponse, id int, website string) HTMLPage {
	return HTMLPage{
		id: id,
		website: website,
		content: response.getContent(),
	}
}

func (t HTMLPage) GetId() int {
	return t.id
}

func (t HTMLPage) GetWebsite() string {
	return t.website
}

func (t HTMLPage) GetContent() []byte {
	return t.content
}
