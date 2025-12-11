package services

import (
	"io/ioutil"
)

type LocalFileService struct {
}

func (s *LocalFileService) PlayLocalFile(path string) (data *[]byte, err error) {
	content, err := ioutil.ReadFile(path)
	data = &content
	return data, err
}
