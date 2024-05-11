package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	pkgjwt "jwt/pkg/jwt"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("%s <key-path>\n", os.Args[0])
		os.Exit(1)
	}

	validator, err := pkgjwt.NewValidator(os.Args[1])
	if err != nil {
		panic(err)
	}

	backend, err := NewBackend(validator)
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", ":8083")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, backend)
	reflection.Register(s)

	fmt.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

type Backend struct {
	pb.UnimplementedGreeterServer
	validator *pkgjwt.Validator
}

func NewBackend(validator *pkgjwt.Validator) (*Backend, error) {
	return &Backend{
		validator: validator,
	}, nil
}

func (b *Backend) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	token, err := b.tokenFromContextMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get token: %w", err)
	}

	roles := token.Claims.(jwt.MapClaims)["roles"]

	return &pb.HelloReply{
		Message: fmt.Sprintf(
			"%s!  backend. You have roles %v",
			in.GetName(), roles),
	}, nil
}

func (b *Backend) tokenFromContextMetadata(ctx context.Context) (*jwt.Token, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("found no metadata in context")
	}
	tokens := headers.Get("authorization")
	if len(tokens) < 1 {
		return nil, errors.New("found no token in metadata")
	}
	tokenString := strings.TrimPrefix(tokens[0], "Bearer ")

	token, err := b.validator.GetToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("error : %w", err)
	}

	return token, nil
}














