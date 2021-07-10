package data_store

import (
	"os"
	"path/filepath"
	"sync"
)

type User struct {
	Name   string
	Number int
}

type JsonStore struct {
	waitGroup          *sync.WaitGroup
	max_file_size_byte int64
	current_index      int
	out_pattern        string
}

func NewJsonStore(
	out_dir string,
	file_name string,
	max_file_size_in_bytes int,
	waitGroup *sync.WaitGroup,
	users chan User) {

	waitGroup.Add(1)

	err := os.MkdirAll(out_dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	path_pattern = filepath.Join(out_dir, file_name+"_{}.json")

	worker := JsonlinesStore{
		waitGroup:          waitGroup,
		out_pattern:        path_pattern,
		max_file_size_byte: 30000,
		current_index:      0,
	}

	go worker.start(wg, users)
}
