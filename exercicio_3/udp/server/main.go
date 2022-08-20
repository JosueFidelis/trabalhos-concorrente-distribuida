package main

import (
	"bytes"
	"fmt"
	"net"
	"sort"
	"strconv"
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
	req = bytes.Trim(req, "\x00")

	fmt.Println("Received request:", string(req))

	processReply(req, &rep)

	_, err = conn.WriteTo(rep, addr)
	logErr(err)

	fmt.Println("Sent reply:", string(rep))
}

func processReply(req []byte, rep *[]byte) {
	numStrList := strings.Split(string(req), ",")
	var numList = []int{}

	for _, i := range numStrList {
		j, err := strconv.Atoi(i)
		logErr(err)
		numList = append(numList, j)
	}

	sort.Ints(numList)

	//Converte em String de novo
	*rep = []byte(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(numList)), ", "), "[]"))
}
