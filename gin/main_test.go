package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestRouter(t *testing.T) *gin.Engine {
	t.Helper()

	uploadDir := t.TempDir()
	setenv(t, "GIN_UPLOAD_DIR", uploadDir)
	setenv(t, "GIN_LOG_FILE", "-")
	gin.SetMode(gin.TestMode)

	return setupRouter()
}

func setenv(t *testing.T, key, value string) {
	t.Helper()

	old, ok := os.LookupEnv(key)
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("设置环境变量失败: %v", err)
	}
	t.Cleanup(func() {
		if ok {
			_ = os.Setenv(key, old)
			return
		}
		_ = os.Unsetenv(key)
	})
}

func performRequest(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestPingAndQuery(t *testing.T) {
	router := newTestRouter(t)

	w := performRequest(router, http.MethodGet, "/ping", "")
	if w.Code != http.StatusOK {
		t.Fatalf("/ping 状态码 = %d，期望 %d", w.Code, http.StatusOK)
	}
	if !strings.Contains(w.Body.String(), "pong") {
		t.Fatalf("/ping 响应 = %s，期望包含 pong", w.Body.String())
	}

	w = performRequest(router, http.MethodGet, "/query?name=tom&age=18", "")
	if w.Code != http.StatusOK {
		t.Fatalf("/query 状态码 = %d，期望 %d", w.Code, http.StatusOK)
	}
	if !strings.Contains(w.Body.String(), "tom") {
		t.Fatalf("/query 响应 = %s，期望包含 tom", w.Body.String())
	}
}

func TestBindRoutes(t *testing.T) {
	router := newTestRouter(t)

	body := `{"user":"jack","password":"123"}`
	w := performRequest(router, http.MethodPost, "/modelBind", body)
	if w.Code != http.StatusOK {
		t.Fatalf("/modelBind 状态码 = %d，期望 %d，响应 %s", w.Code, http.StatusOK, w.Body.String())
	}

	w = performRequest(router, http.MethodGet, "/urlBind/tom/550e8400-e29b-41d4-a716-446655440000", "")
	if w.Code != http.StatusOK {
		t.Fatalf("/urlBind 状态码 = %d，期望 %d，响应 %s", w.Code, http.StatusOK, w.Body.String())
	}

	w = performRequest(router, http.MethodGet, "/bindQuery?keyword=gin&page=2&page_size=20", "")
	if w.Code != http.StatusOK {
		t.Fatalf("/bindQuery 状态码 = %d，期望 %d，响应 %s", w.Code, http.StatusOK, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"page":2`) {
		t.Fatalf("/bindQuery 响应 = %s，期望包含 page=2", w.Body.String())
	}
}

func TestHeaderFormCookieAndAuth(t *testing.T) {
	router := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/bindHeader", nil)
	req.Header.Set("X-Request-Id", "req-1")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("/bindHeader 状态码 = %d，期望 %d，响应 %s", w.Code, http.StatusOK, w.Body.String())
	}

	req = httptest.NewRequest(http.MethodPost, "/form", strings.NewReader("name=tom&message=hello"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("/form 状态码 = %d，期望 %d，响应 %s", w.Code, http.StatusOK, w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/cookie/read", nil)
	req.AddCookie(&http.Cookie{Name: "gin_user", Value: "jack"})
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("/cookie/read 状态码 = %d，期望 %d，响应 %s", w.Code, http.StatusOK, w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/basic/secret", nil)
	req.SetBasicAuth("gin", "123456")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("/basic/secret 状态码 = %d，期望 %d，响应 %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestUploadRoutes(t *testing.T) {
	router := newTestRouter(t)

	w := performRequest(router, http.MethodPost, "/upload", "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("缺少文件时状态码 = %d，期望 %d", w.Code, http.StatusBadRequest)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "hello.txt")
	if err != nil {
		t.Fatalf("创建上传字段失败: %v", err)
	}
	if _, err := io.WriteString(part, "hello gin"); err != nil {
		t.Fatalf("写入上传内容失败: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("关闭 multipart writer 失败: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("上传文件状态码 = %d，期望 %d，响应 %s", w.Code, http.StatusOK, w.Body.String())
	}

	if _, err := os.Stat(filepath.Join(uploadDir(), "hello.txt")); err != nil {
		t.Fatalf("上传文件未保存: %v", err)
	}
}
