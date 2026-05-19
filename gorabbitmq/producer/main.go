package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gorabbitmq"
)

func main() {
	queueName := os.Getenv("RABBITMQ_QUEUE")
	if queueName == "" {
		queueName = "GoTestMessage"
	}
	rabbitmq := gorabbitmq.NewRabbitMQSimple(queueName)

	for i := 0; i <= 100; i++ {
		rabbitmq.PublishSimple("Hello go!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
