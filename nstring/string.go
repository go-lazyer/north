package nstring

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

// 首字母大写驼峰
func ToUpperCamelCase(str string) string {
	if len(str) == 0 {
		return ""
	}
	strArr := strings.Split(str, "_")
	var sb bytes.Buffer
	for _, s := range strArr {
		if len(s) == 0 {
			continue
		}
		sb.WriteString(strings.ToUpper(s[0:1]) + s[1:])
	}
	return sb.String()
}

// 首字母小写驼峰
func ToLowerCamelCase(str string) string {
	if len(str) == 0 {
		return ""
	}
	strArr := strings.Split(str, "_")
	var sb bytes.Buffer
	for n, s := range strArr {
		if len(s) == 0 {
			continue
		}
		if n == 0 {
			sb.WriteString(s)
		} else {
			sb.WriteString(strings.ToUpper(s[0:1]) + s[1:])
		}
	}
	return sb.String()
}

// 首字母大写
func ToUpperFirst(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.ToUpper(str[0:1]) + str[1:]
}

// 压缩json
func CompactJson(jsonStr string) (string, error) {
	if jsonStr == "" {
		return "", nil
	}
	jsonStr = strings.ReplaceAll(jsonStr, "\n", "")

	// 压缩json
	var requestBuffer bytes.Buffer
	if err := json.Compact(&requestBuffer, []byte(jsonStr)); err != nil {
		return "", errors.New("json压缩失败: " + err.Error())
	}
	return requestBuffer.String(), nil
}
