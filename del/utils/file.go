package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"lacon-go-tiny-scrapy/logger"
	"os"
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

func CreateFolder(folderName string) {
	logger.INFO.Printf("Creating folder %s/n", folderName)
	err := os.MkdirAll(folderName, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func CreateFile(fileName string) *os.File {

	logger.INFO.Printf("Creating new file %s", fileName)

	file, err := os.Create(fileName)
	if err != nil {
		logger.ERROR.Println("Cannot create file")
		panic(err)
	}

	logger.INFO.Printf("File %s created", fileName)
	return file
}
