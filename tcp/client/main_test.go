package main

import (
	"net"
	"testing"
)

func TestSendTCPMessage(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("监听失败: %v", err)
	}
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		buf := make([]byte, 64)
		n, _ := conn.Read(buf)
		_, _ = conn.Write(buf[:n])
	}()

	got, err := sendTCPMessage(listener.Addr().String(), "hello")
	if err != nil {
		t.Fatalf("发送 TCP 消息失败: %v", err)
	}
	if got != "hello" {
		t.Fatalf("响应 = %q，期望 hello", got)
	}
}
