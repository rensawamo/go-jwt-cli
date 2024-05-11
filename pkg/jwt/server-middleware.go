package jwt

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// コンテキストのメタデータからトークンを取り出す
// grpcのサービス関数の ctxのheaderにauthentizationを設定して毎回 読み込む
func (m *Middleware) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.New(codes.Unauthenticated, "no  provided").Err()
	}
	tokens := headers.Get("authorization")
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
