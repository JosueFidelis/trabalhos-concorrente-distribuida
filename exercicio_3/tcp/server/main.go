package main

import (
	"bufio"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
)

func logErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var currPort int = 8080
	dstream, err := net.Listen("tcp", ":8080")
	// defer dstream.Close()

	logErr(err)

	for {
		con, err := dstream.Accept()
		logErr(err)

		currPort++

		_, err = con.Write([]byte(strconv.Itoa(currPort) + "\n"))
		logErr(err)

		go handleConnections(currPort)
	}
}

func sortData(data string) string {

	slc := strings.Split(data, ",")

	var parsedSlc = []int{}

	//parse to int
	for _, i := range slc {
		j, err := strconv.Atoi(i)
		logErr(err)
		parsedSlc = append(parsedSlc, j)
	}
	sort.Ints(parsedSlc)

	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(parsedSlc)), " "), "[]")
}

func handleConnections(currPort int) {
	Dstream, err := net.Listen("tcp", ":"+strconv.Itoa(currPort))
	logErr(err)

	con, err := Dstream.Accept()
	logErr(err)

	for i := 0; i < 10000; i++ {
		data, err := bufio.NewReader(con).ReadString('\n')

		logErr(err)
		//fmt.Println("Got: ", data)

		data = sortData(strings.TrimSpace(data))

		_, err = con.Write([]byte(data + "\n"))
		logErr(err)
	}
	con.Close()
}
