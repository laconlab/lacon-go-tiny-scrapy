package selector

import (
	"lacon-go-tiny-scrapy/logger"
	"lacon-go-tiny-scrapy/utils"
)

func newSelectorConfig() configuration {
	var config configuration
	utils.LoadConfig(utils.AppConfigurationFileName, &config)
	return config
}

type configuration struct {
	Selector struct {
		BufferSize int `yaml:"buffer-size"`
	}
}

func (c *configuration) getBufferSize() int {
	return c.Selector.BufferSize
}

func newWebsiteConfig() map[string]websiteConfig {
	var ws websitesConfigs
	utils.LoadConfig(utils.AppConfigurationFileName, &ws)
	logger.INFO.Printf("Using website configuration %v\n", ws)
	return ws.Website
}

type websitesConfigs struct {
	Website map[string]websiteConfig `yaml:"websites"`
}

func (ws *websitesConfigs) getWebsitesSettings() map[string]websiteConfig {
	return ws.Website
}

type websiteConfig struct {
	URLTemplate string `yaml:"url-template"`
	StartIndex  int    `yaml:"start-index"`
	EndIndex    int    `yaml:"end-index"`
}

func (ws *websiteConfig) getStartIndex() int {
	return ws.StartIndex
}

func (ws *websiteConfig) getEndIndex() int {
	return ws.EndIndex
}

func (ws *websiteConfig) getURLTemplate() string {
	return ws.URLTemplate
}
