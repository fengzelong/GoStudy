package main

import (
	pb "GoStudy/proto/grpc"
	"fmt"
	"net"
	"net/http"

	"GoStudy/internal/config"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// GrpcService 实现 proto 中约定的接口。
type GrpcService struct{}

// SayHello 实现 gRPC Hello 服务。
func (g GrpcService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	resp := new(pb.HelloResponse)
	resp.Message = fmt.Sprintf("Hello %s, your sex: %s", in.Name, in.Sex)

	return resp, nil
}

// SumFunc 实现 gRPC 求和服务。
func (g GrpcService) SumFunc(ctx context.Context, in *pb.SumRequest) (*pb.SumResponse, error) {
	resp := new(pb.SumResponse)
	total := in.A + in.B
	resp.Total = total
	return resp, nil
}

// GrpcSvc 是注册到服务端的 gRPC 服务实例。
var GrpcSvc = GrpcService{}

func main() {
	addr := config.Env("GRPC_ADDR", "127.0.0.1:50052")
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}

	svc := grpc.NewServer()
	pb.RegisterGrpcServer(svc, GrpcSvc)

	go startTrace()

	fmt.Println("Listen on " + addr)
	svc.Serve(listen)
}

func startTrace() {
	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		return true, true
	}
	traceAddr := config.Env("GRPC_TRACE_ADDR", ":50051")
	go http.ListenAndServe(traceAddr, nil)
	fmt.Println("Trace listen on " + traceAddr)
}
