package main

import "fmt"

// Define the Producer-Consumer Problem:
// A producer that produces jobs and a consumer that consumes them.
// The producer and consumer share a buffer of a limited size.
// The producer puts its jobs into the buffer and consumer takes a job from the buffer and processes it.
// thinking of the jobs as a scheduled tasks or data, or incoming requests.
func main() {
	link := make(chan string)
	done := make(chan bool)
	go single_producer(link)
	go single_consumer(link, done)
	<-done
}

func single_producer(link chan<- string) {
	for _, msg := range messages {
		link <- msg
	}
	close(link)
}

func single_consumer(link <-chan string, done chan<- bool) {
	for msg := range link {
		fmt.Println(msg)
	}
	done <- true
}
