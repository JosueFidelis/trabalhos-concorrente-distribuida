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
	port := 8080
	currPort := port

	greetConn := createConn(port)
	defer greetConn.Close()

	for {
		handleFirstContact(greetConn, &currPort)
	}
}

func createConn(port int) net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", "localhost:"+strconv.Itoa(port))
	logErr(err)

	conn, err := net.ListenUDP("udp", addr)
	logErr(err)

	fmt.Println("Servidor UDP aguardando requests na porta:", port, "...")

	return *conn
}

func handleFirstContact(greetConn net.UDPConn, currPort *int) {
	req := make([]byte, 1024)
	rep := make([]byte, 1024)

	_, addr, err := greetConn.ReadFromUDP(req)
	logErr(err)
	//fmt.Println("Request de apresentação recebido:", string(req))

	*currPort += 1

	rep = []byte(strconv.Itoa(*currPort))
	createClientConnection(*currPort)

	_, err = greetConn.WriteTo(rep, addr)
	logErr(err)
	fmt.Println("Resposta de apresentação enviada:", string(rep))
}

func createClientConnection(port int) {
	conn := createConn(port)
	go handleClientConnection(conn)
}

func handleClientConnection(conn net.UDPConn) {
	defer conn.Close()

	for {
		req := make([]byte, 1024)
		rep := make([]byte, 1024)

		_, addr, err := conn.ReadFromUDP(req)
		logErr(err)
		req = bytes.Trim(req, "\x00")

		//fmt.Println("Received request:", string(req))

		processReply(req, &rep)

		_, err = conn.WriteTo(rep, addr)
		logErr(err)

		//fmt.Println("Sent reply:", string(rep))
	}
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
