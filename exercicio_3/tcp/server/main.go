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
		return
	}
}

func main() {
	dstream, err := net.Listen("tcp", ":8080")
	// defer dstream.Close()

	logErr(err)

	for {
		con, err := dstream.Accept()
		logErr(err)

		go handle(con)
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

func handle(con net.Conn) {
	for i := 0; i < 2; i++ {
		data, err := bufio.NewReader(con).ReadString('\n')

		logErr(err)
		fmt.Println("Got:       ", data)

		data = sortData(strings.TrimSpace(data))

		_, err = con.Write([]byte(data + "\n"))
		logErr(err)
	}
	con.Close()
}
