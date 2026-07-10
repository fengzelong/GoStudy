package main

import (
	"os"
	"testing"
)

func TestRedisPingWithConfiguredAddr(t *testing.T) {
	if os.Getenv("REDIS_ADDR") == "" {
		t.Skip("未设置 REDIS_ADDR，跳过 Redis 真实连接测试")
	}

	conn := pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		t.Fatalf("连接 Redis 失败: %v", err)
	}
	if _, err := conn.Do("PING"); err != nil {
		t.Fatalf("Redis PING 失败: %v", err)
	}
}
