package nsql

import "database/sql"

func ScanToMap(rows *sql.Rows) ([]map[string]any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]any, len(columns))   // 存放每一列的值的指针
	scanArgs := make([]any, len(columns)) // Scan() 的参数
	for i := range values {
		scanArgs[i] = &values[i] // 每个值的地址
	}

	var result []map[string]any

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		record := make(map[string]any)
		for i, col := range columns {
			val := values[i]
			// 处理 nil 值（比如数据库中的 NULL）
			if val == nil {
				record[col] = nil
			} else {
				// 根据实际类型做断言或转换（可选）
				if bv, ok := val.([]byte); ok {
					record[col] = string(bv)
				} else {
					record[col] = val
				}
			}
		}

		result = append(result, record)
	}

	return result, nil
}

// ScanToSlice 将 rows 扫描为 [][]any
// 返回: [][列值]，例如 [[1 "alice" "a@b.com"] [2 "bob" "b@c.com"]]
func ScanToSlice(rows *sql.Rows) ([][]any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 准备接收值的切片
	values := make([]any, len(columns))
	scanArgs := make([]any, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var result [][]any

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		// 复制当前行，避免被下一次 Scan 覆盖
		rowCopy := make([]any, len(values))
		for i, v := range values {
			// 可选：将 []byte 转为 string
			if bv, ok := v.([]byte); ok {
				rowCopy[i] = string(bv)
			} else {
				rowCopy[i] = v
			}
		}

		result = append(result, rowCopy)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
