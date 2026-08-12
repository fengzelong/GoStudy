package cache

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const (
	ModeOff    = "off"
	ModeMemory = "memory"
	ModeRedis  = "redis"
)

// Store 为服务层提供小型缓存抽象，避免业务代码依赖具体缓存客户端。
type Store interface {
	Get(ctx context.Context, key string) ([]byte, bool, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Health(ctx context.Context) string
	Close() error
}

// New 根据配置创建缓存。off 和 memory 不依赖外部服务。
func New(mode string, addr string, password string) (Store, error) {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "", ModeOff:
		return NewNoop(), nil
	case ModeMemory:
		return NewMemory(), nil
	case ModeRedis:
		if strings.TrimSpace(addr) == "" {
			return nil, fmt.Errorf("redis address is required when APP_CACHE=redis")
		}
		return NewRedis(addr, password)
	default:
		return nil, fmt.Errorf("unsupported cache: %s", mode)
	}
}
