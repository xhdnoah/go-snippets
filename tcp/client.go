package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
)

var host = flag.String("host", "localhost", "host")
var port = flag.String("port", "8888", "port")

var planets = [...]string{"Mercury", "Venus", "Earth", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune"}

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", *host+":"+*port)

	if err != nil {
		fmt.Println("error connecting: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("connecting to " + *host + ":" + *port)

	var wg sync.WaitGroup
	wg.Add(2)

	go handleWrite(conn, &wg)
	go handleRead(conn, &wg)

	wg.Wait()
}

func handleWrite(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, p := range planets {
		_, err := conn.Write([]byte("Hello " + p + "\r\n"))
		if err != nil {
			fmt.Println("error to send message: ", err.Error())
			break
		}
	}
}

func handleRead(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	reader := bufio.NewReader(conn)
	for i := 1; i <= 10; i++ {
		line, err := reader.ReadString(byte('\n'))
		if err != nil {
			fmt.Println("error to read message: ", err)
			return
		}
		fmt.Println(line)
	}
}
