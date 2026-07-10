package main

import (
	"fmt"
	"net"

	"GoStudy/internal/config"
)

func sendUDPMessage(addr string, sendData []byte) ([]byte, *net.UDPAddr, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, nil, err
	}

	socket, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, nil, err
	}
	defer socket.Close()

	if _, err := socket.Write(sendData); err != nil {
		return nil, nil, err
	}

	data := make([]byte, 4096)
	n, remoteAddr, err := socket.ReadFromUDP(data)
	if err != nil {
		return nil, nil, err
	}
	return data[:n], remoteAddr, nil
}

func main() {
	addr := config.Env("UDP_ADDR", "127.0.0.1:30000")
	sendData := []byte("Hello server")
	data, remoteAddr, err := sendUDPMessage(addr, sendData)
	if err != nil {
		fmt.Println("接收数据失败，err:", err)
		return
	}
	fmt.Printf("recv:%v addr:%v count:%v\n", string(data), remoteAddr, len(data))
}
