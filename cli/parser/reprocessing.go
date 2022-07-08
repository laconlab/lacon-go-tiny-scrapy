// use for old data
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/parser"
	"gopkg.in/yaml.v2"
)

type record struct {
	Id           string `json:"id"`
	Site         string `json:"website"`
	rawHtml      []byte
	Result       map[string]string `json:"content"`
	DownloadDate string            `json:"download_date"`
}

func (r *record) GetWebsite() string {
	return r.Site
}

func (r *record) GetContent() []byte {
	return r.rawHtml
}

func (r *record) SetParsedContent(res map[string]string) {
	r.Result = res
}

func main() {
	cfg, err := ioutil.ReadFile("resources/application.yml")
	if err != nil {
		log.Fatal(err)
	}

	var loadPath string
	flag.StringVar(&loadPath, "lp", "./data/", "load path")

	var savePath string
	flag.StringVar(&savePath, "sp", "./data/", "save path")

	flag.Parse()

	records := make(chan *record)

	go func() {
		err := filepath.Walk(loadPath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !strings.HasSuffix(path, ".html") {
					return nil
				}

				dat, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				comps := strings.Split(info.Name(), "_")
				site := comps[0]
				id := strings.Replace(comps[1], ".html", "", 1)
				records <- &record{
					Id:           id,
					Site:         site,
					rawHtml:      dat,
					Result:       nil,
					DownloadDate: info.ModTime().String(),
				}

				return nil
			})

		if err != nil {
			log.Fatalf("Failed in walk %v \n", err)
		}

		close(records)
	}()

	rules := &parser.ParserRules{}
	if err := yaml.Unmarshal(cfg, rules); err != nil {
		log.Fatal(err)
	}

	p := parser.NewParser(rules)

	out := make(chan []byte, 100)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for rec := range records {
			wg.Add(1)
			go func(r *record) {
				defer wg.Done()
				p.Parse(r)

				if _, ok := r.Result["title"]; !ok {
					return
				}

				if _, ok := r.Result["url"]; !ok {
					return
				}

				if _, ok := r.Result["publish_date"]; !ok {
					return
				}

				if _, ok := r.Result["text"]; !ok {
					return
				}

				res, err := json.Marshal(r)
				if err != nil {
					return
				}

				out <- res

			}(rec)
		}

	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	f, err := os.Create(savePath + "result.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	for res := range out {
		f.Write(res)
		f.WriteString("\n")
	}
	f.Sync()

}
