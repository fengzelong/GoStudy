package main

import (
	"os"

	"gorabbitmq"
)

func main() {
	queueName := os.Getenv("RABBITMQ_QUEUE")
	if queueName == "" {
		queueName = "GoTestMessage"
	}
	rabbitmq := gorabbitmq.NewRabbitMQSimple(queueName)
	rabbitmq.ConsumeSimple()
}
