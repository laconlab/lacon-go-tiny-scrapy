package parser

import (
	"bytes"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Parser struct {
	rules *ParserRules
}

func NewParser(rules *ParserRules) *Parser {
	return &Parser{rules}
}

func (p *Parser) Parse(page WebPage) {
	website := page.GetWebsite()
	rules := p.rules.GetByName(website)

	result := make(map[string]string)
	for _, rule := range rules {
		if res, ok := extract(page.GetContent(), rule); ok {
			result[rule.Name] = res
		}
	}
	page.SetParsedContent(result)
}

func extract(cnt []byte, rule *Rule) (string, bool) {
	reader := bytes.NewReader(cnt)
	html, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Printf("Error while extracting %+v %+v\n", rule, err)
		return "", false
	}

	var result []string
	var found bool
	html.Find(rule.Tag).Each(func(_ int, sec *goquery.Selection) {
		if res, ok := findInSection(sec, rule); ok {
			result = append(result, strings.TrimSpace(res))
			found = true
			return
		}
	})

	return strings.Join(result, "\n"), found
}

func findInSection(selection *goquery.Selection, rule *Rule) (string, bool) {
	var res string
	val, found := selection.Attr(rule.AttrName)
	if rule.AttrName == "" {
		found = true
	}

	if found && val == rule.AttrId && rule.Extract != "" {
		res, found = selection.Attr(rule.Extract)
	} else if found && val == rule.AttrId {
		filterSelection(selection, rule.Filters)
		res = selection.Text()
	} else {
		found = false
	}

	return res, found
}

func filterSelection(selection *goquery.Selection, tags []string) {
	for _, tag := range tags {
		selection.Find(tag).Each(func(_ int, sec *goquery.Selection) {
			sec.Remove()
		})
	}
}
