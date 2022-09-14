package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func fibonacciRPC(data string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

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

	corrId := randomString(32)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for {
		err = ch.PublishWithContext(ctx,
			"",          // exchange
			"rpc_queue", // routing key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: corrId,
				ReplyTo:       q.Name,
				Body:          []byte(data),
			})
		failOnError(err, "Failed to publish a message")

		for d := range msgs {
			if corrId == d.CorrelationId {
				res := string(d.Body)
				failOnError(err, "Failed to convert body to integer")
				fmt.Println(res)
				break
			}
		}
		time.Sleep(time.Second)
	}

	//return
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	data := "53,15,56,41,33,78,42,51,11,8,78,95,33,91,4,36,50,46,56,63,31,84,7,4,44,58,67,66,10,39,75,78,67,95,56,43,57,63,91,45,40,16,38,48,77,17,8,42,75,1,20,29,46,69,62,82,34,1,50,80,31,61,6,39,20,63,84,76,37,26,2,13,13,43,18,8,46,86,81,49,60,12,44,18,3,17,39,48,64,47,53,95,22,94,19,25,3,57,43,59"

	log.Printf(" [x] Requesting fib(%d)", data)
	fibonacciRPC(data)
	//failOnError(err, "Failed to handle RPC request")

	//log.Printf(" [.] Got %d", res)
}
