package main

import (
	"errors"
	"log"
	"net/http"
	"net/rpc"
)

// 定义传入参数和返回参数的数据结构
type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

// 定义服务对象
type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

// 实现 RPC 服务器
func main() {
	// 生成 Arith 对象并注册，通过 HTTP 暴露
	rpc.Register(new(Arith))
	// 通过 http.Handle 在预定义路径 /_goRPC_ 上注册处理器
	// 最终会被添加到 net/http 包中的 DefaultServeMux
	rpc.HandleHTTP()

	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
