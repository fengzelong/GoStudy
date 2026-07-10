package main

import (
	pb "GoStudy/proto/grpc"
	"testing"

	"golang.org/x/net/context"
)

func TestGrpcServiceSayHello(t *testing.T) {
	resp, err := GrpcSvc.SayHello(context.Background(), &pb.HelloRequest{Name: "Tom", Sex: "male"})
	if err != nil {
		t.Fatalf("SayHello 返回错误: %v", err)
	}
	if resp.Message != "Hello Tom, your sex: male" {
		t.Fatalf("Message = %q", resp.Message)
	}
}

func TestGrpcServiceSumFunc(t *testing.T) {
	resp, err := GrpcSvc.SumFunc(context.Background(), &pb.SumRequest{A: 2, B: 3})
	if err != nil {
		t.Fatalf("SumFunc 返回错误: %v", err)
	}
	if resp.Total != 5 {
		t.Fatalf("Total = %d，期望 5", resp.Total)
	}
}
