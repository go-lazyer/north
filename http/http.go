package nhttp

import (
	"bytes"
	"errors"
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

func (r *Request) Do(method, u string, reader io.Reader) ([]byte, int, error) {

	if u == "" {
		return nil, 0, errors.New("url is  null")
	}
	if method == "" {
		return nil, 0, errors.New("method is  null")
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
		return nil, 0, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}
	return resBody, res.StatusCode, nil
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
func (r *Request) Get(u string) ([]byte, int, error) {
	return r.Do("GET", u, nil)
}
func (r *Request) PostJson(u string, json string) ([]byte, int, error) {
	if r.header == nil {
		r.header = make(map[string]string)
	}
	r.header["content-type"] = CONTENT_TYPE_JSON

	return r.Do("POST", u, strings.NewReader(json))
}
func (r *Request) PostForm(u string, data url.Values) ([]byte, int, error) {
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
func (r *Request) PostStream(u string, bin []byte) ([]byte, int, error) {
	if r.header == nil {
		r.header = make(map[string]string)
	}

	r.header["Content-Type"] = CONTENT_TYPE_STREAM

	return r.Do("POST", u, bytes.NewReader(bin))
}

// 二进制
func (r *Request) PutStream(u string, bin []byte) ([]byte, int, error) {
	if r.header == nil {
		r.header = make(map[string]string)
	}

	r.header["Content-Type"] = CONTENT_TYPE_STREAM

	return r.Do("PUT", u, bytes.NewReader(bin))
}

func (r *Request) PostData(u string, fileName string, fileHeader *multipart.FileHeader, data map[string]string) ([]byte, int, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 2. 添加文件字段
	file, err := fileHeader.Open()
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	// 创建表单文件字段（字段名可自定义，如 "file"）
	part, err := writer.CreateFormFile(fileName, fileHeader.Filename)
	if err != nil {
		return nil, 0, err
	}

	// 将文件内容写入表单字段
	if _, err := io.Copy(part, file); err != nil {
		return nil, 0, err
	}

	// 3. 添加其他字段
	for key, value := range data {
		if err := writer.WriteField(key, value); err != nil {
			return nil, 0, err
		}
	}

	// 4. 关闭写入器，完成表单的构建
	if err := writer.Close(); err != nil {
		return nil, 0, err
	}
	if r.header == nil {
		r.header = make(map[string]string)
	}
	r.header["Content-Type"] = writer.FormDataContentType()

	return r.Do("POST", u, strings.NewReader(body.String()))
}

func (r *Request) Delete(u string) ([]byte, int, error) {
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
	resp, _, err := r.Get(url)
	if err != nil {
		return err
	}

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, bytes.NewReader(resp))
	if err != nil {
		return err
	}
	return nil
}
