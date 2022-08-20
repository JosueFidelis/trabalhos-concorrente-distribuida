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
	//connect to this socket
	connClient, _ := net.Dial("tcp", "127.0.0.1:8080")

	text := "35,87,52,35,79,62,42,29,23,9,87,29,72,51,80,21,69,8,70,90\n"
	//            10000
	for i := 0; i < 2; i++ {

		// reader := bufio.NewReader(os.Stdin)
		// fmt.Print("Text to send: ")
		// text, err := reader.ReadString('\n')
		// logErr(err)

		//send to socket
		fmt.Println(text)
		fmt.Fprint(connClient, text)

		//listen for reply
		message, _ := bufio.NewReader(connClient).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}
