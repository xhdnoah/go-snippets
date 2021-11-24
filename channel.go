package snippets

import (
	"fmt"
	"sync"
	"time"
)

// using range while reading multiple values
func foo() chan int {
	ch := make(chan int)

	go func() {
		ch <- 1
		ch <- 2
		ch <- 3
		close(ch)
	}()

	return ch
}

func main() {
	for n := range foo() {
		fmt.Println(n)
	}
	fmt.Println("channel is now closed")
}

// channels are often used to implement timeout
func timeout() {
	// empty struct is useful in channels when you notify about some event
	// but you don't need to pass any information about it, only a fact.
	// Best solution is to pass an empty struct beacuse it will only
	// increment a counter in the channel but not assign memory.
	ch := make(chan struct{}, 1)

	go func() {
		time.Sleep(10 * time.Second)
		ch <- struct{}{}
	}()

	select {
	case <-ch:
		// Work completed before timeout
	case <-time.After(1 * time.Second):
		// Work was not completed after 1 sec
	}
}

// coordinating goroutines
// imagine a goroutine with a two-step process
// where the main thread needs to do some work between each step
func coordinate() {
	ch := make(chan struct{})
	go func() {
		// Wait for main thread's singal to begin step one 阻塞
		<-ch

		// Perform work
		time.Sleep(1 * time.Second)

		// Signal to main thread that step one has completed
		ch <- struct{}{}

		// Wait for main thread's singal to begin step two
		<-ch

		// Perform work
		time.Sleep(1 * time.Second)

		// Signal to main thread that work has completed
		ch <- struct{}{}
	}()

	// Notify goroutines that step one can begin
	ch <- struct{}{}

	// Wait for notification from goroutines that step one has completed
	<-ch

	// Perform some work before we notify
	// the goroutine that step two can begin
	time.Sleep(1 * time.Second)

	// Notify goroutine that step two can begin
	ch <- struct{}{}

	// Wait for notification from goroutines that step two has completed
	<-ch
}

// Buffered vs unbuffered
func bufferedUnbuffered(buffered bool) {
	var ch chan int
	if buffered {
		ch = make(chan int, 3)
	} else {
		ch = make(chan int)
	}

	go func() {
		for i := 0; i < 7; i++ {
			// If the channel is buffered, then while there's an empty
			// "slot" in the channel, sending to it will not be a
			// blocking operation. If the channel is full, however, we'll
			// have to wait until a "slot" frees up.
			// If the channel is unbuffered, sending will block until
			// there's a receiver ready to take the value. This is great
			// for goroutine synchronization, not so much for queueing
			// tasks for instance in a webserver, as the request will
			// hang until the worker is ready to take our task
			fmt.Println(">", "Sending", i, "...")
			ch <- i
			fmt.Println(">", i, "Sent")
			time.Sleep(25 * time.Millisecond)
		}
		// we'll close the channel, so that the range over channel below can terminate.
		close(ch)
	}()

	for i := range ch {
		// For each task sent on the channel, we would perform some
		// task, In this case, we'll assume the job is to
		// "sleep 100ms"
		fmt.Println("<", i, "received, performing 100ms job")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("<", i, "job done")
	}
}

// Blocking & unblocking
// By default communication over the channels is sync
// When you send some value there must be a receiver
// Otherwise the code will get fatal error as follows:
func deadlock() {
	msg := make(chan string)
	msg <- "Hey"
	// fatal error: all goroutines are asleep - deadlock
	go func() {
		fmt.Println(<-msg)
	}()
}

// solution: use buffered channels
func buffered() {
	msg := make(chan string, 1)
	msg <- "Hey"
	go func() {
		fmt.Println(<-msg)
	}()
	time.Sleep(time.Second * 1)
}

// waiting for work to finish
// A common technique for using channels is to create
// some number of workers (or consumers) to read from
// the channel. Using a sync.WaitGroup is an easy way to
// wait for those workers to finish running
func workers() {
	numberPiecesOfWork := 20
	numWorkers := 5

	workCh := make(chan int)
	wg := &sync.WaitGroup{}

	// start workers
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker(workCh, wg)
	}

	for i := 0; i < numberPiecesOfWork; i++ {
		work := i % 10 // invent some work
		workCh <- work
	}

	// tell workers that no more work is coming
	close(workCh)

	// wait for workers to finish
	wg.Wait()

	fmt.Println("done")
}

func worker(workCh <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for work := range workCh {
		doWork(work)
	}
}

func doWork(work int) {
	time.Sleep(time.Duration(work) * time.Second)
	fmt.Println("slept for", work, "milliseconds")
}

func fibonacci() {
	c := make(chan int)
	quit := make(chan struct{})
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- struct{}{}
	}()

	func() {
		x, y := 0, 1
		// select, same as the switch, is not a loop
		// will only select one case to process
		// unless an infinite for loop could be added
		for {
			select {
			case c <- x:
				x, y = y, x+y
			case <-quit:
				fmt.Println("quit")
				return
			}
		}
	}()
}

// 0
// 1
// 1
// 2
// 3
// 5
// 8
// 13
// 21
// 34
// quit
