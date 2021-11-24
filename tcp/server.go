package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

var host = flag.String("host", "", "host")
var port = flag.String("port", "8888", "port")

func main() {
	flag.Parse()
	var l net.Listener
	var err error
	l, err = net.Listen("tcp", *host+":"+*port)
	if err != nil {
		fmt.Println("error listening: ", err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("listening on " + *host + ":" + *port)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("error accepting: ", err)
			os.Exit(1)
		}
		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	for {
		io.Copy(conn, conn)
	}
}
