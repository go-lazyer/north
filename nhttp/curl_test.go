package nhttp

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestToCurl_SingleQuoteEscape(t *testing.T) {
	// 测试 header 中包含单引号的情况
	req := httptest.NewRequest("GET", "http://example.com/api", nil)
	req.Header.Set("X-Special", "it's a test")
	result := RequestToCurl(req)
	if !strings.Contains(result, `'\''`) {
		t.Errorf("期望单引号被转义，实际 %q", result)
	}
}

func TestRequestToCurl_BodyWithSingleQuote(t *testing.T) {
	// 测试 body 中包含单引号
	req := httptest.NewRequest("POST", "http://example.com/api", strings.NewReader(`{"msg":"it's working"}`))
	req.Header.Set("Content-Type", "application/json")
	result := RequestToCurl(req)
	if !strings.Contains(result, `'\''`) {
		t.Errorf("期望 body 中单引号被转义，实际 %q", result)
	}
}

func TestRequestToCurl_EmptyBody(t *testing.T) {
	// 测试空 body
	req := httptest.NewRequest("POST", "http://example.com/api", strings.NewReader(""))
	result := RequestToCurl(req)
	if !strings.Contains(result, "-d ''") {
		t.Errorf("期望包含 -d ''，实际 %q", result)
	}
}
