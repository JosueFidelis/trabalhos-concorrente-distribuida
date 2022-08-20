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
	msg := "1,9,3,8,1,4,1,9,5,1,6,3,1,3,5,1,3,4,4,0"

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
	fmt.Println("Sent request: ", string(req[:]))
}

func rcvRep(conn net.UDPConn) {
	rep := make([]byte, 1024)

	_, _, err := conn.ReadFromUDP(rep)
	logErr(err)

	fmt.Println("Received reply: ", string(rep[:]))
}
