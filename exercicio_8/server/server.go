package main

import (
	"log"
	"net"
	"os"
	"sync"

	proto "example.com/go-sushibar-grpc/sushibar"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	glog "google.golang.org/grpc/grpclog"
)

var grpcLog glog.LoggerV2
var seatCount int = 0

func init() {
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

type Connection struct {
	stream proto.Broadcast_CreateStreamServer
	id     string
	error  chan error
}

type Server struct {
	Connection []*Connection
	Seats      []string
	proto.UnimplementedBroadcastServer
}

func (s *Server) CreateStream(pconn *proto.Connect, stream proto.Broadcast_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     pconn.User.Id,
		error:  make(chan error),
	}

	s.Connection = append(s.Connection, conn)

	return <-conn.error
}

func (s *Server) removeElement(id string) {
	index := -1
	for i := 0; i < len(s.Connection); i++ {
		if id == s.Connection[i].id {
			index = i
		}
	}

	if index == len(s.Connection)-1 {
		s.Connection = s.Connection[:index]
	} else {
		s.Connection = append(s.Connection[:index], s.Connection[index+1:]...)
	}

	for i := 1; i < 6; i++ {
		if s.Seats[i] == id {
			s.Seats[i] = "0"
			seatCount--
			grpcLog.Info("Imma head out: ", id)
		}
	}

	grpcLog.Info("Occupied seats: ", seatCount)

	if seatCount == 0 {
		s.Seats[0] = "0"
	}
}

func (s *Server) addElement(id string) {
	if s.Seats[0] == "-1" || seatCount >= 6 {
		return
	}

	availableSeat := -1
	for i := 1; i < 6; i++ {
		if s.Seats[i] == id {
			return
		} else if s.Seats[i] == "0" {
			availableSeat = i
		}
	}
	s.Seats[availableSeat] = id
	seatCount++

	grpcLog.Info("adding: ", id)
	grpcLog.Info("in seat: ", availableSeat)
	grpcLog.Info("Occupied seats: ", seatCount)

	if seatCount == 5 {
		s.Seats[0] = "-1"
	}
	msg := s.Seats[0] + "," + s.Seats[1] + "," + s.Seats[2] + "," + s.Seats[3] + "," + s.Seats[4] + "," + s.Seats[5]
	grpcLog.Info("seats: ", msg)
}

func (s *Server) BroadcastMessage(ctx context.Context, msg *proto.Message) (*proto.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	if msg.Content == "-1" {
		s.removeElement(msg.Id)

	} else {
		s.addElement(msg.Id)
	}

	msg.Content = s.Seats[0] + "," + s.Seats[1] + "," + s.Seats[2] + "," + s.Seats[3] + "," + s.Seats[4] + "," + s.Seats[5]
	grpcLog.Info("Sending: ", msg.Content)

	for _, conn := range s.Connection {
		wait.Add(1)

		go func(msg *proto.Message, conn *Connection) {
			defer wait.Done()

			err := conn.stream.Send(msg)
			//grpcLog.Info("Sending message to: ", conn.stream)

			if err != nil {
				grpcLog.Errorf("Error with Stream: %v - Error: %v", conn.stream, err)
				conn.error <- err
			}

		}(msg, conn)
	}

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
	return &proto.Close{}, nil
}

func main() {
	var connections []*Connection
	var seats []string
	firstPosition := "0"
	seats = append(seats, firstPosition, firstPosition, firstPosition, firstPosition, firstPosition, firstPosition)

	server := &Server{
		Connection: connections,
		Seats:      seats,
	}

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error creating the server %v", err)
	}

	grpcLog.Info("Starting server at port :8080")

	proto.RegisterBroadcastServer(grpcServer, server)
	grpcServer.Serve(listener)
}
