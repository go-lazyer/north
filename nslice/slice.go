package nslice

// 切割slice 每shareNums一份,最后一个shareNums可能比shareNums少
func Cut[T any](slice []T, shareNums int) [][]T {

	resSlice := make([][]T, 0)

	for i := 0; i < len(slice)/shareNums+1; i++ {
		startIndex := i * shareNums
		endIndex := (i + 1) * shareNums
		if endIndex > len(slice) {
			endIndex = len(slice)
		}
		resSlice = append(resSlice, slice[startIndex:endIndex])
	}
	return resSlice
}
