package data_store

import (
	"os"
	"path/filepath"
	"sync"
	"encoding/json"
	"lacon-go-tiny-scrapy/logger"
	"lacon-go-tiny-scrapy/utils"
)

type User struct {
	Name   string
	Number string
}


func StoreJsonLines(out_dir string, file_name string, waitGroup *sync.WaitGroup, users chan User) {
	utils.CreateFolder(out_dir)
	file := utils.CreateFile(filepath.Join(out_dir, file_name + ".json"))

	waitGroup.Add(1)
	go startStoring(file, waitGroup, users)
}


func startStoring(file *os.File, waitGroup *sync.WaitGroup, users chan User) {
	defer waitGroup.Done()
	defer file.Close()

	for user := range users {
		json, _ := json.Marshal(user)
		file.WriteString(string(json) + "\n")
	}

	logger.INFO.Println("Saving completed, closing file")
}


