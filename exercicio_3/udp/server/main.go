package main

import (
	"fmt"
	"net"
	"strings"
)

func logErr(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {

	addr, err := net.ResolveUDPAddr("udp", "localhost:8081")
	logErr(err)

	conn, err := net.ListenUDP("udp", addr)
	logErr(err)

	defer conn.Close()

	fmt.Println("Servidor UDP aguardando requests...")

	for {
		handle(*conn)
	}

}

func handle(conn net.UDPConn) {

	req := make([]byte, 1024)
	rep := make([]byte, 1024)

	_, addr, err := conn.ReadFromUDP(req)
	logErr(err)

	fmt.Println(string(req[:]))

	rep = []byte(strings.ToUpper(string(req)))

	fmt.Println(string(rep[:]))

	_, err = conn.WriteTo(rep, addr)
	logErr(err)
}
