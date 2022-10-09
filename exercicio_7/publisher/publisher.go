package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"task",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	for i := 0; i < 11000; i++ {
		messageToSend := fillMessage()

		err = ch.PublishWithContext(ctx,
			"task", // exchange
			"",     // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        messageToSend,
			})
		failOnError(err, "Failed to publish a message")
	}
}

func fillMessage() []byte {
	currTime := time.Now()
	timeFormat := time.RFC3339Nano

	var messageStr string = currTime.Format(timeFormat)

	messageInBytes := []byte(messageStr)
	padding := make([]byte, (1024 - len(messageInBytes)))
	messageToSend := append(messageInBytes[:], padding[:]...)

	return messageToSend
}
