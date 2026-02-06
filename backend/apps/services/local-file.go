package services

import (
	"clipM3u8Media/goApi"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type LocalFileService struct {
}

func (s *LocalFileService) PlayLocalFile(path, m3u8Path string) (data *[]byte, err error) {
	isVideo := false
	isVideo, err = s.checkIsVideo(path)
	if err != nil {
		return nil, err
	}
	if isVideo {
		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			m3u8Handler := goApi.M3u8Handler{}

			content, err := m3u8Handler.CheckM3u8File(m3u8Path)
			m3u8Info, _, err1 := m3u8Handler.ParseM3u8File(m3u8Path, &content)
			err = err1
			if err != nil {
				return nil, err
			}
			// listSlice := m3u8Info.ExtList
			sliceLen := len(m3u8Info.ExtList)
			sliceIndex := 0
			sliceIndexStr := "0"
			sliceIndexStr, err = s.getSliceIndexFromPath(path)
			if err != nil {
				return nil, err
			}
			sliceIndex, err = strconv.Atoi(sliceIndexStr)
			if err != nil {
				return nil, err
			}
			if sliceIndex >= sliceLen {
				return nil, errors.New("slice index is error")
			}

			curSlice := []string{m3u8Info.ExtList[sliceIndex]}
			m3u8Handler.DoGetM3u8SliceVideo(m3u8Path, m3u8Info, curSlice, goApi.OptTypeVideo)
		}
	}
	content, err := ioutil.ReadFile(path)
	data = &content
	return data, err
}

func (s *LocalFileService) checkIsVideo(path string) (bool, error) {
	if path == "" {
		return false, errors.New("path is empty")
	}
	path = strings.ToLower(path)
	return strings.Contains(path, ".mp4") && !strings.Contains(path, ".jpg") && !strings.Contains(path, ".jpeg") && !strings.Contains(path, ".png"), nil
}

func (s *LocalFileService) getSliceIndexFromPath(path string) (string, error) {
	path = strings.ReplaceAll(path, "\\", "/")
	pathArr := strings.Split(path, "/")
	if len(pathArr) < 2 {
		return "", errors.New("path is error")
	}
	fileName := pathArr[len(pathArr)-1]
	return strings.Split(fileName, ".")[0], nil
}
