package main

import (
	"encoding/json"
	"fmt"
	"ggrpc"
	"ggrpc/codec"
	"log"
	"net"
	"time"
)

func startServer(addr chan string) {
	// pick a free port
	listener, err := net.Listen("tcp", ":8085")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", listener.Addr())
	addr <- listener.Addr().String()
	ggrpc.Accept(listener)
}

func main() {
	addr := make(chan string)
	go startServer(addr)

	// in fact, following code is like a simple geerpc client
	conn, _ := net.Dial("tcp", <-addr)
	defer func() {
		conn.Close()
	}()

	time.Sleep(time.Second)
	json.NewEncoder(conn).Encode(ggrpc.DefaultOption) // 将 DefaultOption 编码并写入conn
	cc := codec.NewGobCodec(conn)                     // 创建一个 GobCodec 实例
	// send request & receive response
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		cc.Write(h, fmt.Sprintf("ggrpc req %d", h.Seq)) // 发送请求(header 和 body)
		// server返回的时候，会将header和body一起返回，所以这里需要先读取header，再读取body
		cc.ReadHeader(h)
		var reply string
		cc.ReadBody(&reply) // 读取响应体
		log.Println("reply:", reply)
	}
}
