package event

import (
	"context"
	"sync"
)

// Message 是内存发布器记录的事件，便于测试和本地观察。
type Message struct {
	Name    string
	Payload interface{}
}

type MemoryPublisher struct {
	mu       sync.RWMutex
	messages []Message
}

func NewMemory() *MemoryPublisher { return &MemoryPublisher{} }

func (p *MemoryPublisher) Publish(ctx context.Context, name string, payload interface{}) error {
	p.mu.Lock()
	p.messages = append(p.messages, Message{Name: name, Payload: payload})
	p.mu.Unlock()
	return nil
}

// Messages 返回已发布事件的快照，供测试使用。
func (p *MemoryPublisher) Messages() []Message {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return append([]Message(nil), p.messages...)
}

func (p *MemoryPublisher) Health(ctx context.Context) string { return "ok" }
func (p *MemoryPublisher) Close() error                      { return nil }
