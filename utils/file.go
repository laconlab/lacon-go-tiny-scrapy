package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"lacon-go-tiny-scrapy/logger"
)


func loadFile(fileName string) []byte {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.ERROR.Fatal("Failed to loadFile", err)
	}
	return file
}

func setSettings(settings interface{}, file []byte) {
	err := yaml.Unmarshal(file, settings)
	if err != nil {
		logger.ERROR.Fatal("Failed to setSettings", err)
	}
}
