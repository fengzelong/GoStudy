package main

import (
	"net"
	"testing"
)

func TestSendUDPMessage(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("解析地址失败: %v", err)
	}
	server, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatalf("监听 UDP 失败: %v", err)
	}
	defer server.Close()

	go func() {
		buf := make([]byte, 128)
		n, remote, err := server.ReadFromUDP(buf)
		if err != nil {
			return
		}
		_, _ = server.WriteToUDP(buf[:n], remote)
	}()

	got, _, err := sendUDPMessage(server.LocalAddr().String(), []byte("hello"))
	if err != nil {
		t.Fatalf("发送 UDP 消息失败: %v", err)
	}
	if string(got) != "hello" {
		t.Fatalf("响应 = %q，期望 hello", string(got))
	}
}
