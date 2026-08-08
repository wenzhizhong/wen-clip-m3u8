package utils

import "fmt"

// 微秒转“时:分:秒”
func MicrosecondToTime(microsecond int64) string {
	hour := microsecond / 1000000 / 60 / 60
	minute := (microsecond / 1000000 / 60) % 60
	second := (microsecond / 1000000) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
}
