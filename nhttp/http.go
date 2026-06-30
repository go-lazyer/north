package nhttp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-lazyer/north/nfile"
)

const (
	CONTENT_TYPE_STREAM = "application/octet-stream"
	CONTENT_TYPE_FORM   = "application/x-www-form-urlencoded"
	CONTENT_TYPE_JSON   = "application/json"
	CONTENT_TYPE_NDJSON = "application/x-ndjson"
	CONTENT_TYPE_TEXT   = "text/plain"
	CONTENT_TYPE_DATA   = "multipart/form-data"
)

var (
	defaultClient = sync.OnceValue(func() *http.Client {
		return &http.Client{
			Timeout: 30 * time.Second, // 默认全局超时
		}
	})
)

type BasicAuth struct {
	Username string
	Password string
}

type Request struct {
	header        map[string]string
	timeout       int
	basicAuth     BasicAuth
	client        *http.Client
	checkRedirect func(req *http.Request, via []*http.Request) error
}

type Response struct {
	Body   []byte
	Status int
	Header http.Header
	Cookie []*http.Cookie
}

func NewRequest() *Request {
	return &Request{client: httpClient()}
}

func httpClient() *http.Client {
	return defaultClient()
}

func (r *Request) Do(method, u string, reader io.Reader) (Response, error) {

	if u == "" {
		return Response{}, errors.New("url is null")
	}
	if method == "" {
		return Response{}, errors.New("method is null")
	}

	req, err := http.NewRequest(method, u, reader)
	if err != nil {
		return Response{}, err
	}

	if r.timeout > 0 {
		ctx, cancel := context.WithTimeout(req.Context(), time.Duration(r.timeout)*time.Second)
		defer cancel()
		req = req.WithContext(ctx)
	}

	if r.client == nil {
		r.client = httpClient()
	}
	if len(r.header) != 0 {
		for k, v := range r.header {
			req.Header.Set(k, v)
		}
	}
	if r.basicAuth.Username != "" && r.basicAuth.Password != "" {
		req.SetBasicAuth(r.basicAuth.Username, r.basicAuth.Password)
	}

	// 单次请求可用的自定义 CheckRedirect，不影响全局 client
	c := r.client
	if r.checkRedirect != nil {
		cc := *r.client // 浅拷贝
		cc.CheckRedirect = r.checkRedirect
		c = &cc
	}
	res, err := c.Do(req)
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
	// log.Printf("HTTP %s %s\nStatus: %d\nHeaders: %v\nBody: %s\n", method, u, res.StatusCode, res.Header, string(resBody))
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
			return fmt.Sprintf("读取请求体失败: %v", err)
		}
		if seeker, ok := reader.(io.Seeker); ok {
			_, _ = seeker.Seek(0, io.SeekStart)
		}

		// 特殊处理：空请求体
		if len(bodyBytes) == 0 {
			curlCmd.WriteString("-d ''")
			return curlCmd.String()
		}

		// 转义单引号和换行符
		escapedBody := strings.ReplaceAll(string(bodyBytes), "'", `'\''`)
		escapedBody = strings.ReplaceAll(escapedBody, "\n", `\n`)

		// 根据内容类型决定格式化方式
		contentType := req.header["content-type"]
		if contentType == "" {
			contentType = req.header["Content-Type"]
		}
		if strings.Contains(contentType, CONTENT_TYPE_FORM) {
			fmt.Fprintf(&curlCmd, "-d '%s'", escapedBody)
		} else {
			fmt.Fprintf(&curlCmd, "-d '%s'", escapedBody)
		}
	}
	return curlCmd.String()
}

func (r *Request) Header(header map[string]string) *Request {
	r.header = make(map[string]string, len(header))
	for k, v := range header {
		r.header[k] = v
	}
	return r
}

func (r *Request) BasicAuth(username, password string) *Request {
	r.basicAuth = BasicAuth{Username: username, Password: password}
	return r
}

// CheckRedirect 为当前请求设置自定义重定向策略，不影响全局 client
func (r *Request) CheckRedirect(fn func(req *http.Request, via []*http.Request) error) *Request {
	r.checkRedirect = fn
	return r
}

func (r *Request) Timeout(timeout int) *Request {
	r.timeout = timeout
	return r
}
func (r *Request) Get(u string) (Response, error) {
	return r.Do("GET", u, nil)
}
func (r *Request) PatchJson(u string, json string) (Response, error) {
	if r.header == nil {
		r.header = make(map[string]string)
	}
	r.header["content-type"] = CONTENT_TYPE_JSON

	return r.Do("PATCH", u, strings.NewReader(json))
}

func (r *Request) PostJson(u string, json string) (Response, error) {
	if r.header == nil {
		r.header = make(map[string]string)
	}
	r.header["content-type"] = CONTENT_TYPE_JSON

	return r.Do("POST", u, strings.NewReader(json))
}

// 换行分隔的 JSON（Newline Delimited JSON），每一行是一个独立 JSON，专门用于流式传输、大数据、日志、事件推送
func (r *Request) PostNdjson(u string, json string) (Response, error) {
	if r.header == nil {
		r.header = make(map[string]string)
	}
	r.header["content-type"] = CONTENT_TYPE_NDJSON

	return r.Do("POST", u, strings.NewReader(json))
}
func (r *Request) PostText(u string, json string) (Response, error) {
	if r.header == nil {
		r.header = make(map[string]string)
	}
	r.header["content-type"] = CONTENT_TYPE_TEXT

	return r.Do("POST", u, strings.NewReader(json))
}
func (r *Request) PostForm(u string, data url.Values) (Response, error) {
	if r.header == nil {
		r.header = make(map[string]string)
	}
	if data == nil {
		data = url.Values{}
	}
	r.header["content-type"] = CONTENT_TYPE_FORM

	return r.Do("POST", u, bytes.NewBufferString(data.Encode()))
}

// 二进制
func (r *Request) PostStream(u string, bin []byte) (Response, error) {
	if r.header == nil {
		r.header = make(map[string]string)
	}

	r.header["content-type"] = CONTENT_TYPE_STREAM

	return r.Do("POST", u, bytes.NewReader(bin))
}

// 二进制
func (r *Request) PutStream(u string, bin []byte) (Response, error) {
	if r.header == nil {
		r.header = make(map[string]string)
	}

	r.header["content-type"] = CONTENT_TYPE_STREAM

	return r.Do("PUT", u, bytes.NewReader(bin))
}
func (r *Request) PutJson(u string, json string) (Response, error) {
	if r.header == nil {
		r.header = make(map[string]string)
	}
	r.header["content-type"] = CONTENT_TYPE_JSON

	return r.Do("PUT", u, strings.NewReader(json))
}

func (r *Request) PostData(u string, fileName string, fileHeader *multipart.FileHeader, data map[string]string) (Response, error) {
	if fileHeader == nil {
		return Response{}, errors.New("fileHeader is nil")
	}
	if fileName == "" {
		fileName = "file"
	}

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
	r.header["content-type"] = writer.FormDataContentType()

	return r.Do("POST", u, bytes.NewReader(body.Bytes()))
}

func (r *Request) Delete(u string) (Response, error) {
	return r.Do("DELETE", u, nil)
}
func (r *Request) Download(rawURL string, file string) error {
	// 创建一个文件用于保存
	filePath, _ := filepath.Split(file)
	if filePath != "" {
		if err := nfile.CreateDir(filePath); err != nil {
			return errors.New("create filePath error")
		}
	}
	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := r.Get(rawURL)
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
