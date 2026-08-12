package event

import (
	"context"
	"testing"
)

func TestMemoryPublisher(t *testing.T) {
	publisher := NewMemory()
	if err := publisher.Publish(context.Background(), "task.created", map[string]int64{"id": 1}); err != nil {
		t.Fatalf("publish event: %v", err)
	}
	messages := publisher.Messages()
	if len(messages) != 1 || messages[0].Name != "task.created" {
		t.Fatalf("unexpected messages: %+v", messages)
	}
	if publisher.Health(context.Background()) != "ok" {
		t.Fatal("expected memory publisher health to be ok")
	}
}

func TestNoopPublisherAndInvalidMode(t *testing.T) {
	publisher, err := New(ModeOff, "", "")
	if err != nil {
		t.Fatalf("new noop publisher: %v", err)
	}
	if err := publisher.Publish(context.Background(), "task.created", nil); err != nil {
		t.Fatalf("publish noop event: %v", err)
	}
	if publisher.Health(context.Background()) != "skipped" {
		t.Fatalf("expected skipped health, got %s", publisher.Health(context.Background()))
	}
	if _, err := New("unknown", "", ""); err == nil {
		t.Fatal("expected unsupported publisher error")
	}
	if _, err := New(ModeRabbitMQ, "", ""); err == nil {
		t.Fatal("expected missing RabbitMQ configuration error")
	}
}
