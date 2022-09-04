package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	pb "example.com/go-sortarray-grpc/sortarray"
	"google.golang.org/grpc"
)

func logErr(err error) {
	if err != nil {
		panic(err)
	}
}

func logRtt(numberOfClientsRunning int, rttMean float64, rttSd float64) {
	fileName := fmt.Sprintf("log_with_%d_clients.txt", numberOfClientsRunning)
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Printf("Mean: %f nanoseconds\n", rttMean)
	log.Printf("Standard deviation: %f nanoseconds\n", rttSd)
}

const (
	address = "localhost:50051"
)

func main() {
	numberOfClientsRunning, _ := strconv.Atoi(os.Args[1])
	numberOfIterations := 10000
	numberOfIterationsToDiscard := 1000

	var rttArray = make([]float64, numberOfIterations)
	rttMean := 0.0

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	logErr(err)

	defer conn.Close()
	client := pb.NewSortArrayClient(conn)

	data := "53,15,56,41,33,78,42,51,11,8,78,95,33,91,4,36,50,46,56,63,31,84,7,4,44,58,67,66,10,39,75,78,67,95,56,43,57,63,91,45,40,16,38,48,77,17,8,42,75,1,20,29,46,69,62,82,34,1,50,80,31,61,6,39,20,63,84,76,37,26,2,13,13,43,18,8,46,86,81,49,60,12,44,18,3,17,39,48,64,47,53,95,22,94,19,25,3,57,43,59"
	for i := 0; i < numberOfIterations+numberOfIterationsToDiscard; i++ {
		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		start := time.Now()
		_, err = client.Sort(ctx, &pb.Arr{Data: data})
		logErr(err)
		elapsed := time.Since(start).Nanoseconds()

		if i-numberOfIterationsToDiscard >= 0 {
			rttArray[i-numberOfIterationsToDiscard] = float64(elapsed)
			rttMean += float64(elapsed)
		}
	}

	if numberOfClientsRunning != -1 {
		rttMean /= float64(numberOfIterations)
		rttSd := 0.0

		for i := 0; i < numberOfIterations; i++ {
			rttSd += math.Pow(rttArray[i]-rttMean, 2)
		}

		rttSd = math.Sqrt(rttSd / float64(numberOfIterations))

		logRtt(numberOfClientsRunning, rttMean, rttSd)
	}
}
