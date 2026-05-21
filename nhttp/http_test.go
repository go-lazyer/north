package nhttp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestToCurl(t *testing.T) {
	// 测试 GET 请求
	req := httptest.NewRequest("GET", "http://example.com/api?foo=bar", nil)
	result := RequestToCurl(req)
	if !strings.Contains(result, "curl -X GET") {
		t.Errorf("期望包含 'curl -X GET'，实际 %q", result)
	}
	if !strings.Contains(result, "http://example.com/api?foo=bar") {
		t.Errorf("期望包含 URL，实际 %q", result)
	}

	// 测试 POST 请求带 Body
	req = httptest.NewRequest("POST", "http://example.com/api", strings.NewReader(`{"key":"value"}`))
	req.Header.Set("Content-Type", "application/json")
	result = RequestToCurl(req)
	if !strings.Contains(result, "curl -X POST") {
		t.Errorf("期望包含 'curl -X POST'，实际 %q", result)
	}
	if !strings.Contains(result, "-H 'Content-Type: application/json'") {
		t.Errorf("期望包含 Content-Type header，实际 %q", result)
	}
	if !strings.Contains(result, "-d") {
		t.Errorf("期望包含 -d 参数，实际 %q", result)
	}

	// 测试带自定义 Header
	req = httptest.NewRequest("GET", "http://example.com/api", nil)
	req.Header.Set("Authorization", "Bearer token123")
	req.Header.Set("X-Custom", "value")
	result = RequestToCurl(req)
	if !strings.Contains(result, "-H 'Authorization: Bearer token123'") {
		t.Errorf("期望包含 Authorization header，实际 %q", result)
	}
	if !strings.Contains(result, "-H 'X-Custom: value'") {
		t.Errorf("期望包含 X-Custom header，实际 %q", result)
	}

	// 测试空 Body
	req = httptest.NewRequest("DELETE", "http://example.com/api/1", nil)
	result = RequestToCurl(req)
	if !strings.Contains(result, "curl -X DELETE") {
		t.Errorf("期望包含 'curl -X DELETE'，实际 %q", result)
	}
}

func TestNewRequest(t *testing.T) {
	r := NewRequest()
	if r == nil {
		t.Fatal("NewRequest 返回 nil")
	}
	if r.client == nil {
		t.Fatal("NewRequest 的 client 应不为 nil")
	}
}

func TestRequest_Header(t *testing.T) {
	r := NewRequest()
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer test",
	}
	result := r.Header(headers)
	if result == nil {
		t.Fatal("Header 返回 nil")
	}
	if len(r.header) != 2 {
		t.Errorf("期望 2 个 header，实际 %d", len(r.header))
	}
	if r.header["Content-Type"] != "application/json" {
		t.Errorf("期望 Content-Type 为 application/json，实际 %q", r.header["Content-Type"])
	}
}

func TestRequest_BasicAuth(t *testing.T) {
	r := NewRequest()
	result := r.BasicAuth("user", "pass")
	if result == nil {
		t.Fatal("BasicAuth 返回 nil")
	}
	if r.basicAuth.Username != "user" {
		t.Errorf("期望 Username 为 'user'，实际 %q", r.basicAuth.Username)
	}
	if r.basicAuth.Password != "pass" {
		t.Errorf("期望 Password 为 'pass'，实际 %q", r.basicAuth.Password)
	}
}

func TestRequest_Timeout(t *testing.T) {
	r := NewRequest()
	result := r.Timeout(10)
	if result == nil {
		t.Fatal("Timeout 返回 nil")
	}
	if r.timeout != 10 {
		t.Errorf("期望 timeout 为 10，实际 %d", r.timeout)
	}
}

func TestRequest_ToCurl(t *testing.T) {
	r := NewRequest()
	r.Header(map[string]string{"Content-Type": "application/json"})
	result := r.ToCurl("POST", "http://example.com/api", strings.NewReader(`{"key":"value"}`))
	if !strings.Contains(result, "curl -X POST") {
		t.Errorf("期望包含 'curl -X POST'，实际 %q", result)
	}
	if !strings.Contains(result, "http://example.com/api") {
		t.Errorf("期望包含 URL，实际 %q", result)
	}
	if !strings.Contains(result, "-H 'Content-Type: application/json'") {
		t.Errorf("期望包含 Content-Type header，实际 %q", result)
	}
	if !strings.Contains(result, "-d") {
		t.Errorf("期望包含 -d 参数，实际 %q", result)
	}

	// 无 Body 的请求
	r2 := NewRequest()
	result2 := r2.ToCurl("GET", "http://example.com/api", nil)
	if !strings.Contains(result2, "curl -X GET") {
		t.Errorf("期望包含 'curl -X GET'，实际 %q", result2)
	}
	if strings.Contains(result2, "-d") {
		t.Errorf("GET 请求不应包含 -d 参数，实际 %q", result2)
	}
}

func TestRequest_Do_InvalidURL(t *testing.T) {
	r := NewRequest()
	// 空 URL
	_, err := r.Do("GET", "", nil)
	if err == nil {
		t.Errorf("空 URL 应返回错误")
	}

	// 空 method
	_, err = r.Do("", "http://example.com", nil)
	if err == nil {
		t.Errorf("空 method 应返回错误")
	}
}

func TestRequest_Do_WithMockServer(t *testing.T) {
	// 创建 mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	r := NewRequest()
	resp, err := r.Get(server.URL)
	if err != nil {
		t.Fatalf("Get 失败: %v", err)
	}
	if resp.Status != 200 {
		t.Errorf("期望状态码 200，实际 %d", resp.Status)
	}
	if string(resp.Body) != "OK" {
		t.Errorf("期望 Body 为 'OK'，实际 %q", string(resp.Body))
	}
}

func TestRequest_PostJson_WithMockServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("期望 POST 方法，实际 %s", r.Method)
		}
		if r.Header.Get("Content-Type") != CONTENT_TYPE_JSON {
			t.Errorf("期望 Content-Type 为 %s，实际 %s", CONTENT_TYPE_JSON, r.Header.Get("Content-Type"))
		}
		w.WriteHeader(200)
		w.Write([]byte("created"))
	}))
	defer server.Close()

	r := NewRequest()
	resp, err := r.PostJson(server.URL, `{"name":"test"}`)
	if err != nil {
		t.Fatalf("PostJson 失败: %v", err)
	}
	if resp.Status != 200 {
		t.Errorf("期望状态码 200，实际 %d", resp.Status)
	}
}

func TestRequest_Delete_WithMockServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("期望 DELETE 方法，实际 %s", r.Method)
		}
		w.WriteHeader(204)
	}))
	defer server.Close()

	r := NewRequest()
	resp, err := r.Delete(server.URL)
	if err != nil {
		t.Fatalf("Delete 失败: %v", err)
	}
	if resp.Status != 204 {
		t.Errorf("期望状态码 204，实际 %d", resp.Status)
	}
}
