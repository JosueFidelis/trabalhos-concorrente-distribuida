package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func logTmp(messageSize int, tmpMean float64) {
	fileName := fmt.Sprintf("log_with_%d_bytes_size.txt", messageSize)
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Printf("Mean: %f nanoseconds\n", tmpMean)
	log.Printf("Size: %d bytes\n", messageSize)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	messageSize := getMessageSize(os.Args)
	timeFormat := time.RFC3339Nano
	fmt.Println(messageSize)

	numberOfsamples := 10000
	samplesToDiscard := 1000
	tmp := 0.0

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

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"task", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	samplesReceived := 0

	for msg := range msgs {
		messageStr := getMessageFromBody(msg.Body)

		if messageStr != "no" {
			samplesReceived++

			if samplesReceived > samplesToDiscard {
				messageTime, err := time.Parse(timeFormat, messageStr)
				failOnError(err, "time in wrong format")

				elapsed := time.Since(messageTime).Nanoseconds()
				tmp += float64(elapsed)
			}
		}

		if samplesReceived >= numberOfsamples+samplesToDiscard {
			break
		}
	}

	tmp /= float64(numberOfsamples)

	logTmp(messageSize, tmp)
}

func getMessageFromBody(message []byte) string {
	messageInbytes := bytes.Trim(message, "\x00")
	messageStr := string(messageInbytes)

	return messageStr
}

func getMessageSize(args []string) int {
	var size int
	if (len(args) < 2) || os.Args[1] == "" {
		panic("lack of size argument")
	}

	size, _ = strconv.Atoi(args[1])
	return size
}
