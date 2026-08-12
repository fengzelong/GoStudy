package cache

import (
	"context"
	"sync"
	"time"
)

type entry struct {
	value     []byte
	expiresAt time.Time
}

type memoryStore struct {
	mu      sync.RWMutex
	entries map[string]entry
}

// NewMemory 创建适合本地学习和测试的进程内缓存。
func NewMemory() Store {
	return &memoryStore{entries: make(map[string]entry)}
}

func (s *memoryStore) Get(ctx context.Context, key string) ([]byte, bool, error) {
	s.mu.RLock()
	item, ok := s.entries[key]
	s.mu.RUnlock()
	if !ok || (!item.expiresAt.IsZero() && time.Now().After(item.expiresAt)) {
		if ok {
			s.mu.Lock()
			delete(s.entries, key)
			s.mu.Unlock()
		}
		return nil, false, nil
	}
	return append([]byte(nil), item.value...), true, nil
}

func (s *memoryStore) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	item := entry{value: append([]byte(nil), value...)}
	if ttl > 0 {
		item.expiresAt = time.Now().Add(ttl)
	}
	s.mu.Lock()
	s.entries[key] = item
	s.mu.Unlock()
	return nil
}

func (s *memoryStore) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	delete(s.entries, key)
	s.mu.Unlock()
	return nil
}

func (s *memoryStore) Health(ctx context.Context) string { return "ok" }
func (s *memoryStore) Close() error                      { return nil }

type noopStore struct{}

// NewNoop 创建关闭状态的缓存实现。
func NewNoop() Store { return noopStore{} }

func (noopStore) Get(ctx context.Context, key string) ([]byte, bool, error) { return nil, false, nil }
func (noopStore) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return nil
}
func (noopStore) Delete(ctx context.Context, key string) error { return nil }
func (noopStore) Health(ctx context.Context) string            { return "skipped" }
func (noopStore) Close() error                                 { return nil }
