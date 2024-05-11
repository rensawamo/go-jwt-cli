package main

import (
	"context"
	_ "embed"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"

	pkgjwt "jwt/pkg/jwt"

	"github.com/golang-jwt/jwt"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf(" %s <key-path>\n", os.Args[0])
		os.Exit(1)
	}

	backend, err := NewBackend()
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", ":8083")
	if err != nil {
		panic(err)
	}

	// 指定された公開鍵パスを使用してミドルウェアを作成
	middleware, err := pkgjwt.NewMiddleware(os.Args[1])
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(
		// token認証を行うミドルウェアを追加
		grpc.UnaryInterceptor(middleware.UnaryServerInterceptor),
	)
	pb.RegisterGreeterServer(s, backend)

	fmt.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

type Backend struct {
	pb.UnimplementedGreeterServer
}

func NewBackend() (*Backend, error) {
	return &Backend{}, nil
}

// SayHelloはhelloworld.GreeterServerを実装しており、コンテキストに有効なトークンが必要
func (b *Backend) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	token := pkgjwt.ShouldContextGetToken(ctx)

	// クレームからロールを取得
	roles := token.Claims.(jwt.MapClaims)["roles"]

	return &pb.HelloReply{
		Message: fmt.Sprintf("%s! backend your roll is %v", in.GetName(), roles),
	}, nil
}