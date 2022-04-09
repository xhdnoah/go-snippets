package main

import (
	"fmt"
	"time"
)

func ticker() {
	ch := make(chan bool)
	d := time.NewTicker(1 * time.Second)

	go func() {
		time.Sleep(7 * time.Second)
		ch <- true
	}()

	for {
		select {
		case <-ch:
			fmt.Println("Completed!")
			return
		case tm := <-d.C:
			fmt.Println("The current time is: ", tm)
		}
	}
}

func stop() {
	ch := make(chan bool)
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ch:
				return
			case tm := <-ticker.C:
				fmt.Println("The current time is: ", tm)
			}
		}
	}()

	time.Sleep(7 * time.Second)
	ticker.Stop()
	ch <- true
	fmt.Println("Ticker is turned off!")
}
