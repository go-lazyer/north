package nfile

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err == nil {
		return true
	}
	if os.IsExist(err) {
		return true
	}
	return false
}

// 只能创建目录，不能创建文件
func CreateDir(path string) error {
	if IsExist(path) {
		return nil
	}
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// 按照行读取文件，转换为切片（禁止读取大文件）
func ReadToLines(pathName string) ([]string, error) {
	file, err := os.Open(pathName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bufio := bufio.NewReader(file)
	lines := make([]string, 0)
	for {
		line, _, c := bufio.ReadLine()
		if c == io.EOF {
			break
		}
		lines = append(lines, string(line))
	}
	return lines, nil
}

// 读取文件转为字符串（禁止读取大文件）
func ReadToString(pathName string) (string, error) {
	content, err := os.ReadFile(pathName)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// 读取文件夹中的文件，不包含子文件夹
func ReadFolder(fileFullPath string) ([]string, error) {
	files, err := os.ReadDir(fileFullPath)
	if err != nil {
		return nil, err
	}
	myFiles := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		myFiles = append(myFiles, filepath.Join(fileFullPath, file.Name()))
	}
	return myFiles, nil
}

// 遍历文件夹中的所有文件,包含子文件夹中的内容
func EachFolder(rootPath string) ([]string, error) {

	myFiles := make([]string, 0)
	filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}
		if d.IsDir() {
			fmt.Printf("Directory: %s\n", path)
		} else {
			myFiles = append(myFiles, path)
			fmt.Printf("File: %s\n", path)
		}
		return nil
	})

	return myFiles, nil
}
