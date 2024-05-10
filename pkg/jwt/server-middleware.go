package jwt

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (m *Middleware) UnaryServerInterceptor(ctx context.Context,  info *grpc.UnaryServerInfo, req interface{},handler grpc.UnaryHandler) (resp interface{}, err error) {
		// コンテキストのメタデータからトークンを取り出す
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.New(codes.Unauthenticated, "no  provided").Err()
	}
	tokens := headers.Get("jwtTokens")
	if len(tokens) < 1 {
		return nil, status.New(codes.Unauthenticated, "no rovided").Err()
	}
	// 最初だけを使用し、繰り返されるヘッダーは無視する
	tokenString := tokens[0] 

	token, err := m.GetToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// トークンをコンテキストに保存し、後で使用できるようにする。
	ctx = ContextWithToken(ctx, token)

	fmt.Println("validation ok and set token")

	// コンテキストを更新して、次のハンドラーを呼び出す。
	return handler(ctx, req)
}
