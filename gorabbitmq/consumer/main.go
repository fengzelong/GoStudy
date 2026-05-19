package main

import "gorabbitmq"

func main() {
	rabbitmq := gorabbitmq.NewRabbitMQSimple("GoTestMessage")
	rabbitmq.ConsumeSimple()
}
