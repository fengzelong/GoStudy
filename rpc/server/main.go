package main

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"

	"GoStudy/rpc/arith"
)

func main() {
	fmt.Println("rpc server start")

	if err := rpc.Register(new(arith.Arith)); err != nil {
		log.Panicln(err)
	}
	rpc.HandleHTTP()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Panicln(err)
	}
}
