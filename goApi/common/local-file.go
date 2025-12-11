package common

import (
	"fmt"
	"os"
	"path/filepath"
)

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
