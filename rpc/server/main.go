package main

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"

	"GoStudy/internal/config"
	"GoStudy/rpc/arith"
)

func main() {
	fmt.Println("rpc server start")

	if err := rpc.Register(new(arith.Arith)); err != nil {
		log.Panicln(err)
	}
	rpc.HandleHTTP()

	addr := config.Env("RPC_ADDR", ":8080")
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Panicln(err)
	}
}
