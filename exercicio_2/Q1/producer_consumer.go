package main

import (
	"fmt"
	"math/rand"
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
var productVessel int

func produce(n int) {
	productVessel = n
	fmt.Printf("produtor produziu %d\n", productVessel)

	for getProductQueue.queueSize > 0 {
		consumeCond.Broadcast()
	}
}

func producer(marketIsOpen *bool) {
	for *marketIsOpen {
		produceCond.L.Lock()

		produceCond.Wait()

		getProductQueue.setAvailability(false)
		fmt.Println("Produtor acordou")
		produceCond.L.Unlock()

		produce(rand.Intn(100))

		getProductQueue.setAvailability(true)
		fmt.Println("Produtor limpou o produto")
	}
}

func consume(consumerId int) int {
	consumeCond.L.Lock()

	getProductQueue.changeQueueSize(getProductQueue.queueSize + 1)

	fmt.Printf("consumidor %d requisitou ao produtor\n", consumerId)

	produceCond.Signal()
	consumeCond.Wait()
	fmt.Printf("consumidor %d acordou\n", consumerId)

	product := productVessel

	getProductQueue.changeQueueSize(getProductQueue.queueSize - 1)
	fmt.Printf("consumidor %d consumiu %d\n", consumerId, product)
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

	//log all entries
}
