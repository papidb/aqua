package config

import (
	"fmt"
	"log"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ListenForNotifications(env Env, msgHandler func(msg string)) {
	fmt.Println("Connecting to RabbitMQ...")
	url := urlFromEnv(env)
	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ:", url)
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		env.RabbitMQQueue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for msg := range msgs {
		msgHandler(string(msg.Body))
	}
}

func urlFromEnv(env Env) string {
	_, err := strconv.Atoi(env.RabbitMQPort)
	if err != nil {
		panic(fmt.Sprintf("Invalid AMQP port: %s", env.RabbitMQPort))
	}
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", env.RabbitMQUser, env.RabbitMQPassword, env.RabbitMQHost, env.RabbitMQPort)
}
