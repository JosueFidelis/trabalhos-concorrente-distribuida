package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	size, flag := getFlagAndSize(os.Args)

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
		messageToSend := fillMessage(size, flag)

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

	log.Printf(" [x] Sent %s", flag)
}

func fillMessage(size int, flag string) []byte {
	currTime := time.Now()
	timeFormat := time.RFC3339Nano

	var messageStr string
	if flag == "1" {
		messageStr = currTime.Format(timeFormat)
	} else {
		messageStr = "no"
	}

	messageInBytes := []byte(messageStr)
	padding := make([]byte, (size - len(messageInBytes)))
	messageToSend := append(messageInBytes[:], padding[:]...)

	return messageToSend
}

func getFlagAndSize(args []string) (int, string) {
	var size int
	var flag string
	if (len(args) < 2) || os.Args[1] == "" {
		panic("lack of size argument")
	}

	if len(args) < 3 || os.Args[2] == "" {
		flag = "0"
	} else {
		flag = args[2]
	}

	fmt.Println(flag)

	size, _ = strconv.Atoi(args[1])
	return size, flag
}
