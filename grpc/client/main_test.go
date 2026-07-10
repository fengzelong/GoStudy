package main

import (
	pb "GoStudy/proto/grpc"
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type testGrpcService struct{}

func (testGrpcService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello " + in.Name}, nil
}

func (testGrpcService) SumFunc(ctx context.Context, in *pb.SumRequest) (*pb.SumResponse, error) {
	return &pb.SumResponse{Total: in.A + in.B}, nil
}

func TestCallSumWithLocalGRPCServer(t *testing.T) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	pb.RegisterGrpcServer(server, testGrpcService{})

	go func() {
		if err := server.Serve(listener); err != nil {
			t.Errorf("启动本地 gRPC server 失败: %v", err)
		}
	}()
	defer server.Stop()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("连接本地 gRPC server 失败: %v", err)
	}
	defer conn.Close()

	total, err := callSum(ctx, pb.NewGrpcClient(conn), 7, 8)
	if err != nil {
		t.Fatalf("callSum 返回错误: %v", err)
	}
	if total != 15 {
		t.Fatalf("total = %d，期望 15", total)
	}
}
