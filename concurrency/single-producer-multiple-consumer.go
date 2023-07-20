package main

import "fmt"

func single2multi() {
	const consumerCount int = 3
	jobs := make(chan string)
	done := make(chan bool)

	go singleProducer(jobs)
	for i := 0; i < consumerCount; i++ {
		go multiConsumer(i, jobs, done)
	}
	<-done
}

func singleProducer(jobs chan<- string) {
	for _, msg := range messages {
		jobs <- msg
	}
	close(jobs)
}

func multiConsumer(worker int, jobs <-chan string, done chan<- bool) {
	for msg := range jobs {
		fmt.Printf("message %v is consumed by worker %v.\n", msg, worker)
	}
	done <- true
}
