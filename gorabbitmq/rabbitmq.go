package gorabbitmq

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

const defaultMQURL = "amqp://guest:guest@127.0.0.1:5672/"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机名称
	Exchange string
	// 绑定键名称
	Key string
	// 连接地址
	MqUrl string
}

// NewRabbitMQ 创建 RabbitMQ 实例。
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	mqURL := os.Getenv("RABBITMQ_URL")
	if mqURL == "" {
		mqURL = defaultMQURL
		fmt.Println("未设置 RABBITMQ_URL，使用本地默认连接示例")
	}
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, MqUrl: mqURL}
}

// Destory 关闭 channel 和 connection。
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

// failOnErr 处理 RabbitMQ 错误。
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// NewRabbitMQSimple 创建简单模式下的 RabbitMQ 实例。
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MqUrl)
	rabbitmq.failOnErr(err, "连接 RabbitMQ 失败")

	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "打开 channel 失败")
	return rabbitmq
}

// PublishSimple 使用简单模式发送消息。
func (r *RabbitMQ) PublishSimple(message string) {
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// ConsumeSimple 使用简单模式消费消息。
func (r *RabbitMQ) ConsumeSimple() {
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("收到消息: %s", d.Body)
		}
	}()

	log.Printf("等待消息中，按 CTRL+C 退出")
	<-forever
}
