package main

import (
	"fmt"
	"log"
	"net/rpc"

	"GoStudy/rpc/arith"
)

func main() {
	conn, err := rpc.DialHTTP("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	req := arith.Request{A: 9, B: 2}
	var res arith.Response

	if err := conn.Call("Arith.Multiply", req, &res); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)

	if err := conn.Call("Arith.Divide", req, &res); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d / %d = %d remainder %d\n", req.A, req.B, res.Quo, res.Rem)
}
