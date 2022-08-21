package main

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
)

func logErr(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	msg := "35,87,52,35,79,62,42,29,23,9,87,29,72,51,80,21,69,8,70,90"
	port := 8080

	greetConn := connect(port)

	defer greetConn.Close()

	sendGreeting(greetConn)
	conn := connect(rcvGreetPort(greetConn))

	sendMsg(conn, msg)

	rcvRep(conn)
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
	fmt.Println("Sent request:", string(req))
}

func rcvGreetPort(conn net.UDPConn) int {
	rep := make([]byte, 1024)

	_, _, err := conn.ReadFromUDP(rep)
	logErr(err)
	rep = bytes.Trim(rep, "\x00")

	fmt.Println("Received reply:", string(rep))

	connPort, err := strconv.Atoi(string(rep))
	logErr(err)

	return connPort
}

func sendMsg(conn net.UDPConn, msg string) {
	req := make([]byte, 1024)

	req = []byte(msg)

	_, err := conn.Write(req)
	logErr(err)
	fmt.Println("Sent request:", string(req))
}

func rcvRep(conn net.UDPConn) {
	rep := make([]byte, 1024)

	_, _, err := conn.ReadFromUDP(rep)
	logErr(err)

	fmt.Println("Received reply:", string(rep))
}
