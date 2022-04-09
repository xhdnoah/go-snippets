package main

import (
	"fmt"
	"time"

	"github.com/go-zookeeper/zk"
)

var (
	path = "/hello"
)

func create(conn *zk.Conn) {
	var flags int32 = 0
	acls := zk.WorldACL(zk.PermAll) // 获取访问控制权限
	// 创建持久节点
	s, err := conn.Create(path, []byte("world"), flags, acls)
	if err != nil {
		println("failed to create: %v\n", err)
	}
	println("create %s successfully", s)

	// 创建临时节点，创建此节点的会话结束后立即清除此节点
	ephemeral, err := conn.Create("/ephemeral", []byte("1"), zk.FlagEphemeral, acls)
	if err != nil {
		panic(err)
	}
	println("Ephemeral node created:", ephemeral)
}

func get(conn *zk.Conn) {
	res, _, err := conn.Get(path)
	if err != nil {
		panic(err)
	}
	fmt.Println("result: ", string(res)) // "world"
}

// 删改与增不同在于其函数中的 version 参数，用于 CAS 支持以保证原子性
func set(conn *zk.Conn) {
	_, state, _ := conn.Get(path)
	_, err := conn.Set(path, []byte("alice"), state.Version)
	if err != nil {
		panic(err)
	}

	data, _, _ := conn.Get(path)
	fmt.Println("\nnew value: ", string(data)) // "alice"
}

func del(conn *zk.Conn) {
	exists, state, err := conn.Exists(path)
	fmt.Printf("\npath[%s] exists: %v\n", path, exists) // true

	err = conn.Delete(path, state.Version)
	if err != nil {
		panic(err)
	}
	fmt.Printf("path[%s] is deleted", path) // path[/hello] is deleted

	exists, _, err = conn.Exists(path)
	fmt.Printf("\npath[%s] exists: %v\n", path, exists) // false
}

func main() {
	hosts := []string{"172.17.0.2", "172.17.0.3", "172.17.0.4"}
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	create(conn)
	get(conn)
	set(conn)
	del(conn)
}
