package main

import (
	pb "GoStudy/proto/grpc"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

func main() {
	// 连接
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer conn.Close()

	// 初始化客户端
	client := pb.NewGrpcClient(conn)

	// 调用方法
	//req := &pb.HelloRequest{Name: "lily", Sex: "female"}
	//res, err1 := c.SayHello(context.Background(), req)
	//if err1 != nil {
	//	grpclog.Fatalln(err1)
	//}
	//fmt.Println(res.Message)

	var idx int32 = 0
	ctx := context.Background()
	for idx < 50 {
		idx++
		req := &pb.SumRequest{A: idx * 2, B: (idx * 2) + 1}
		res, err1 := client.SumFunc(ctx, req)
		if err1 != nil {
			grpclog.Fatalln(err1)
		}
		fmt.Println(res.Total)
		//time.Sleep(1 * time.Second)
	}
}
