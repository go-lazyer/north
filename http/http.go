package nhttp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	nfile "github.com/go-lazyer/north/file"
)

const (
	CONTENT_TYPE_STREAM = "application/octet-stream"
	CONTENT_TYPE_FORM   = "application/x-www-form-urlencoded"
	CONTENT_TYPE_JSON   = "application/json"
	CONTENT_TYPE_DATA   = "multipart/form-data"
)

func (r *Request) Do(method, u string, reader io.Reader) (Response, error) {

	if u == "" {
		return Response{}, errors.New("url is  null")
	}
	if method == "" {
		return Response{}, errors.New("method is  null")
	}

	req, _ := http.NewRequest(method, u, reader)

	client := http.Client{
		Timeout: time.Duration(r.timeout) * time.Second, // 超时加在这里，是每次调用的超时
		// CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// 	return http.ErrUseLastResponse // 禁止自动跟随重定向
		// },
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.Header.Set("Cookie", "")
			req.Header.Set("Referer", "")
			req.Header.Set("Content-Type", "")
			return nil
		},
	}
	if len(r.header) != 0 {
		for k, v := range r.header {
			req.Header.Add(k, v)
		}
	}
	if r.basicAuth.Username != "" && r.basicAuth.Password != "" {
		req.SetBasicAuth(r.basicAuth.Username, r.basicAuth.Password)
	}
	res, err := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return Response{}, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return Response{}, err
	}
	resp := Response{
		Body:   resBody,
		Status: res.StatusCode,
		Header: res.Header,
		Cookie: res.Cookies(),
	}
	return resp, nil
}

func (req *Request) ToCurl(method, u string, reader io.Reader) string {
	var curlCmd strings.Builder

	// 1. 添加基础命令和方法
	curlCmd.WriteString("curl -X ")
	curlCmd.WriteString(method)

	// 2. 处理 URL（包含查询参数）
	curlCmd.WriteString(" '")
	curlCmd.WriteString(u)
	curlCmd.WriteString("' \\\n")

	// 3. 处理请求头
	for key, value := range req.header {
		// 转义单引号防止命令中断（' -> '\'）
		escapedValue := strings.ReplaceAll(value, "'", `'\''`)
		fmt.Fprintf(&curlCmd, "-H '%s: %s' \\\n", key, escapedValue)
	}

	// 4. 处理请求体
	if reader != nil {
		// 复制原始 Body（避免读取后丢失）
		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			return fmt.Sprintf("读取请求体失败: %w", err)
		}
		// defer func() {
		// 	// 重置 Body 以便后续使用
		// 	reader = io.NopCloser(bytes.NewBuffer(bodyBytes))
		// }()

		// 特殊处理：空请求体
		if len(bodyBytes) == 0 {
			curlCmd.WriteString("-d ''")
			return curlCmd.String()
		}

		// 转义单引号和换行符
		escapedBody := strings.ReplaceAll(string(bodyBytes), "'", `'\''`)
		escapedBody = strings.ReplaceAll(escapedBody, "\n", `\n`)

		// 根据内容类型决定格式化方式
		contentType := req.header["Content-Type"]
		if strings.Contains(contentType, CONTENT_TYPE_FORM) {
			fmt.Fprintf(&curlCmd, "-d '%s'", escapedBody)
		} else {
			fmt.Fprintf(&curlCmd, "-d '%s'", escapedBody)
		}
	}
	return curlCmd.String()
}

type BasicAuth struct {
	Username string
	Password string
}

type Request struct {
	header    map[string]string
	timeout   int
	basicAuth BasicAuth
}

type Response struct {
	Body   []byte
	Status int
	Header http.Header
	Cookie []*http.Cookie
}

func NewRequest() *Request {
	return &Request{}
}
func (r *Request) Header(header map[string]string) *Request {
	r.header = header
	return r
}

func (r *Request) BasicAuth(username, passowrd string) *Request {
	r.basicAuth = BasicAuth{Username: username, Password: passowrd}
	return r
}

func (r *Request) Timeout(timeout int) *Request {
	r.timeout = timeout
	return r
}
func (r *Request) Get(u string) (Response, error) {
	return r.Do("GET", u, nil)
}
func (r *Request) PostJson(u string, json string) (Response, error) {
	if r.header == nil {
		r.header = make(map[string]string)
	}
	r.header["content-type"] = CONTENT_TYPE_JSON

	return r.Do("POST", u, strings.NewReader(json))
}
func (r *Request) PostForm(u string, data url.Values) (Response, error) {
	if r.header == nil {
		r.header = make(map[string]string)
	}
	if data == nil {
		data = url.Values{}
	}
	r.header["Content-Type"] = CONTENT_TYPE_FORM

	return r.Do("POST", u, bytes.NewBufferString(data.Encode()))
}

// 二进制
func (r *Request) PostStream(u string, bin []byte) (Response, error) {
	if r.header == nil {
		r.header = make(map[string]string)
	}

	r.header["Content-Type"] = CONTENT_TYPE_STREAM

	return r.Do("POST", u, bytes.NewReader(bin))
}

// 二进制
func (r *Request) PutStream(u string, bin []byte) (Response, error) {
	if r.header == nil {
		r.header = make(map[string]string)
	}

	r.header["Content-Type"] = CONTENT_TYPE_STREAM

	return r.Do("PUT", u, bytes.NewReader(bin))
}

func (r *Request) PostData(u string, fileName string, fileHeader *multipart.FileHeader, data map[string]string) (Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 2. 添加文件字段
	file, err := fileHeader.Open()
	if err != nil {
		return Response{}, err
	}
	defer file.Close()

	// 创建表单文件字段（字段名可自定义，如 "file"）
	part, err := writer.CreateFormFile(fileName, fileHeader.Filename)
	if err != nil {
		return Response{}, err
	}

	// 将文件内容写入表单字段
	if _, err := io.Copy(part, file); err != nil {
		return Response{}, err
	}

	// 3. 添加其他字段
	for key, value := range data {
		if err := writer.WriteField(key, value); err != nil {
			return Response{}, err
		}
	}

	// 4. 关闭写入器，完成表单的构建
	if err := writer.Close(); err != nil {
		return Response{}, err
	}
	if r.header == nil {
		r.header = make(map[string]string)
	}
	r.header["Content-Type"] = writer.FormDataContentType()

	return r.Do("POST", u, strings.NewReader(body.String()))
}

func (r *Request) Delete(u string) (Response, error) {
	return r.Do("DELETE", u, nil)
}
func (r *Request) Download(url string, file string) error {
	// 创建一个文件用于保存
	filePath, _ := filepath.Split(file)
	if err := nfile.CreateDir(filePath); err != nil {
		return errors.New("create filePath error")
	}
	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()

	// resp, err := GetConfigurable(url, nil)
	resp, err := r.Get(url)
	if err != nil {
		return err
	}

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, bytes.NewReader(resp.Body))
	if err != nil {
		return err
	}
	return nil
}
