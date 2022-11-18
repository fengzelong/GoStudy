package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
)

// Arith 结构体，用于注册的
type Arith struct{}

// ArithRequest 声明参数结构体
type ArithRequest struct {
	A, B int
}

// ArithResponse 返回给客户端的结果
type ArithResponse struct {
	// 乘积
	Pro int
	// 商
	Quo int
	// 余数
	Rem int
}

// Multiply 乘法
func (a *Arith) Multiply(req ArithRequest, res *ArithResponse) error {
	res.Pro = req.A * req.B
	return nil
}

// Divide 商和余数
func (a *Arith) Divide(req ArithRequest, res *ArithResponse) error {
	if req.B == 0 {
		return errors.New("除数不能为0")
	}
	// 除
	res.Quo = req.A / req.B
	// 取模
	res.Rem = req.A % req.B
	return nil
}

func main() {
	fmt.Println("rpc server start")

	rect := new(Arith)
	rpc.Register(rect)
	rpc.HandleHTTP()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panicln(err)
	}
}
