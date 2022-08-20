package main

import (
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
	rep := make([]byte, 1024)
	n := 30

	addr, err := net.ResolveUDPAddr("udp", "localhost:8081")
	logErr(err)

	conn, err := net.DialUDP("udp", nil, addr)
	logErr(err)

	defer conn.Close()

	for i := 0; i < n; i++ {
		go sendMsg(*conn, i)

		_, _, err = conn.ReadFromUDP(rep)
		logErr(err)
		fmt.Println(string(rep[:]))
	}
}

func sendMsg(conn net.UDPConn, i int) {
	req := make([]byte, 1024)

	req = []byte("Mensagem " + strconv.Itoa(i+1))

	_, err := conn.Write(req)
	logErr(err)
}
