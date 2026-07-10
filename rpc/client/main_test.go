package main

import (
	"net"
	"net/http"
	"net/rpc"
	"testing"

	"GoStudy/rpc/arith"
)

func TestCallArithmetic(t *testing.T) {
	server := rpc.NewServer()
	if err := server.RegisterName("Arith", new(arith.Arith)); err != nil {
		t.Fatalf("注册 RPC 服务失败: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle(rpc.DefaultRPCPath, server)
	mux.Handle(rpc.DefaultDebugPath, server)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("监听失败: %v", err)
	}
	defer listener.Close()

	go http.Serve(listener, mux)

	res, err := callArithmetic(listener.Addr().String(), arith.Request{A: 9, B: 2})
	if err != nil {
		t.Fatalf("RPC 调用失败: %v", err)
	}
	if res.Pro != 18 || res.Quo != 4 || res.Rem != 1 {
		t.Fatalf("响应 = %+v，期望乘积 18、商 4、余数 1", res)
	}
}
