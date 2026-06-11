package repository

import "testing"

func TestNewStoreMemory(t *testing.T) {
	store, err := NewStore(StorageMemory, "")
	if err != nil {
		t.Fatalf("new memory store: %v", err)
	}
	if store == nil {
		t.Fatal("expected memory store")
	}
}

func TestNewStoreUnknown(t *testing.T) {
	_, err := NewStore("unknown", "")
	if err == nil {
		t.Fatal("expected unsupported storage error")
	}
}
