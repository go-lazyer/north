package nslice

import "fmt"

// 切割slice 每shareNums一份,最后一个shareNums可能比shareNums少
func Cut[T any](slice []T, shareNums int) [][]T {
	// 处理无效的shareNums输入
	if shareNums <= 0 {
		return [][]T{}
	}

	sliceLen := len(slice)
	// 计算实际的组数（向上取整）
	numGroups := (sliceLen + shareNums - 1) / shareNums
	resSlice := make([][]T, 0, numGroups) // 预分配内存提高性能

	for i := 0; i < numGroups; i++ {
		startIndex := i * shareNums
		endIndex := startIndex + shareNums
		if endIndex > sliceLen {
			endIndex = sliceLen
		}
		// 确保只添加非空分组
		if startIndex < endIndex {
			resSlice = append(resSlice, slice[startIndex:endIndex])
		}
	}
	return resSlice
}

// 将[]map[string]any 转为 [][]any, 第一行是标题行
func ToCsv(data []map[string]any) [][]string {
	if len(data) == 0 {
		return [][]string{}
	}
	// 提取标题行
	headers := make([]string, 0, len(data[0]))
	for key := range data[0] {
		headers = append(headers, key)
	}

	result := make([][]string, 0, len(data)+1)
	result = append(result, headers)

	// 提取数据行
	for _, record := range data {
		row := make([]string, 0, len(record))
		for _, header := range headers {
			val, ok := record[header]
			if !ok || val == nil {
				row = append(row, "")
			} else {
				row = append(row, fmt.Sprintf("%v", val))
			}
		}
		result = append(result, row)
	}
	return result
}
