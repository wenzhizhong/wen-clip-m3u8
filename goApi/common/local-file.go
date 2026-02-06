package common

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
)

// VideoInfo 结构体
type VideoInfo struct {
	Format struct {
		Duration string `json:"duration"`
		BitRate  string `json:"bit_rate"`
		Size     string `json:"size"`
	} `json:"format"`
	Streams []struct {
		CodecType  string `json:"codec_type"`
		CodecName  string `json:"codec_name"`
		Width      int    `json:"width,omitempty"`
		Height     int    `json:"height,omitempty"`
		Duration   string `json:"duration,omitempty"`
		BitRate    string `json:"bit_rate,omitempty"`
		RFrameRate string `json:"r_frame_rate,omitempty"`
	} `json:"streams"`
}

// 使用通配符删除文件
func RemoveByWildcard(dir, pattern string) error {
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return err
	}

	for _, match := range matches {
		fmt.Printf("删除: %s\n", match)
		if err := os.RemoveAll(match); err != nil {
			return err
		}
	}
	return nil
}

// 通过通配符获取文件列表
func ListByWildcard(dir, pattern string) (matches []string, err error) {
	matches, err = filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return
	}

	return
}

// 获取视频信息
func GetVideoInfoJSON(filePath string) (*VideoInfo, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filePath)

	// 设置 Windows 下不显示窗口
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var info VideoInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func GetSize(path string) (size int64, err error) {
	size = 0
	// 获取文件或者目录大小
	var fileInfo os.FileInfo
	fileInfo, err = os.Stat(path)
	if err != nil {
		return
	}
	if fileInfo.IsDir() {
		files, err1 := ListByWildcard(path, "*")
		if err1 != nil {
			return size, err1
		}
		for _, filePath := range files {
			tmpSize, err2 := GetSize(filePath)
			if err2 != nil {
				return size, err2
			}
			size += tmpSize
		}
	} else {
		size += fileInfo.Size()
	}

	return size, nil
}
