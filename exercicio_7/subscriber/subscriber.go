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

func logPmt(PC int, tmpMean float64) {
	fileName := fmt.Sprintf("log_with_%d_bytes_size.txt", PC)
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Printf("Mean: %f nanoseconds\n", tmpMean)
	log.Printf("Size: %d bytes\n", PC)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	pc := getPC(os.Args)
	timeFormat := time.RFC3339Nano

	numberOfsamples := 10000
	samplesToDiscard := 1000
	pmt := 0.0

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
			samplesReceived, pmt = incrementAndGetPmt(samplesReceived, samplesToDiscard, timeFormat, messageStr, pmt)
		}

		if samplesReceived >= numberOfsamples+samplesToDiscard {
			break
		}
	}

	pmt /= float64(numberOfsamples)

	logPmt(pc, pmt)
}

func incrementAndGetPmt(samplesReceived int, samplesToDiscard int, timeFormat string, messageStr string, pmt float64) (int, float64) {
	samplesReceived++

	if samplesReceived > samplesToDiscard {
		messageTime, err := time.Parse(timeFormat, messageStr)
		failOnError(err, "time in wrong format")

		elapsed := time.Since(messageTime).Nanoseconds()
		pmt += float64(elapsed)
	}
	return samplesReceived, pmt
}

func getMessageFromBody(message []byte) string {
	messageInbytes := bytes.Trim(message, "\x00")
	messageStr := string(messageInbytes)

	return messageStr
}

func getPC(args []string) int {
	var pc int
	if (len(args) < 2) || os.Args[1] == "" {
		panic("lack of pc argument")
	}

	pc, _ = strconv.Atoi(args[1])
	return pc
}
