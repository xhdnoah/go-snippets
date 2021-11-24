package main

import (
	"fmt"
	"net"
	"net/http"
)

type Header struct{}

func (h *Header) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {
	l, _ := net.Listen("tcp", ":8080")
	http.Handle("/headers", &Header{})
	_ = http.Serve(l, nil)
}
