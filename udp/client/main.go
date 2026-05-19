package main

import (
	"fmt"
	"net"

	"GoStudy/internal/config"
)

func main() {
	addr := config.Env("UDP_ADDR", "127.0.0.1:30000")
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("resolve udp addr failed, err:", err)
		return
	}

	socket, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("连接服务端失败，err:", err)
		return
	}
	defer socket.Close()

	sendData := []byte("Hello server")
	_, err = socket.Write(sendData)
	if err != nil {
		fmt.Println("发送数据失败，err:", err)
		return
	}

	data := make([]byte, 4096)
	n, remoteAddr, err := socket.ReadFromUDP(data)
	if err != nil {
		fmt.Println("接收数据失败，err:", err)
		return
	}
	fmt.Printf("recv:%v addr:%v count:%v\n", string(data[:n]), remoteAddr, n)
}
