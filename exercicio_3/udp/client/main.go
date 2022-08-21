package main

import (
	"fmt"
	"net"
)

func logErr(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	msg := "35,87,52,35,79,62,42,29,23,9,87,29,72,51,80,21,69,8,70,90"

	addr, err := net.ResolveUDPAddr("udp", "localhost:8081")
	logErr(err)

	conn, err := net.DialUDP("udp", nil, addr)
	logErr(err)

	defer conn.Close()

	sendMsg(*conn, msg)

	rcvRep(*conn)
}

func sendMsg(conn net.UDPConn, msg string) {
	req := make([]byte, 1024)

	req = []byte(msg)

	_, err := conn.Write(req)
	logErr(err)
	fmt.Println("Sent request:", string(req[:]))
}

func rcvRep(conn net.UDPConn) {
	rep := make([]byte, 1024)

	_, _, err := conn.ReadFromUDP(rep)
	logErr(err)

	fmt.Println("Received reply:", string(rep[:]))
}
