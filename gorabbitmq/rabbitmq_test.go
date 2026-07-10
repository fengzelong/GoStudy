package gorabbitmq

import (
	"os"
	"testing"

	"github.com/streadway/amqp"
)

func TestNewRabbitMQReadsEnv(t *testing.T) {
	oldValue, hadValue := os.LookupEnv("RABBITMQ_URL")
	if err := os.Setenv("RABBITMQ_URL", "amqp://example/"); err != nil {
		t.Fatalf("设置 RABBITMQ_URL 失败: %v", err)
	}
	defer restoreEnv("RABBITMQ_URL", oldValue, hadValue)

	mq := NewRabbitMQ("queue", "exchange", "key")
	if mq.QueueName != "queue" || mq.Exchange != "exchange" || mq.Key != "key" {
		t.Fatalf("RabbitMQ 字段未正确初始化: %#v", mq)
	}
	if mq.MqUrl != "amqp://example/" {
		t.Fatalf("MqUrl = %q，期望读取环境变量", mq.MqUrl)
	}
}

func TestRabbitMQConnectionWithConfiguredURL(t *testing.T) {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		t.Skip("未设置 RABBITMQ_URL，跳过 RabbitMQ 真实连接测试")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		t.Fatalf("连接 RabbitMQ 失败: %v", err)
	}
	defer conn.Close()
}

func restoreEnv(key string, oldValue string, hadValue bool) {
	if hadValue {
		_ = os.Setenv(key, oldValue)
		return
	}
	_ = os.Unsetenv(key)
}
