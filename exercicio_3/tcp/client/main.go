package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

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

func main() {
	//connect to this socket
	numberOfClientsRunning, _ := strconv.Atoi(os.Args[1])
	numberOfIterations := 10000
	numberOfIterationsToDiscard := 1000

	var rttArray = make([]float64, numberOfIterations)
	rttMean := 0.0

	connClient, _ := net.Dial("tcp", "127.0.0.1:8080")

	//listen for reply
	newPort, _ := bufio.NewReader(connClient).ReadString('\n')

	fmt.Print("Port from server: " + newPort)
	connClient.Close()

	newConnClient, _ := net.Dial("tcp", "127.0.0.1:"+strings.TrimSpace(newPort))
	defer newConnClient.Close()

	text := "35,87,52,35,79,62,42,29,23,9,87,29,72,51,80,21,69,8,70,90\n"

	for i := 0; i < numberOfIterations+numberOfIterationsToDiscard; i++ {

		//send to socket
		// fmt.Println(text)
		start := time.Now()
		fmt.Fprint(newConnClient, text)

		//listen for reply
		bufio.NewReader(newConnClient).ReadString('\n')
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
			fmt.Print(rttArray[i])
			fmt.Print("  ;  ")
		}
		fmt.Print("\n")

		rttSd = math.Sqrt(rttSd / float64(numberOfIterations))

		logRtt(numberOfClientsRunning, rttMean, rttSd)
	}

}
