package main

import (
	"net"
	"testing"
	"time"
)

func TestEchoUDPOnce(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("解析地址失败: %v", err)
	}
	server, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatalf("监听 UDP 失败: %v", err)
	}
	defer server.Close()

	done := make(chan string, 1)
	go func() {
		buf := make([]byte, 128)
		message, _, err := echoUDPOnce(server, buf)
		if err != nil {
			done <- err.Error()
			return
		}
		done <- message
	}()

	client, err := net.DialUDP("udp", nil, server.LocalAddr().(*net.UDPAddr))
	if err != nil {
		t.Fatalf("连接 UDP 失败: %v", err)
	}
	defer client.Close()
	if _, err := client.Write([]byte("hello udp")); err != nil {
		t.Fatalf("写入 UDP 失败: %v", err)
	}
	if err := client.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatalf("设置读取超时失败: %v", err)
	}
	buf := make([]byte, 128)
	n, err := client.Read(buf)
	if err != nil {
		t.Fatalf("读取 UDP 回显失败: %v", err)
	}
	if string(buf[:n]) != "hello udp" {
		t.Fatalf("回显 = %q，期望 hello udp", string(buf[:n]))
	}
	if got := <-done; got != "hello udp" {
		t.Fatalf("服务端消息 = %q，期望 hello udp", got)
	}
}
