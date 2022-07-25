package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var n_lorem = 100

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

func readLoremFile(beginning int, end int, active_thread_pointer *chan int) {
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

	*active_thread_pointer <- 1
}

//O professor falou que o início do programa pode enfrentar instabilidade
func iterationsToDiscard() {
	//Cria um canal e preenche ele, só para o readLoremFile conseguir remover algo do canal ao final
	var active_thread_chan = make(chan int, 1)

	for i := 0; i < 100; i++ {
		readLoremFile(1, n_lorem, &active_thread_chan)
		<-active_thread_chan
	}
}

func createLoremThreadTest(n_threads int) {
	var out_file_name strings.Builder

	out_file_name.WriteString(fmt.Sprint(n_threads))
	out_file_name.WriteString("_thread_elapsed_time.txt")

	start := time.Now()
	defer logDuration(out_file_name.String(), start)

	// Monitora threads ativas para que a função createLoremThreadTest espere que todas as threads terminem
	var active_thread_chan = make(chan int, n_threads)

	s_range := 1
	for i := 0; i < n_threads; i++ {
		e_range := (s_range - 1) + n_lorem/n_threads
		go readLoremFile(s_range, e_range, &active_thread_chan)
		s_range += n_lorem / n_threads
	}

	// Tenta esvaziar o canal para garantir que as threads terminaram
	for i := 0; i < n_threads; i++ {
		<-active_thread_chan
	}
}

func main() {
	iterationsToDiscard()

	for i := 0; i < 30; i++ {
		// Single Thread
		createLoremThreadTest(1)
	}

	for i := 0; i < 30; i++ {
		// Two Threads
		createLoremThreadTest(2)
	}

	for i := 0; i < 30; i++ {
		// Five Threads
		createLoremThreadTest(5)
	}

	for i := 0; i < 30; i++ {
		// Ten Threads
		createLoremThreadTest(10)
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
