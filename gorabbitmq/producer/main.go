package main

import (
	"fmt"
	"strconv"
	"time"

	"gorabbitmq"
)

func main() {
	rabbitmq := gorabbitmq.NewRabbitMQSimple("GoTestMessage")

	for i := 0; i <= 100; i++ {
		rabbitmq.PublishSimple("Hello go!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
