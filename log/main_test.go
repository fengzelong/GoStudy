package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLoggerComponents(t *testing.T) {
	if getEncoder() == nil {
		t.Fatal("getEncoder() 不应返回 nil")
	}
	if getConsoleEncoder() == nil {
		t.Fatal("getConsoleEncoder() 不应返回 nil")
	}
	if getWriteSyncer() == nil {
		t.Fatal("getWriteSyncer() 不应返回 nil")
	}
	if getLevelEnabler() != zapcore.InfoLevel {
		t.Fatalf("日志级别 = %v，期望 InfoLevel", getLevelEnabler())
	}
}

func TestGinLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	zap.ReplaceGlobals(zap.NewNop())

	r := gin.New()
	r.Use(GinLogger())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("状态码 = %d，期望 %d", w.Code, http.StatusOK)
	}
}
