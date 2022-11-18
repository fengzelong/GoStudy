package main

func main() {
	rabbitmq := NewRabbitMQSimple("GoTestMessage")
	rabbitmq.ConsumeSimple()
}
