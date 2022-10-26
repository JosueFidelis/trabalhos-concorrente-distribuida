package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	proto "example.com/go-sushibar-grpc/sushibar"

	"log"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var client proto.BroadcastClient
var wait *sync.WaitGroup

var availableSeats string
var id string

func init() {
	wait = &sync.WaitGroup{}
}

func connect(user *proto.User) error {
	var streamerror error

	stream, err := client.CreateStream(context.Background(), &proto.Connect{
		User: user,
	})

	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}

	wait.Add(1)
	go func(str proto.Broadcast_CreateStreamClient) {
		defer wait.Done()

		entered := false
		for !entered {
			msg, err := str.Recv()
			if err != nil {
				streamerror = fmt.Errorf("Error reading message: %v", err)
				break
			}

			availableSeats = msg.Content

			if strings.Contains(availableSeats, id) {
				entered = true
			}

			//fmt.Printf("%v : %s\n", msg.Id, msg.Content)

		}
	}(stream)

	return streamerror
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	done := make(chan int)

	id = getId(os.Args)

	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldnt connect to service: %v", err)
	}

	client = proto.NewBroadcastClient(conn)
	user := &proto.User{
		Id: id,
	}

	println(id)
	connect(user)

	wait.Add(1)
	go func() {
		defer wait.Done()

		content := id
		for content != "-1" {
			if strings.Contains(availableSeats, id) {
				content = "-1"
				time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
			} else if strings.Contains(availableSeats, "-1") {
				continue
			}

			msg := &proto.Message{
				Id:      user.Id,
				Content: content,
			}

			_, err := client.BroadcastMessage(context.Background(), msg)

			if err != nil {
				fmt.Printf("Error Sending Message: %v", err)
				break
			}
		}

	}()

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
}

func getId(args []string) string {
	if (len(args) < 2) || os.Args[1] == "" {
		panic("lack of pc argument")
	}

	id := args[1]
	return id
}
