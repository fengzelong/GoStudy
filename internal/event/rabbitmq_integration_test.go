package event

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/streadway/amqp"
)

func TestRabbitMQPublisherIntegration(t *testing.T) {
	if os.Getenv("APP_INTEGRATION_RABBITMQ") != "1" {
		t.Skip("未设置 APP_INTEGRATION_RABBITMQ=1，跳过 RabbitMQ 集成测试")
	}
	url := os.Getenv("RABBITMQ_URL")
	queuePrefix := os.Getenv("RABBITMQ_QUEUE")
	if url == "" || queuePrefix == "" {
		t.Skip("未设置 RABBITMQ_URL 或 RABBITMQ_QUEUE，跳过 RabbitMQ 集成测试")
	}

	queue := fmt.Sprintf("%s.integration.%d", queuePrefix, time.Now().UnixNano())
	publisher, err := NewRabbitMQ(url, queue)
	if err != nil {
		t.Fatalf("创建 RabbitMQ 发布器失败: %v", err)
	}
	defer publisher.Close()
	if health := publisher.Health(context.Background()); health != "ok" {
		t.Fatalf("RabbitMQ 健康检查失败: %s", health)
	}

	connection, err := amqp.Dial(url)
	if err != nil {
		t.Fatalf("连接 RabbitMQ 消费端失败: %v", err)
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		t.Fatalf("创建 RabbitMQ 消费通道失败: %v", err)
	}
	defer channel.Close()
	defer channel.QueueDelete(queue, false, false, false)

	payload := map[string]int64{"task_id": 1}
	if err := publisher.Publish(context.Background(), "task.created", payload); err != nil {
		t.Fatalf("发布 RabbitMQ 事件失败: %v", err)
	}

	deadline := time.Now().Add(5 * time.Second)
	for {
		delivery, ok, err := channel.Get(queue, false)
		if err != nil {
			t.Fatalf("读取 RabbitMQ 事件失败: %v", err)
		}
		if ok {
			defer delivery.Ack(false)
			var message struct {
				Name    string           `json:"name"`
				Payload map[string]int64 `json:"payload"`
			}
			if err := json.Unmarshal(delivery.Body, &message); err != nil {
				t.Fatalf("解析 RabbitMQ 事件失败: %v", err)
			}
			if message.Name != "task.created" || message.Payload["task_id"] != 1 {
				t.Fatalf("RabbitMQ 事件内容不符合预期: %+v", message)
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("等待 RabbitMQ 事件超时")
		}
		time.Sleep(50 * time.Millisecond)
	}
}
