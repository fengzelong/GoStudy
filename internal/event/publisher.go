package event

import (
	"context"
	"fmt"
	"strings"
)

const (
	ModeOff      = "off"
	ModeMemory   = "memory"
	ModeRabbitMQ = "rabbitmq"
)

// Publisher 发布任务领域事件，默认实现不执行外部 I/O。
type Publisher interface {
	Publish(ctx context.Context, name string, payload interface{}) error
	Health(ctx context.Context) string
	Close() error
}

// New 根据配置创建事件发布器。
func New(mode string, url string, queue string) (Publisher, error) {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "", ModeOff:
		return noopPublisher{}, nil
	case ModeMemory:
		return NewMemory(), nil
	case ModeRabbitMQ:
		if strings.TrimSpace(url) == "" || strings.TrimSpace(queue) == "" {
			return nil, fmt.Errorf("rabbitmq url and queue are required when APP_MQ=rabbitmq")
		}
		return NewRabbitMQ(url, queue)
	default:
		return nil, fmt.Errorf("unsupported message queue: %s", mode)
	}
}

type noopPublisher struct{}

func (noopPublisher) Publish(ctx context.Context, name string, payload interface{}) error { return nil }
func (noopPublisher) Health(ctx context.Context) string                                   { return "skipped" }
func (noopPublisher) Close() error                                                        { return nil }
