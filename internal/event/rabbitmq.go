package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/streadway/amqp"
)

type rabbitPublisher struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      string
}

// NewRabbitMQ 建立 RabbitMQ 连接并声明任务事件队列。
func NewRabbitMQ(url string, queue string) (Publisher, error) {
	connection, err := amqp.DialConfig(url, amqp.Config{Heartbeat: 10 * time.Second})
	if err != nil {
		return nil, err
	}
	channel, err := connection.Channel()
	if err != nil {
		connection.Close()
		return nil, err
	}
	if _, err := channel.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		channel.Close()
		connection.Close()
		return nil, err
	}
	return &rabbitPublisher{connection: connection, channel: channel, queue: queue}, nil
}

func (p *rabbitPublisher) Publish(ctx context.Context, name string, payload interface{}) error {
	body, err := json.Marshal(struct {
		Name    string      `json:"name"`
		Payload interface{} `json:"payload"`
	}{Name: name, Payload: payload})
	if err != nil {
		return err
	}
	return p.channel.Publish("", p.queue, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func (p *rabbitPublisher) Health(ctx context.Context) string {
	if p.connection.IsClosed() {
		return "down"
	}
	return "ok"
}

func (p *rabbitPublisher) Close() error {
	_ = p.channel.Close()
	return p.connection.Close()
}
