package main

import (
	"fmt"
	"sync"
)

func main() {
	const producerCount int = 5
	const consumerCount int = 3
	link := make(chan string)
	wp := &sync.WaitGroup{}
	wc := &sync.WaitGroup{}

	wp.Add(producerCount)
	wc.Add(consumerCount)

	for i := 0; i < producerCount; i++ {
		go multipleProducer(link, i, wp)
	}
	for i := 0; i < consumerCount; i++ {
		go multipleConsumer(link, i, wc)
	}
	wp.Wait()
	close(link)
	wc.Wait()
}

func multipleProducer(link chan<- string, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, msg := range multiMessages[id] {
		link <- msg
	}
}

func multipleConsumer(link <-chan string, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range link {
		fmt.Printf("Message \"%v\" is consumed by consumer %v\n", msg, id)
	}
}
