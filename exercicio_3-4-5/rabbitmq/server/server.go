package main

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func logErr(err error) {
	if err != nil {
		panic(err)
	}
}

func sortData(data string) string {
	slc := strings.Split(data, ",")

	var parsedSlc = []int{}

	//parse to int
	for _, i := range slc {
		j, err := strconv.Atoi(i)
		logErr(err)
		parsedSlc = append(parsedSlc, j)
	}
	sort.Ints(parsedSlc)

	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(parsedSlc)), " "), "[]")
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	logErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	logErr(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	logErr(err)

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	logErr(err)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	logErr(err)

	func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for msg := range msgs {
			array := string(msg.Body)

			response := sortData(array)

			err := ch.PublishWithContext(ctx,
				"",          // exchange
				msg.ReplyTo, // routing key
				false,       // mandatory
				false,       // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: msg.CorrelationId,
					Body:          []byte(response),
				})
			logErr(err)

			msg.Ack(false)
		}
	}()
}
