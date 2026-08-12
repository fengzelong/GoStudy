package cache

import (
	"context"
	"testing"
	"time"
)

func TestMemoryStore(t *testing.T) {
	ctx := context.Background()
	store := NewMemory()
	if err := store.Set(ctx, "user:1", []byte("Alice"), time.Hour); err != nil {
		t.Fatalf("set cache: %v", err)
	}
	value, ok, err := store.Get(ctx, "user:1")
	if err != nil || !ok || string(value) != "Alice" {
		t.Fatalf("unexpected cache value: %q, %t, %v", value, ok, err)
	}
	if err := store.Delete(ctx, "user:1"); err != nil {
		t.Fatalf("delete cache: %v", err)
	}
	if _, ok, err := store.Get(ctx, "user:1"); err != nil || ok {
		t.Fatalf("expected deleted cache key, got %t, %v", ok, err)
	}
}

func TestNoopStoreAndInvalidMode(t *testing.T) {
	ctx := context.Background()
	store, err := New(ModeOff, "", "")
	if err != nil {
		t.Fatalf("new noop cache: %v", err)
	}
	if store.Health(ctx) != "skipped" {
		t.Fatalf("expected skipped health, got %s", store.Health(ctx))
	}
	if err := store.Set(ctx, "key", []byte("value"), time.Hour); err != nil {
		t.Fatalf("set noop cache: %v", err)
	}
	if _, ok, err := store.Get(ctx, "key"); err != nil || ok {
		t.Fatalf("expected noop cache miss, got %t, %v", ok, err)
	}
	if _, err := New("unknown", "", ""); err == nil {
		t.Fatal("expected unsupported cache error")
	}
	if _, err := New(ModeRedis, "", ""); err == nil {
		t.Fatal("expected missing Redis address error")
	}
}
