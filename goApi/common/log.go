package common

import (
	"os"
	"path/filepath"
)

func LogToFile(m3u8Path string, logContent string) (err error) {
	// logFilePath := filepath.Dir(m3u8Path) + "/" + "log.txt"
	logFilePath := filepath.Join(filepath.Dir(m3u8Path), WorkPathName+"log.txt")
	handel, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer handel.Close()
	_, err = handel.WriteString(logContent)
	return err
}
