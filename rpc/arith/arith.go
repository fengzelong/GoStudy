package arith

import "errors"

// Arith 是服务端注册的 RPC 服务。
type Arith struct{}

// Request 是 RPC 客户端和服务端共用的请求参数。
type Request struct {
	A, B int
}

// Response 是 RPC 客户端和服务端共用的返回结果。
type Response struct {
	Pro int
	Quo int
	Rem int
}

// Multiply 计算 A * B。
func (a *Arith) Multiply(req Request, res *Response) error {
	res.Pro = req.A * req.B
	return nil
}

// Divide 计算商和余数。
func (a *Arith) Divide(req Request, res *Response) error {
	if req.B == 0 {
		return errors.New("除数不能为零")
	}
	res.Quo = req.A / req.B
	res.Rem = req.A % req.B
	return nil
}
