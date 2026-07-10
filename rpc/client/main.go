package main

import (
	"fmt"
	"log"
	"net/rpc"

	"GoStudy/internal/config"
	"GoStudy/rpc/arith"
)

func callArithmetic(addr string, req arith.Request) (arith.Response, error) {
	conn, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		return arith.Response{}, err
	}
	defer conn.Close()

	var res arith.Response

	if err := conn.Call("Arith.Multiply", req, &res); err != nil {
		return arith.Response{}, err
	}

	if err := conn.Call("Arith.Divide", req, &res); err != nil {
		return arith.Response{}, err
	}
	return res, nil
}

func main() {
	addr := config.Env("RPC_ADDR", "127.0.0.1:8080")
	req := arith.Request{A: 9, B: 2}
	res, err := callArithmetic(addr, req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
	fmt.Printf("%d / %d = %d remainder %d\n", req.A, req.B, res.Quo, res.Rem)
}
