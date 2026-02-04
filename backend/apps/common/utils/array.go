package utils

import (
	"reflect"
	"sort"
)

// 数组包含
func ArrayContains[T comparable](array []T, value T) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

// 数组分块
func ArrayChunk[T comparable](array []T, size int) [][]T {
	chunked := make([][]T, 0)
	for i := 0; i < len(array); i += size {
		chunk := array[i:min(i+size, len(array))]
		chunked = append(chunked, chunk)
	}
	return chunked
}

// ArraySort 对数组进行排序，支持升序和降序
func ArraySort[T comparable](array []T, sortType int) []T {
	// 创建副本避免修改原数组
	result := make([]T, len(array))
	copy(result, array)

	// 如果数组为空或只有一个元素，直接返回
	if len(result) <= 1 {
		return result
	}

	// 获取数组第一个元素的实际类型并进行相应排序
	firstElement := reflect.ValueOf(result[0])
	switch firstElement.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sortIntSliceInterface(result, sortType)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		sortUintSliceInterface(result, sortType)
	case reflect.Float32, reflect.Float64:
		sortFloatSliceInterface(result, sortType)
	case reflect.String:
		sortStringSliceInterface(result, sortType)
	default:
		// 对于其他类型，无法排序则返回原数组
		return result
	}

	return result
}

// sortIntSliceInterface 排序整数切片（通过interface）
func sortIntSliceInterface[T comparable](slice []T, sortType int) {
	if sortType == 1 { // 升序
		sort.Slice(slice, func(i, j int) bool {
			return convertToFloat64(slice[i]) < convertToFloat64(slice[j])
		})
	} else { // 降序
		sort.Slice(slice, func(i, j int) bool {
			return convertToFloat64(slice[i]) > convertToFloat64(slice[j])
		})
	}
}

// sortUintSliceInterface 排序无符号整数切片（通过interface）
func sortUintSliceInterface[T comparable](slice []T, sortType int) {
	if sortType == 1 { // 升序
		sort.Slice(slice, func(i, j int) bool {
			return convertToFloat64(slice[i]) < convertToFloat64(slice[j])
		})
	} else { // 降序
		sort.Slice(slice, func(i, j int) bool {
			return convertToFloat64(slice[i]) > convertToFloat64(slice[j])
		})
	}
}

// sortFloatSliceInterface 排序浮点数切片（通过interface）
func sortFloatSliceInterface[T comparable](slice []T, sortType int) {
	if sortType == 1 { // 升序
		sort.Slice(slice, func(i, j int) bool {
			return convertToFloat64(slice[i]) < convertToFloat64(slice[j])
		})
	} else { // 降序
		sort.Slice(slice, func(i, j int) bool {
			return convertToFloat64(slice[i]) > convertToFloat64(slice[j])
		})
	}
}

// sortStringSliceInterface 排序字符串切片（通过interface）
func sortStringSliceInterface[T comparable](slice []T, sortType int) {
	if sortType == 1 { // 升序
		sort.Slice(slice, func(i, j int) bool {
			return reflect.ValueOf(slice[i]).String() < reflect.ValueOf(slice[j]).String()
		})
	} else { // 降序
		sort.Slice(slice, func(i, j int) bool {
			return reflect.ValueOf(slice[i]).String() > reflect.ValueOf(slice[j]).String()
		})
	}
}

// convertToFloat64 将可比较类型转换为float64进行比较
func convertToFloat64(v interface{}) float64 {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return rv.Float()
	default:
		return 0.0
	}
}

// min 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
