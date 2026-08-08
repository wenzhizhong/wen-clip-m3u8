package tests

import (
	"clipM3u8Media/backend/apps/common/utils"
	"fmt"
	"testing"
)

func TestArrayDiff(t *testing.T) {
	result := utils.ArrayDiff([]string{"1", "2", "3", "5"}, []string{"2", "3", "4"}, false)
	result2 := utils.ArrayDiff([]string{"2", "3", "4"}, []string{"1", "2", "3", "5"}, false)
	result3 := utils.ArrayDiff([]string{"1", "2", "3", "5"}, []string{"2", "3", "4"}, true)
	result4 := utils.ArrayDiff([]string{"2", "3", "4"}, []string{"1", "2", "3", "5"}, true)
	fmt.Println("result = ", result)
	fmt.Println("result2 = ", result2)
	fmt.Println("result3 = ", result3)
	fmt.Println("result4 = ", result4)
	fmt.Println("==================")
}
