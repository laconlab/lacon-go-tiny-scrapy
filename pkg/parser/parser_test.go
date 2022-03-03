package parser_test

import (
	"testing"

	"github.com/laconlab/lacon-go-tiny-scrapy/pkg/parser"
	"gopkg.in/yaml.v2"
)

type PageMock struct {
	website string
	content []byte
	result  map[string]string
}

func (p *PageMock) GetWebsite() string {
	return p.website
}

func (p *PageMock) GetContent() []byte {
	return p.content
}

func (p *PageMock) SetParsedContent(res map[string]string) {
	p.result = res
}

func TestParsing(t *testing.T) {
	rulesStr := `
    websites:
    -   name: "site"
        rules:
        - component: "title"
          tag: "h1"
          attribute-name: "class"
          attribute-id: "class-id"
        - component: "text"
          tag: "div"
          attribute-name: "class"
          attribute-id: "div-class"
        - component: "exec"
          tag: "p"
          attribute-name: "style"
          attribute-id: "style-1"
          extract: "exec"
        - component: "filtered"
          tag: "div"
          attribute-name: "class"
          attribute-id: "div-class-3"
          filters: ["iframe"]
    `

	webPage := `
	<h1 class="class-id">expected_title</h1>
	<h1 class="class-id-2">unexpected_title</h1>
	<h1 style="style-0">unexpected_title</h1>
	<div class="div-class">
		<p style="style-0">Paragraph 1</p>
		<p>Paragraph 2</p>
		<p>Paragraph 3</p>
	</div>
	<div class="div-class-2">
		<p style="style-1", exec="id-0">Paragraph 4</p>
	</div>
	<div class="div-class-3">
		<p>Paragraph 5</p>
		<iframe>
			<p>Paragraph 6</p>
		</iframe>
	</div>
	`

	rules := &parser.ParserRules{}
	if err := yaml.Unmarshal([]byte(rulesStr), rules); err != nil {
		t.Error(err)
	}

	parser := parser.NewParser(rules)

	page := &PageMock{
		website: "site",
		content: []byte(webPage),
	}
	parser.Parse(page)
	if len(page.result) != 4 {
		t.Errorf("Expected result of size 4, got %d\n", len(page.result))
	}

	if val, ok := page.result["title"]; !ok || val != "expected_title" {
		t.Errorf("Expected title: expected_title, got: %+v\n", page.result)
	}

	if val, ok := page.result["text"]; !ok || len(val) != 39 {
		t.Errorf("Unexpected text %+v\n", page.result)
	}

	if val, ok := page.result["exec"]; !ok || val != "id-0" {
		t.Errorf("Unexpected value of exec: %+v\n", page.result)
	}

	if val, ok := page.result["filtered"]; !ok || val != "Paragraph 5" {
		t.Errorf("Unexpected value of filtered: %+v\n", page.result)
	}
}
