package main

import "fmt"

var workers []*producers

type producers struct {
	id    int
	quit  chan bool
	queue chan string
}

func multi2multiWithoutWg() {
	const (
		producerCount int = 3
		consumerCount int = 3
	)
	jobQ := make(chan string)
	allDone := make(chan bool)
	workerPool := make(chan *producers)

	for i := 0; i < producerCount; i++ {
		workers = append(workers, &producers{
			queue: make(chan string),
			quit:  make(chan bool),
			id:    i,
		})
		go produce(jobQ, workers[i], workerPool)
	}

	go execute(jobQ, workerPool, allDone)

	for i := 0; i < consumerCount; i++ {
		go consume(i, workerPool)
	}
	<-allDone
}

func execute(jobQ chan<- string, workerPool chan *producers, allDone chan<- bool) {
	for _, j := range messages {
		jobQ <- j
	}
	close(jobQ)
	for _, w := range workers {
		w.quit <- true
	}
	close(workerPool)
	allDone <- true
}

func produce(jobQ <-chan string, p *producers, workerPool chan *producers) {
	for {
		select {
		case msg := <-jobQ:
			{
				workerPool <- p
				if len(msg) > 0 {
					fmt.Printf("Job \"%v\" produced by worker %v\n", msg, p.id)
				}
				p.queue <- msg
			}
		case <-p.quit:
			return
		}
	}
}

func consume(cIdx int, workerPool <-chan *producers) {
	for {
		worker := <-workerPool
		if msg, ok := <-worker.queue; ok {
			if len(msg) > 0 {
				fmt.Printf("Message \"%v\" is consumed by consumer %v from worker %v\n", msg, cIdx, worker.id)
			}
		}
	}
}
