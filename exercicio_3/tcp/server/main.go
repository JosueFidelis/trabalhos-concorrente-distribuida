package main

import (
	"bufio"
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
	dstream, err := net.Listen("tcp", ":8080")
	defer dstream.Close()

	logErr(err)

	for {
		con, err := dstream.Accept()
		logErr(err)

		go handle(con)
	}
}

func handle(con net.Conn) {
	for {
		data, err := bufio.NewReader(con).ReadString('\n')

		logErr(err)
		fmt.Println(data)
	}
}
