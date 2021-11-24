package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

func main() {
	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		log.Fatal("dialing: ", err)
	}
	Go(client)
	Call(client)
}

// 异步
func Go(client *rpc.Client) {
	args1 := &Args{7, 8}
	var reply int
	multiplyReply := client.Go("Arith.Multiply", args1, &reply, nil)

	args2 := &Args{15, 6}
	var quo Quotient
	divideReply := client.Go("Arith.Divide", args2, &quo, nil)

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	var multiplyReplied, divideReplied bool
	for !multiplyReplied || !divideReplied {
		select {
		case replyCall := <-multiplyReply.Done:
			if err := replyCall.Error; err != nil {
				fmt.Println("Multiply error: ", err)
			} else {
				fmt.Printf("Multiply: %d*%d=%d\n", args1.A, args1.B, reply)
			}
			multiplyReplied = true
		case replyCall := <-divideReply.Done:
			if err := replyCall.Error; err != nil {
				fmt.Println("Divide error:", err)
			} else {
				fmt.Printf("Divide: %d/%d=%d...%d\n", args2.A, args2.B, quo.Quo, quo.Rem)
			}
			divideReplied = true
		case <-ticker.C:
			fmt.Println("tick")
		}
	}
}

func Call(client *rpc.Client) {
	args := &Args{7, 8}
	var reply int
	err := client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("Multiply error:", err)
	}
	fmt.Printf("Synchronous Multiply: %d*%d=%d\n", args.A, args.B, reply)

	args = &Args{15, 6}
	var quo Quotient
	err = client.Call("Arith.Divide", args, &quo)
	if err != nil {
		log.Fatal("Divide error:", err)
	}
	fmt.Printf("Synchronous Divide: %d/%d=%d...%d\n", args.A, args.B, quo.Quo, quo.Rem)
}
