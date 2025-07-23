package nslice

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
