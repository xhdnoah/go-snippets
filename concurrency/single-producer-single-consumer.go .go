package main

import (
	"fmt"
	"time"
)

// Define the Producer-Consumer Problem:
// A producer that produces jobs and a consumer that consumes them.
// The producer and consumer share a buffer of a limited size.
// The producer puts its jobs into the buffer and consumer takes a job from the buffer and processes it.
// thinking of the jobs as a scheduled tasks or data, or incoming requests.
func single2single() {
	// 无缓冲通道强制要求生产者和消费者之间的消息传递同步进行 两者必须配合否则会产生 deadlock
	// 有缓冲通道提供一定的异步性 解耦生产者和消费者 需要注意控制缓冲区大小以避免资源浪费或缓冲区溢出
	link := make(chan string)
	// queue := make(chan string, 3)
	done := make(chan bool)
	go producer(link)
	go consumer(link, done)
	<-done
}

func producer(link chan<- string) {
	for _, msg := range messages {
		link <- msg             // 数据发送到 channel
		time.Sleep(time.Second) // 模拟生产耗时
	}
	close(link) // 关闭 channel 表示生产结束
}

func consumer(link <-chan string, done chan<- bool) {
	for msg := range link {
		fmt.Println(msg)
	}
	done <- true // 通知主程序消费完成
}
