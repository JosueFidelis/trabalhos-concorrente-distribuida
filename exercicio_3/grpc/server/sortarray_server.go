package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"

	pb "example.com/go-sortarray-grpc/sortarray"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func logErr(err error) {
	if err != nil {
		panic(err)
	}
}

type sortArrayServer struct {
	pb.UnimplementedSortArrayServer
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

func (s *sortArrayServer) Sort(ctx context.Context, in *pb.Arr) (*pb.Arr, error) {
	data := in.GetData()
	sortedData := sortData(data)
	return &pb.Arr{Data: sortedData}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterSortArrayServer(server, &sortArrayServer{})
	log.Printf("server listening ar %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
