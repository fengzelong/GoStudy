package main

import (
	"os"
	"testing"
)

func TestEnvUsesDefaultAndOverride(t *testing.T) {
	key := "WEBSOCKET_TEST_ADDR"
	oldValue, hadValue := os.LookupEnv(key)
	defer restoreEnv(key, oldValue, hadValue)

	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("清理环境变量失败: %v", err)
	}
	if got := env(key, "127.0.0.1:8080"); got != "127.0.0.1:8080" {
		t.Fatalf("env 默认值 = %q", got)
	}

	if err := os.Setenv(key, "127.0.0.1:9090"); err != nil {
		t.Fatalf("设置环境变量失败: %v", err)
	}
	if got := env(key, "127.0.0.1:8080"); got != "127.0.0.1:9090" {
		t.Fatalf("env 环境变量值 = %q", got)
	}
}

func restoreEnv(key string, oldValue string, hadValue bool) {
	if hadValue {
		_ = os.Setenv(key, oldValue)
		return
	}
	_ = os.Unsetenv(key)
}
