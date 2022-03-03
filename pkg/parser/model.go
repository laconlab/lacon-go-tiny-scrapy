package parser

import (
	"log"
)

type WebPage interface {
	GetWebsite() string
	GetContent() []byte
	SetParsedContent(map[string]string)
}

type ParserRules struct {
	Rules []*ParserRule `yaml:"websites"`
	cache map[string][]*Rule
}

func (p *ParserRules) GetByName(website string) []*Rule {
	if rule, ok := p.cache[website]; ok {
		return rule
	}

	var found *ParserRule
	for _, rule := range p.Rules {
		if rule.Name == website {
			found = rule
			break
		}
	}

	if found == nil {
		log.Fatal("Cannot find rules for ", website)
		return nil
	}

	if p.cache == nil {
		p.cache = make(map[string][]*Rule)
	}

	p.cache[website] = found.Rules

	return found.Rules
}

type ParserRule struct {
	Name  string  `yaml:"name"`
	Rules []*Rule `yaml:"rules"`
}

type Rule struct {
	Name     string   `yaml:"component"`
	Tag      string   `yaml:"tag"`
	AttrName string   `yaml:"attribute-name"`
	AttrId   string   `yaml:"attribute-id"`
	Extract  string   `yaml:"extract"`
	Filters  []string `yaml:"filters"`
}
