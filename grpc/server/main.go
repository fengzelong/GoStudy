package main

import (
	pb "GoStudy/proto/grpc"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
	"net/http"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

// GrpcService 实现约定的接口
type GrpcService struct{}

// SayHello grpc服务go实现
func (g GrpcService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	resp := new(pb.HelloResponse)
	resp.Message = fmt.Sprintf("Hello %s, your sex: %s", in.Name, in.Sex)

	return resp, nil
}

// SumFunc grpc服务go实现
func (g GrpcService) SumFunc(ctx context.Context, in *pb.SumRequest) (*pb.SumResponse, error) {
	resp := new(pb.SumResponse)
	total := in.A + in.B
	resp.Total = total
	return resp, nil
}

// GrpcSvc Hello服务
var GrpcSvc = GrpcService{}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}

	// 实例化grpc Server
	svc := grpc.NewServer()

	// 注册HelloService
	pb.RegisterGrpcServer(svc, GrpcSvc)

	// 开启trace
	go startTrace()

	fmt.Println("Listen on " + Address)
	svc.Serve(listen)
}

// startTrace trace
func startTrace() {
	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		return true, true
	}
	go http.ListenAndServe(":50051", nil)
	fmt.Println("Trace listen on 50051")
}
