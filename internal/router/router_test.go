package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"GoStudy/internal/auth"
	"GoStudy/internal/repository"
	"GoStudy/internal/response"
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

	rec = performRequest(r.Engine(), http.MethodPost, "/api/v1/tasks", `{"title":"发布企业骨架"}`, token)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected task create 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHealthDependencies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := repository.NewMemoryStore()
	r := New(Dependencies{
		AppName:      "test",
		Env:          "test",
		TokenManager: auth.NewManager("test-secret", time.Hour),
		UserService:  service.NewUserService(store),
		TaskService:  service.NewTaskService(store, store),
		Health: func() map[string]string {
			return map[string]string{"storage": "skipped", "cache": "ok", "mq": "skipped"}
		},
	})
	rec := performRequest(r.Engine(), http.MethodGet, "/health", "", "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"cache":"ok"`) {
		t.Fatalf("unexpected dependency health response: %d: %s", rec.Code, rec.Body.String())
	}
	assertResponseCode(t, rec.Body.String(), response.CodeOK)
}

func TestEnterpriseRoutesCurrentUserOwnershipAndPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store := repository.NewMemoryStore()
	r := New(Dependencies{
		AppName:      "test",
		Env:          "test",
		TokenManager: auth.NewManager("test-secret", time.Hour),
		UserService:  service.NewUserService(store),
		TaskService:  service.NewTaskService(store, store),
	})

	aliceToken := registerAndLogin(t, r.Engine(), "Alice", "alice@example.com")
	bobToken := registerAndLogin(t, r.Engine(), "Bob", "bob@example.com")

	rec := performRequest(r.Engine(), http.MethodGet, "/api/v1/me", "", "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthenticated me request to fail, got %d", rec.Code)
	}
	assertResponseCode(t, rec.Body.String(), 40101)

	rec = performRequest(r.Engine(), http.MethodGet, "/api/v1/me", "", aliceToken)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"email":"alice@example.com"`) {
		t.Fatalf("expected Alice profile, got %d: %s", rec.Code, rec.Body.String())
	}

	for _, title := range []string{"A1", "A2", "A3"} {
		rec = performRequest(r.Engine(), http.MethodPost, "/api/v1/tasks", `{"title":"`+title+`"}`, aliceToken)
		if rec.Code != http.StatusCreated {
			t.Fatalf("create Alice task %s: %d: %s", title, rec.Code, rec.Body.String())
		}
	}
	rec = performRequest(r.Engine(), http.MethodPost, "/api/v1/tasks", `{"title":"B1"}`, bobToken)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create Bob task: %d: %s", rec.Code, rec.Body.String())
	}
	bobTaskID := extractTaskID(t, rec.Body.String())

	rec = performRequest(r.Engine(), http.MethodGet, "/api/v1/tasks?page=2&page_size=2", "", aliceToken)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"total":3`) || !strings.Contains(rec.Body.String(), `"title":"A3"`) {
		t.Fatalf("expected paginated Alice tasks, got %d: %s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), `"title":"B1"`) {
		t.Fatalf("task list leaked Bob task: %s", rec.Body.String())
	}

	rec = performRequest(r.Engine(), http.MethodPatch, "/api/v1/tasks/"+strconv.FormatInt(bobTaskID, 10)+"/complete", "", aliceToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected cross-user completion to be forbidden, got %d: %s", rec.Code, rec.Body.String())
	}
	assertResponseCode(t, rec.Body.String(), 40301)

	rec = performRequest(r.Engine(), http.MethodGet, "/api/v1/users?page=0", "", aliceToken)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid pagination to fail, got %d: %s", rec.Code, rec.Body.String())
	}
	assertResponseCode(t, rec.Body.String(), 40001)
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

func registerAndLogin(t *testing.T, engine http.Handler, name string, email string) string {
	t.Helper()
	registerBody := `{"name":"` + name + `","email":"` + email + `","password":"secret1"}`
	rec := performRequest(engine, http.MethodPost, "/api/v1/users", registerBody, "")
	if rec.Code != http.StatusCreated {
		t.Fatalf("register %s: %d: %s", email, rec.Code, rec.Body.String())
	}
	loginBody := `{"email":"` + email + `","password":"secret1"}`
	rec = performRequest(engine, http.MethodPost, "/api/v1/auth/login", loginBody, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("login %s: %d: %s", email, rec.Code, rec.Body.String())
	}
	token := extractToken(rec.Body.String())
	if token == "" {
		t.Fatalf("missing token for %s", email)
	}
	return token
}

func extractTaskID(t *testing.T, body string) int64 {
	t.Helper()
	var payload struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &payload); err != nil || payload.Data.ID == 0 {
		t.Fatalf("extract task ID from %s: %v", body, err)
	}
	return payload.Data.ID
}

func assertResponseCode(t *testing.T, body string, want int) {
	t.Helper()
	var payload struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Code != want {
		t.Fatalf("expected response code %d, got %d: %s", want, payload.Code, body)
	}
}
