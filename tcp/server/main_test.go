package main

import (
	"net"
	"testing"
	"time"
)

func TestProcessEcho(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()

	go process(serverConn)

	if _, err := clientConn.Write([]byte("hello tcp")); err != nil {
		t.Fatalf("写入失败: %v", err)
	}

	if err := clientConn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatalf("设置读取超时失败: %v", err)
	}
	buf := make([]byte, 64)
	n, err := clientConn.Read(buf)
	if err != nil {
		t.Fatalf("读取失败: %v", err)
	}
	if string(buf[:n]) != "hello tcp" {
		t.Fatalf("回显 = %q，期望 hello tcp", string(buf[:n]))
	}
}
