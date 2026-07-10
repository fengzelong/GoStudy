package main

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"

	"GoStudy/internal/config"
	"GoStudy/rpc/arith"
)

func newArithRPCServer() (*rpc.Server, error) {
	server := rpc.NewServer()
	if err := server.RegisterName("Arith", new(arith.Arith)); err != nil {
		return nil, err
	}
	return server, nil
}

func main() {
	fmt.Println("rpc server start")

	if err := rpc.RegisterName("Arith", new(arith.Arith)); err != nil {
		log.Panicln(err)
	}
	rpc.HandleHTTP()

	addr := config.Env("RPC_ADDR", ":8080")
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Panicln(err)
	}
}
