package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func logDuration(fileName string, start time.Time) {
	elapsed := time.Since(start).Microseconds()

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Printf("%d\n", elapsed)
}

func readLoremFile(beginning int, end int) {
	baseFile := "lorem"
	var file strings.Builder

	for i := beginning; i <= end; i++ {
		file.WriteString(fmt.Sprint(baseFile, i))
		file.WriteString(".txt")

		dat, err := os.ReadFile(fmt.Sprintf("./files/%s", file.String()))
		if err != nil {
			panic(err)
		}
		_ = string(dat)

		file.Reset()
	}

}

//O professor falou que o início do programa pode enfrentar instabilidade
func iterationsToDiscard() {
	for i := 0; i < 100; i++ {
		readLoremFile(1, 100)
	}
}

func singleThreadTest() {
	start := time.Now()
	defer logDuration("single_thread_time_elapsed.txt", start)

	readLoremFile(1, 100)
}

func main() {
	iterationsToDiscard()

	for i := 0; i < 30; i++ {
		singleThreadTest()
	}
}

/*func fileWriter() {
	dat, err := os.ReadFile(fmt.Sprintf("./files/%s", "lorem1.txt"))
	if err != nil {
		panic(err)
	}

	fileContent := string(dat)

	d1 := []byte(fileContent)

	baseFile := "lorem"
	var file strings.Builder

	for i := 2; i <= 100; i++ {
		file.WriteString(fmt.Sprint(baseFile, i))
		file.WriteString(".txt")

		os.WriteFile(fmt.Sprintf("./files/%s", file.String()), d1, 0644)

		file.Reset()
	}
}*/
