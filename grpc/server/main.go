package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

// 定义helloService并实现约定的接口
type helloService struct{}

// HelloService Hello服务
var HelloService = helloService{}

//func SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
//	resp := new(pb.HelloResponse)
//	resp.Message = fmt.Sprintf("Hello %s.", in.Name)
//
//	return resp, nil
//}

func main() {
	//fmt.Println("git commit")
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}

	// 实例化grpc Server
	s := grpc.NewServer()

	// 注册HelloService
	//pb.RegisterHelloServer(s, HelloService)

	grpclog.Println("Listen on " + Address)
	s.Serve(listen)
}
