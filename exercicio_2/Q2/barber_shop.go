package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

func logIteration(barberShopEntry *chan string) {
	f, err := os.OpenFile("historico.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	var logEntry string
	for len(*barberShopEntry) != 0 {
		logEntry = <-*barberShopEntry
		log.Printf("%s\n", logEntry)
	}

}

func cutHair(isSleeping *bool, barberShopEntry *chan string, clientNumber int) {
	if *isSleeping {
		*barberShopEntry <- fmt.Sprintf("cliente %d acordou o barbeiro", clientNumber)
		*isSleeping = false
	}
	time.Sleep(time.Duration(rand.Intn(4)) * time.Nanosecond)
	*barberShopEntry <- fmt.Sprintf("cliente %d foi atendido", clientNumber)
}

func barber(isClosed *bool, barberShopQueue *chan int, barberShopEntry *chan string) {
	var clientNumber int
	isSleeping := false
	for !*isClosed {
		if len(*barberShopQueue) == 0 {
			if !isSleeping {
				*barberShopEntry <- "Barbeiro dorme"
				isSleeping = true
			}
		} else {
			clientNumber = <-*barberShopQueue
			if clientNumber != -1 {
				cutHair(&isSleeping, barberShopEntry, clientNumber)
			}
		}
	}
}

func client(number int, barberShopQueue *chan int, barberShopEntry *chan string, syncBarberShop *sync.Mutex) {
	if len(*barberShopQueue) == cap(*barberShopQueue) {
		*barberShopEntry <- fmt.Sprintf("cliente %d foi embora, pois a barbearia estava cheia", number)
		return
	}

	syncBarberShop.Lock()
	defer syncBarberShop.Unlock()
	*barberShopEntry <- fmt.Sprintf("cliente %d sentou na fila", number)
	*barberShopQueue <- number
}

func main() {
	iterations := 100
	queue_size := 6
	var barberShopQueue = make(chan int, queue_size)
	var barberShopEntry = make(chan string, 3*iterations)
	var syncBarberShop sync.Mutex

	isClosed := false
	go barber(&isClosed, &barberShopQueue, &barberShopEntry)

	for currIteration := 0; currIteration < iterations; currIteration++ {
		time.Sleep(time.Duration(rand.Intn(3)) * time.Nanosecond)
		go client(currIteration, &barberShopQueue, &barberShopEntry, &syncBarberShop)
	}

	//Guarda as cadeiras da fila
	for i := 0; i < queue_size; i++ {
		barberShopQueue <- -1
	}

	isClosed = true

	logIteration(&barberShopEntry)
}
