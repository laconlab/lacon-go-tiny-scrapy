package result // TODO find better name

import (
	"encoding/json"
	"log"
)

type FullWebsiteResult struct {
	Id            int               `json:"id"`
	Website       string            `json:"website"`
	Url           string            `json:"url"`
	DownloadDate  string            `json:"download_date"`
	RawContent    []byte            `json:"html"`
	ParsedContent map[string]string `json:"content"` // TODO think about it
}

func NewFullWebsiteResultChan() chan *FullWebsiteResult {
	out := make(chan *FullWebsiteResult)
	go func(out chan *FullWebsiteResult) {
		for {
			out <- &FullWebsiteResult{
				Id:            0,
				Website:       "",
				Url:           "",
				DownloadDate:  "",
				RawContent:    []byte{},
				ParsedContent: map[string]string{},
			}
		}
	}(out)
	return out
}

func (w *FullWebsiteResult) SetId(id int) {
	w.Id = id
}

func (w *FullWebsiteResult) SetWebsite(website string) {
	w.Website = website
}

func (w *FullWebsiteResult) GetWebsite() string {
	return w.Website
}

func (w *FullWebsiteResult) SetUrl(url string) {
	w.Url = url
}

func (w *FullWebsiteResult) GetUrl() string {
	return w.Url
}

func (w *FullWebsiteResult) SetDownloadDate(date string) {
	w.DownloadDate = date
}

func (w *FullWebsiteResult) SetRawContent(cnt []byte) {
	w.RawContent = cnt
}

func (w *FullWebsiteResult) GetRawContent() []byte {
	return w.RawContent
}

func (w *FullWebsiteResult) SetParsedContent(cnt map[string]string) {
	w.ParsedContent = cnt
}

func (w *FullWebsiteResult) GetRawWebsiteAsJSON() []byte {
	b, err := json.Marshal(w)
	if err != nil {
		log.Println("Cannot convert to json", err)
		return nil
	}
	return b
}
