package parser

import (
	"flag"
	"lacon-go-tiny-scrapy/utils"
)

var ruleFieldsPerWebsite websitesFields
var filename = flag.String("parsing", "parsing-config.yml", "Location of the parsing rules file.")

func init() {
	utils.LoadConfig(filename, &ruleFieldsPerWebsite)
}

type websitesFields struct {
	Website map[string]fields `yaml:"website"`
}

func getFieldsByWebsite(website string) map[string]rule {
	return ruleFieldsPerWebsite.Website[website]
}
