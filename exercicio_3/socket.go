package main

import {
	"bufio"
	"fmt"
	"net"
}

func main() {
	dstream, err := net.Listen("tcp", ":8080")
	defer dstream.Close()

	if(err != nil){
		fmt.println(err)
		return
	}

	for {
		con, err := dstream.Accept()
		if(err != nil){
			fmt.println(err)
			return
		}
		
		go handle(con)
	}
}

func handle(con net.Conn) {

}