package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/golang", nil)
	w := httptest.NewRecorder()

	home(w, req)

	if w.Body.String() != "Hi there, I love golang!" {
		t.Fatalf("响应 = %q", w.Body.String())
	}
}

func TestUserHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/jack", nil)
	w := httptest.NewRecorder()

	user(w, req)

	if w.Body.String() != "Hi there, user jack!" {
		t.Fatalf("响应 = %q", w.Body.String())
	}
}

func TestNewMuxRoutes(t *testing.T) {
	mux := newMux()

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/user/tom", nil))
	if w.Body.String() != "Hi there, user tom!" {
		t.Fatalf("/user/tom 响应 = %q", w.Body.String())
	}
}
