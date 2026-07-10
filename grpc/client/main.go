package main

import (
	pb "GoStudy/proto/grpc"
	"context"
	"fmt"

	"GoStudy/internal/config"

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
		total, err := callSum(ctx, client, idx*2, (idx*2)+1)
		if err != nil {
			grpclog.Fatalln(err)
		}
		fmt.Println(total)
	}
}

func callSum(ctx context.Context, client pb.GrpcClient, a int32, b int32) (int32, error) {
	req := &pb.SumRequest{A: a, B: b}
	res, err := client.SumFunc(ctx, req)
	if err != nil {
		return 0, err
	}
	return res.Total, nil
}
