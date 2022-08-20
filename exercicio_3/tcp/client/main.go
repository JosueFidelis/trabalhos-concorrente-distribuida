package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func logErr(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	//connect to this socket
	connClient, _ := net.Dial("tcp", "127.0.0.1:8080")

	for {

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, err := reader.ReadString('\n')
		logErr(err)

		//send to socket
		fmt.Fprint(connClient, text+"\n")

		//listen for reply
		//message, _ := bufio.NewReader(connClient).ReadString('\n')
		//fmt.Print("Message from server: " + message)
	}
}
