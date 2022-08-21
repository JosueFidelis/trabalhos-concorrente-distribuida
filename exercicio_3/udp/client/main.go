package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"time"
)

func logErr(err error) {
	if err != nil {
		fmt.Println(err)
		return
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

func main() {
	numberOfClientsRunning, _ := strconv.Atoi(os.Args[1])
	numberOfIterations := 10000
	numberOfIterationsToDiscard := 1000

	var rttArray = make([]float64, numberOfIterations)
	rttMean := 0.0

	msg := "35,87,52,35,79,62,42,29,23,9,87,29,72,51,80,21,69,8,70,90"
	port := 8080

	greetConn := connect(port)
	defer greetConn.Close()

	greetReceived := false
	var conn net.UDPConn
	var connPort int

	for !greetReceived {
		sendGreeting(greetConn)
		connPort = rcvGreetPort(greetConn, &greetReceived)
	}

	conn = connect(connPort)

	for i := 0; i < numberOfIterations+numberOfIterationsToDiscard; i++ {
		start := time.Now()
		sendMsg(conn, msg)
		rcvRep(conn)
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

func connect(port int) net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", "localhost:"+strconv.Itoa(port))
	logErr(err)

	conn, err := net.DialUDP("udp", nil, addr)
	logErr(err)

	fmt.Println("Cliente UDP conectando na porta:", port, "...")

	return *conn
}

func sendGreeting(conn net.UDPConn) {
	req := make([]byte, 1024)

	req = []byte("hi")

	_, err := conn.Write(req)
	logErr(err)
	//fmt.Println("Sent request:", string(req))
}

func rcvGreetPort(conn net.UDPConn, greetReceived *bool) int {
	rep := make([]byte, 1024)
	timeOutSec := 1
	conn.SetReadDeadline(time.Now().Add(time.Duration(timeOutSec) * time.Second))

	_, _, err := conn.ReadFromUDP(rep)
	logErr(err)

	if !errors.Is(err, os.ErrDeadlineExceeded) {
		*greetReceived = true
	}

	if *greetReceived {
		rep = bytes.Trim(rep, "\x00")

		//fmt.Println("Received reply:", string(rep))

		connPort, err := strconv.Atoi(string(rep))
		logErr(err)

		return connPort
	}

	return -1
}

func sendMsg(conn net.UDPConn, msg string) {
	req := make([]byte, 1024)

	req = []byte(msg)

	_, err := conn.Write(req)
	logErr(err)
	//fmt.Println("Sent request:", string(req))
}

func rcvRep(conn net.UDPConn) {
	rep := make([]byte, 1024)

	_, _, err := conn.ReadFromUDP(rep)
	logErr(err)

	//fmt.Println("Received reply:", string(rep))
}
