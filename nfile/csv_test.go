package nfile

import (
	"path/filepath"
	"testing"
)

func TestWriterToCsvFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.csv")

	content := [][]string{
		{"name", "age", "city"},
		{"张三", "25", "北京"},
		{"李四", "30", "上海"},
	}

	// 写入 CSV
	if err := WriterToCsvFile(tmpFile, content); err != nil {
		t.Fatalf("WriterToCsvFile 失败: %v", err)
	}

	// 验证文件存在
	if !IsExist(tmpFile) {
		t.Fatalf("CSV 文件应存在")
	}

	// 读取并验证内容
	lines, err := ReadToLines(tmpFile)
	if err != nil {
		t.Fatalf("读取 CSV 文件失败: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("期望 3 行，实际 %d 行", len(lines))
	}

	// 测试空内容
	emptyFile := filepath.Join(tmpDir, "empty.csv")
	if err := WriterToCsvFile(emptyFile, [][]string{}); err != nil {
		t.Fatalf("WriterToCsvFile 空内容失败: %v", err)
	}

	// 测试覆盖写入（TRUNC 模式）
	if err := WriterToCsvFile(tmpFile, [][]string{{"only"}}); err != nil {
		t.Fatalf("WriterToCsvFile 覆盖写入失败: %v", err)
	}
	lines, err = ReadToLines(tmpFile)
	if err != nil {
		t.Fatalf("读取覆盖后的 CSV 文件失败: %v", err)
	}
	if len(lines) != 1 {
		t.Fatalf("覆盖后期望 1 行，实际 %d 行", len(lines))
	}
	if lines[0] != "only" {
		t.Errorf("覆盖后期望 'only'，实际 %q", lines[0])
	}

	// // 测试无效路径
	// _, err := WriterToCsvFile("/nonexistent/path/test.csv", content)
	// if err == nil {
	// 	t.Errorf("写入无效路径应返回错误")
	// }
}
