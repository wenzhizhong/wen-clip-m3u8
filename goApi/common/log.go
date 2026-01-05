package common

import (
	"os"
	"path/filepath"
)

func LogToFile(m3u8Path string, logContent string) (err error) {
	logFilePath := filepath.Dir(m3u8Path) + "/" + "log.txt"
	haldel, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer haldel.Close()
	_, err = haldel.WriteString(logContent)
	return err
}
