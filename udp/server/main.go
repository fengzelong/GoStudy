package main

import (
	"fmt"
	"net"

	"GoStudy/internal/config"
)

func echoUDPOnce(listen *net.UDPConn, data []byte) (string, *net.UDPAddr, error) {
	n, remoteAddr, err := listen.ReadFromUDP(data)
	if err != nil {
		return "", nil, err
	}
	message := string(data[:n])
	_, err = listen.WriteToUDP(data[:n], remoteAddr)
	if err != nil {
		return "", remoteAddr, err
	}
	return message, remoteAddr, nil
}

func main() {
	addr := config.Env("UDP_ADDR", "127.0.0.1:30000")
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("resolve udp addr failed, err:", err)
		return
	}

	listen, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	defer listen.Close()

	fmt.Println("UDP server listen on", addr)
	for {
		var data [1024]byte
		message, remoteAddr, err := echoUDPOnce(listen, data[:])
		if err != nil {
			fmt.Println("read udp failed, err:", err)
			continue
		}
		fmt.Printf("data:%v addr:%v count:%v\n", message, remoteAddr, len(message))
	}
}
