package main

import (
	"testing"
)

func TestNewArithRPCServer(t *testing.T) {
	server, err := newArithRPCServer()
	if err != nil {
		t.Fatalf("创建 RPC server 失败: %v", err)
	}
	if server == nil {
		t.Fatal("server 不应为 nil")
	}
}
