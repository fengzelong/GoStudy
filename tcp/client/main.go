package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"GoStudy/internal/config"
)

func sendTCPMessage(addr string, message string) (string, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(message)); err != nil {
		return "", err
	}
	buf := [512]byte{}
	n, err := conn.Read(buf[:])
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func main() {
	addr := config.Env("TCP_ADDR", "127.0.0.1:20000")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	defer conn.Close()
	inputReader := bufio.NewReader(os.Stdin)
	for {
		input, _ := inputReader.ReadString('\n')
		inputInfo := strings.Trim(input, "\r\n")
		if strings.ToUpper(inputInfo) == "Q" {
			return
		}
		_, err = conn.Write([]byte(inputInfo))
		if err != nil {
			return
		}
		buf := [512]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("recv failed, err:", err)
			return
		}
		fmt.Println(string(buf[:n]))
	}
}
