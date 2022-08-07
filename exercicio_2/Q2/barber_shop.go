package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

func logIteration(barber_shop_entry *chan string) {
	f, err := os.OpenFile("historico.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	var logEntry string
	for len(*barber_shop_entry) != 0 {
		logEntry = <-*barber_shop_entry
		log.Printf("%s\n", logEntry)
	}

}

func cutHair(isSleeping *bool, barber_shop_entry *chan string, clientNumber int) {
	if *isSleeping {
		*barber_shop_entry <- fmt.Sprintf("cliente %d acordou o barbeiro", clientNumber)
		*isSleeping = false
	}
	time.Sleep(time.Duration(rand.Intn(4)) * time.Nanosecond)
	*barber_shop_entry <- fmt.Sprintf("cliente %d foi atendido", clientNumber)
}

func barber(isClosed *bool, barber_shop_queue *chan int, barber_shop_entry *chan string) {
	var clientNumber int
	isSleeping := false
	for !*isClosed {
		if len(*barber_shop_queue) == 0 {
			if !isSleeping {
				*barber_shop_entry <- "Barbeiro dorme"
				isSleeping = true
			}
		} else {
			clientNumber = <-*barber_shop_queue
			if clientNumber != -1 {
				cutHair(&isSleeping, barber_shop_entry, clientNumber)
			}
		}
	}
}

func client(number int, barber_shop_queue *chan int, barber_shop_entry *chan string, sync_barber_shop *sync.Mutex) {
	if len(*barber_shop_queue) == cap(*barber_shop_queue) {
		*barber_shop_entry <- fmt.Sprintf("cliente %d foi embora, pois a barbearia estava cheia", number)
		return
	}

	sync_barber_shop.Lock()
	defer sync_barber_shop.Unlock()
	*barber_shop_entry <- fmt.Sprintf("cliente %d sentou na fila", number)
	*barber_shop_queue <- number
}

func main() {
	iterations := 100
	queue_size := 6
	var barber_shop_queue = make(chan int, queue_size)
	var barber_shop_entry = make(chan string, 3*iterations)
	var sync_barber_shop sync.Mutex

	isClosed := false
	go barber(&isClosed, &barber_shop_queue, &barber_shop_entry)

	for curr_iteration := 0; curr_iteration < iterations; curr_iteration++ {
		time.Sleep(time.Duration(rand.Intn(3)) * time.Nanosecond)
		go client(curr_iteration, &barber_shop_queue, &barber_shop_entry, &sync_barber_shop)
	}

	//Guarda as cadeiras da fila
	for i := 0; i < queue_size; i++ {
		barber_shop_queue <- -1
	}

	isClosed = true

	logIteration(&barber_shop_entry)
}
