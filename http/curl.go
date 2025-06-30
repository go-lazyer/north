package nhttp

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ToCurl 将 *http.Request 对象转换为 cURL 命令字符串
func RequestToCurl(req *http.Request) string {
	var curlCmd strings.Builder

	// 1. 添加基础命令和方法
	curlCmd.WriteString("curl -X ")
	curlCmd.WriteString(req.Method)

	// 2. 处理 URL（包含查询参数）
	curlCmd.WriteString(" '")
	curlCmd.WriteString(req.URL.String())
	curlCmd.WriteString("' \\\n")

	// 3. 处理请求头
	for key, values := range req.Header {
		for _, value := range values {
			// 转义单引号防止命令中断（' -> '\'')
			escapedValue := strings.ReplaceAll(value, "'", `'\''`)
			fmt.Fprintf(&curlCmd, "-H '%s: %s' \\\n", key, escapedValue)
		}
	}

	// 4. 处理请求体
	if req.Body != nil {
		// 复制原始 Body（避免读取后丢失）
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return fmt.Sprintf("读取请求体失败: %w", err)
		}
		defer func() {
			// 重置 Body 以便后续使用
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}()

		// 特殊处理：空请求体
		if len(bodyBytes) == 0 {
			curlCmd.WriteString("-d ''")
			return curlCmd.String()
		}

		// 转义单引号和换行符
		escapedBody := strings.ReplaceAll(string(bodyBytes), "'", `'\''`)
		escapedBody = strings.ReplaceAll(escapedBody, "\n", `\n`)

		// 根据内容类型决定格式化方式
		contentType := req.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			fmt.Fprintf(&curlCmd, " -d '%s'", escapedBody)
		} else {
			fmt.Fprintf(&curlCmd, " -d '%s'", escapedBody)
		}
	}
	return curlCmd.String()
}
