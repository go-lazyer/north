package nfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsExist(t *testing.T) {
	// 创建临时文件
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_exist.txt")
	if err := os.WriteFile(tmpFile, []byte("hello"), 0644); err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}

	// 测试存在的文件
	if !IsExist(tmpFile) {
		t.Errorf("IsExist(%s) 期望 true，实际 false", tmpFile)
	}

	// 测试存在的目录
	if !IsExist(tmpDir) {
		t.Errorf("IsExist(%s) 期望 true，实际 false", tmpDir)
	}

	// 测试不存在的路径
	if IsExist(filepath.Join(tmpDir, "not_exist.txt")) {
		t.Errorf("IsExist(不存在的路径) 期望 false，实际 true")
	}
}

func TestCreateDir(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建新目录
	newDir := filepath.Join(tmpDir, "new_sub_dir")
	if err := CreateDir(newDir); err != nil {
		t.Fatalf("CreateDir 失败: %v", err)
	}
	if !IsExist(newDir) {
		t.Errorf("CreateDir 后目录应存在")
	}

	// 创建已存在的目录（应返回 nil）
	if err := CreateDir(newDir); err != nil {
		t.Errorf("CreateDir 已存在的目录应返回 nil，实际 %v", err)
	}

	// 创建多层嵌套目录
	nestedDir := filepath.Join(tmpDir, "a", "b", "c")
	if err := CreateDir(nestedDir); err != nil {
		t.Fatalf("CreateDir 嵌套目录失败: %v", err)
	}
	if !IsExist(nestedDir) {
		t.Errorf("CreateDir 嵌套目录后应存在")
	}
}

func TestReadToLines(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "lines.txt")

	content := "line1\nline2\nline3"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}

	lines, err := ReadToLines(tmpFile)
	if err != nil {
		t.Fatalf("ReadToLines 失败: %v", err)
	}

	expected := []string{"line1", "line2", "line3"}
	if len(lines) != len(expected) {
		t.Fatalf("期望 %d 行，实际 %d 行", len(expected), len(lines))
	}
	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("第 %d 行期望 %q，实际 %q", i, expected[i], line)
		}
	}

	// 测试不存在的文件
	_, err = ReadToLines(filepath.Join(tmpDir, "not_exist.txt"))
	if err == nil {
		t.Errorf("读取不存在的文件应返回错误")
	}
}

func TestReadToString(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "string.txt")

	content := "hello world"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}

	result, err := ReadToString(tmpFile)
	if err != nil {
		t.Fatalf("ReadToString 失败: %v", err)
	}
	if result != content {
		t.Errorf("期望 %q，实际 %q", content, result)
	}

	// 测试不存在的文件
	_, err = ReadToString(filepath.Join(tmpDir, "not_exist.txt"))
	if err == nil {
		t.Errorf("读取不存在的文件应返回错误")
	}
}

func TestReadFolder(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建一些文件和子目录
	os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "b.txt"), []byte("b"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".hidden"), []byte("h"), 0644) // 隐藏文件
	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)                   // 子目录

	files, err := ReadFolder(tmpDir)
	if err != nil {
		t.Fatalf("ReadFolder 失败: %v", err)
	}

	// 应只包含非隐藏文件，不包含子目录
	expectedCount := 2 // a.txt, b.txt
	if len(files) != expectedCount {
		t.Fatalf("期望 %d 个文件，实际 %d 个", expectedCount, len(files))
	}

	// 检查不包含隐藏文件
	for _, f := range files {
		if filepath.Base(f) == ".hidden" {
			t.Errorf("ReadFolder 不应包含隐藏文件 .hidden")
		}
	}

	// 检查不包含子目录路径
	for _, f := range files {
		if filepath.Base(f) == "subdir" {
			t.Errorf("ReadFolder 不应包含子目录")
		}
	}
}

func TestEachFolder(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建嵌套结构
	os.WriteFile(filepath.Join(tmpDir, "root.txt"), []byte("r"), 0644)
	subDir := filepath.Join(tmpDir, "sub")
	os.Mkdir(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "sub_file.txt"), []byte("s"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".hidden_root"), []byte("h"), 0644)

	files, err := EachFolder(tmpDir)
	if err != nil {
		t.Fatalf("EachFolder 失败: %v", err)
	}

	// 应包含子目录中的文件，不包含隐藏文件
	for _, f := range files {
		base := filepath.Base(f)
		if base == ".hidden_root" {
			t.Errorf("EachFolder 不应包含隐藏文件")
		}
	}

	// 应包含 root.txt 和 sub_file.txt
	foundRoot := false
	foundSub := false
	for _, f := range files {
		if filepath.Base(f) == "root.txt" {
			foundRoot = true
		}
		if filepath.Base(f) == "sub_file.txt" {
			foundSub = true
		}
	}
	if !foundRoot {
		t.Errorf("EachFolder 应包含 root.txt")
	}
	if !foundSub {
		t.Errorf("EachFolder 应包含 sub/sub_file.txt")
	}
}
