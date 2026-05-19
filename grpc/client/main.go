package main

import (
	pb "GoStudy/proto/grpc"
	"fmt"

	"GoStudy/internal/config"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
)

func main() {
	addr := config.Env("GRPC_ADDR", "127.0.0.1:50052")
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer conn.Close()

	client := pb.NewGrpcClient(conn)

	var idx int32 = 0
	ctx := context.Background()
	for idx < 50 {
		idx++
		req := &pb.SumRequest{A: idx * 2, B: (idx * 2) + 1}
		res, err := client.SumFunc(ctx, req)
		if err != nil {
			grpclog.Fatalln(err)
		}
		fmt.Println(res.Total)
	}
}
