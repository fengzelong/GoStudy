package main

import (
	"bufio"
	"fmt"
	"net"

	"GoStudy/internal/config"
)

func process(conn net.Conn) {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("收到 client 端发来的数据:", recvStr)
		conn.Write([]byte(recvStr))
	}
}

func main() {
	addr := config.Env("TCP_ADDR", "127.0.0.1:20000")
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	defer listen.Close()

	fmt.Println("TCP server listen on", addr)
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn)
	}
}
