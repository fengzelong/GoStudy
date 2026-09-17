package cache

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestRedisStoreIntegration(t *testing.T) {
	if os.Getenv("APP_INTEGRATION_REDIS") != "1" {
		t.Skip("未设置 APP_INTEGRATION_REDIS=1，跳过 Redis 集成测试")
	}
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		t.Skip("未设置 REDIS_ADDR，跳过 Redis 集成测试")
	}

	store, err := NewRedis(addr, os.Getenv("REDIS_PASSWORD"))
	if err != nil {
		t.Fatalf("创建 Redis 缓存失败: %v", err)
	}
	defer store.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if health := store.Health(ctx); health != "ok" {
		t.Fatalf("Redis 健康检查失败: %s", health)
	}

	key := fmt.Sprintf("gostudy:integration:cache:%d", time.Now().UnixNano())
	defer store.Delete(context.Background(), key)
	if err := store.Set(ctx, key, []byte("integration-value"), time.Minute); err != nil {
		t.Fatalf("写入 Redis 缓存失败: %v", err)
	}
	value, ok, err := store.Get(ctx, key)
	if err != nil || !ok || string(value) != "integration-value" {
		t.Fatalf("读取 Redis 缓存失败: value=%q ok=%t err=%v", value, ok, err)
	}
	if err := store.Delete(ctx, key); err != nil {
		t.Fatalf("删除 Redis 缓存失败: %v", err)
	}
	if _, ok, err := store.Get(ctx, key); err != nil || ok {
		t.Fatalf("期望 Redis 缓存已删除: ok=%t err=%v", ok, err)
	}
}
