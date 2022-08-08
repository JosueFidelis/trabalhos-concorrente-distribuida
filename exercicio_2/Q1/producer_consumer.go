package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Queue struct {
	queueSize  int
	isAvaiable bool
	sync.Mutex
}

func (q *Queue) Init() {
	q.changeQueueSize(0)
	q.setAvailability(true)
}

func (q *Queue) setAvailability(status bool) {
	q.Mutex.Lock()

	q.isAvaiable = status

	q.Mutex.Unlock()
}

func (q *Queue) changeQueueSize(newSize int) {
	q.Mutex.Lock()

	q.queueSize = newSize

	q.Mutex.Unlock()
}

var consumeCond = sync.NewCond(&sync.Mutex{})
var produceCond = sync.NewCond(&sync.Mutex{})
var getProductQueue *Queue
var entries chan string
var productVessel int

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

func produce(n int) {
	productVessel = n
	entries <- fmt.Sprintf("produtor produziu %d", productVessel)

	for getProductQueue.queueSize > 0 {
		consumeCond.Broadcast()
	}
}

func producer(marketIsOpen *bool) {
	for *marketIsOpen {
		produceCond.L.Lock()

		produceCond.Wait()

		getProductQueue.setAvailability(false)
		entries <- fmt.Sprintf("Produtor acordou")
		produceCond.L.Unlock()

		produce(rand.Intn(100))

		getProductQueue.setAvailability(true)
		entries <- fmt.Sprintf("Produtor limpou o produto")
	}
}

func consume(consumerId int) int {
	consumeCond.L.Lock()

	entries <- fmt.Sprintf("consumidor %d requisitou ao produtor", consumerId)
	getProductQueue.changeQueueSize(getProductQueue.queueSize + 1)

	produceCond.Signal()
	consumeCond.Wait()
	entries <- fmt.Sprintf("consumidor %d acordou", consumerId)

	product := productVessel

	getProductQueue.changeQueueSize(getProductQueue.queueSize - 1)
	entries <- fmt.Sprintf("consumidor %d consumiu %d", consumerId, product)
	consumeCond.L.Unlock()

	return product
}

func consumer(consumerId int, wg *sync.WaitGroup) {
	defer wg.Done()
	for !getProductQueue.isAvaiable {
	}
	consume(consumerId)
}

func main() {
	marketIsOpen := true
	numberOfConsumers := 100
	rand.Seed(time.Now().UnixNano())
	productVessel = 0

	entries = make(chan string, 6*numberOfConsumers)

	getProductQueue = new(Queue)
	getProductQueue.Init()

	var wg sync.WaitGroup
	wg.Add(numberOfConsumers)

	go producer(&marketIsOpen)

	for i := 0; i < numberOfConsumers; i++ {
		go consumer(i, &wg)
	}

	wg.Wait()
	marketIsOpen = false

	logIteration(&entries)
}
