package main

import "fmt"

const consumerCount int = 3

func main() {
	jobs := make(chan string)
	done := make(chan bool)

	go single_producer(jobs)
	for i := 0; i < consumerCount; i++ {
		go multiple_consumer(i, jobs, done)
	}
	<-done
}

func single_producer(jobs chan<- string) {
	for _, msg := range messages {
		jobs <- msg
	}
	close(jobs)
}

func multiple_consumer(worker int, jobs <-chan string, done chan<- bool) {
	for msg := range jobs {
		fmt.Printf("message %v is consumed by worker %v.\n", msg, worker)
	}
	done <- true
}
