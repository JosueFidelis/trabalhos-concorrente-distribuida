package main

import (
	"math/rand"
	"sync"
)

func produce(productVessel *int) {
	*productVessel = rand.Intn(100)
}

func producer(productVessel *int, marketIsOpen *bool, consumeCond *sync.Cond) {
	for *marketIsOpen {

	}
}

func consume(productVessel *int, consumeCond *sync.Cond) int {
	consumeCond.L.Lock()

	for *productVessel == -1 {
		consumeCond.Wait()
	}

	product := *productVessel

	return product
}

func consumer(consumerId int, productVessel *int, consumeCond *sync.Cond, wg *sync.WaitGroup) {
	defer wg.Done()

	consume(productVessel, consumeCond)
}

func main() {
	marketIsOpen := true
	productVessel := -1

	numberOfConsumers := 100

	consumeCond := sync.NewCond(&sync.Mutex{})

	var wg sync.WaitGroup
	wg.Add(numberOfConsumers)

	go producer(&productVessel, &marketIsOpen, consumeCond)

	for i := 0; i < numberOfConsumers; i++ {
		go consumer(i, &productVessel, consumeCond, &wg)
	}

	wg.Wait()
	marketIsOpen = false

	//log all entries
}
