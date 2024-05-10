package jwt

import (
	"context"
	"fmt"

	//https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// meta dataにトークンを設定して サーバに送る
func (m *Middleware) UnaryClientInterceptor(ctx context.Context,  req, reply interface{},method string, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// コンテキストから tokenを取得する
	token, err := ContextGetToken(ctx)
	if err != nil {
		return fmt.Errorf("token not set in context: %w", err)
	}

	// grpcコンテキストに認証トークンを追加する。
	ctx = metadata.NewOutgoingContext(ctx,
		metadata.New(
			map[string]string{
				"jwt": token.Raw,
			},
		),
	)

	fmt.Println("* gRPC CLIENT set token")

	// grpc実行

	return invoker(ctx, method, req, reply, cc, opts...)
}
