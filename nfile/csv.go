package nfile

import (
	"encoding/csv"
	"os"
)

func WriterToCsvFile(pathName string, content [][]string) error {
	file, err := os.OpenFile(pathName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	// 创建 CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush() // 确保所有数据都写入文件

	for _, record := range content {
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}
