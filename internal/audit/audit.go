package audit

import (
	"context"
	"sync"
	"time"
)

// Entry 是关键业务操作的审计记录。
type Entry struct {
	Action    string    `json:"action"`
	ActorID   int64     `json:"actor_id"`
	Resource  string    `json:"resource"`
	CreatedAt time.Time `json:"created_at"`
}

// Logger 隔离审计记录的具体去向，当前默认保存在内存中。
type Logger interface {
	Record(ctx context.Context, entry Entry) error
}

// Reader 提供审计记录读取能力，便于后续替换为持久化实现。
type Reader interface {
	List(ctx context.Context) ([]Entry, error)
}

// Store 同时提供审计记录的写入和读取能力。
type Store interface {
	Logger
	Reader
}

type MemoryLogger struct {
	mu      sync.RWMutex
	entries []Entry
}

func NewMemoryLogger() *MemoryLogger { return &MemoryLogger{} }

func (l *MemoryLogger) Record(ctx context.Context, entry Entry) error {
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}
	l.mu.Lock()
	l.entries = append(l.entries, entry)
	l.mu.Unlock()
	return nil
}

// Entries 返回审计记录快照，供测试和后续管理接口使用。
func (l *MemoryLogger) Entries() []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return append([]Entry(nil), l.entries...)
}

// List 返回审计记录快照，满足管理端查询和后续持久化抽象。
func (l *MemoryLogger) List(ctx context.Context) ([]Entry, error) {
	return l.Entries(), nil
}
