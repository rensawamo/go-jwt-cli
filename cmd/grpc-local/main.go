package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
)



func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <JWT token>", os.Args[0])
	}
	token := os.Args[1] // JWTトークンをコマンドライン引数から取得

	ctx := context.Background()
	conn, err := grpc.Dial("localhost:8083", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("接続に失敗しました: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	// JWTトークンをメタデータに追加
	// md := metadata.Pairs("authorization", "Bearer "+token)
	md := metadata.Pairs("authorization", token)
	// メタデータをコンテキストに追加
	ctx = metadata.NewOutgoingContext(ctx, md)
	

	// Contact the server and print out its response
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	resp, err := c.SayHello(ctx, &pb.HelloRequest{Name: "world"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	fmt.Println(resp.Message)
}