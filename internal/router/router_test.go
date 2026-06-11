package router

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"GoStudy/internal/auth"
	"GoStudy/internal/repository"
	"GoStudy/internal/service"

	"github.com/gin-gonic/gin"
)

func TestEnterpriseRoutesAuthFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store := repository.NewMemoryStore()
	r := New(Dependencies{
		AppName:      "test",
		Env:          "test",
		TokenManager: auth.NewManager("test-secret", time.Hour),
		UserService:  service.NewUserService(store),
		TaskService:  service.NewTaskService(store, store),
	})

	rec := performRequest(r.Engine(), http.MethodGet, "/health", "", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected health 200, got %d", rec.Code)
	}
	if rec.Header().Get("X-Request-ID") == "" {
		t.Fatal("expected X-Request-ID response header")
	}

	registerBody := `{"name":"Alice","email":"alice@example.com","password":"secret1"}`
	rec = performRequest(r.Engine(), http.MethodPost, "/api/v1/users", registerBody, "")
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected register 201, got %d: %s", rec.Code, rec.Body.String())
	}

	rec = performRequest(r.Engine(), http.MethodPost, "/api/v1/tasks", `{"title":"需要鉴权","owner_id":1}`, "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized task create, got %d", rec.Code)
	}

	loginBody := `{"email":"alice@example.com","password":"secret1"}`
	rec = performRequest(r.Engine(), http.MethodPost, "/api/v1/auth/login", loginBody, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected login 200, got %d: %s", rec.Code, rec.Body.String())
	}

	token := extractToken(rec.Body.String())
	if token == "" {
		t.Fatalf("expected token in response: %s", rec.Body.String())
	}

	rec = performRequest(r.Engine(), http.MethodPost, "/api/v1/tasks", `{"title":"发布企业骨架","owner_id":1}`, token)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected task create 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func performRequest(engine http.Handler, method string, path string, body string, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	return rec
}

func extractToken(body string) string {
	marker := `"token":"`
	start := strings.Index(body, marker)
	if start == -1 {
		return ""
	}
	start += len(marker)
	end := strings.Index(body[start:], `"`)
	if end == -1 {
		return ""
	}
	return body[start : start+end]
}
