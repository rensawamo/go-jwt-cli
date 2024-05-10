package jwt

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt"
)

// middlewareContextKey は、コンテキストの値が一意であることを保証するためのカスタム型
type middlewareCtxKey string

// tokenContextKeyは、解析されたトークンに使用されるキー
const tokenCtxKey middlewareCtxKey = "token"

// ContextWithTokenは、与えられたトークンを与えられたコンテキストに追加
func ContextWithToken(ctx context.Context, token *jwt.Token) context.Context {
	return context.WithValue(ctx, tokenCtxKey, token)
}

// ContextGetToken はコンテキストからトークンを取得
// トークンが見つからないか無効な場合はエラーを返す
// トークンの請求や署名の検証はしない
// これは公開鍵を必要とするため、トークンを設定したプロセスで処理される
// トークンを最初に設定したプロセスによって処理される
func ContextGetToken(ctx context.Context) (*jwt.Token, error) {
	val := ctx.Value(tokenCtxKey)
	if val == nil {
		return nil, errors.New("no token found in context")
	}

	t, ok := val.(*jwt.Token)
	if !ok {
		return nil, errors.New("unexpected token type in context")
	}

	return t, nil
}

func MustContextGetToken(ctx context.Context) *jwt.Token {
	t, err := ContextGetToken(ctx)
	if err != nil {
		panic(err)
	}

	return t
}
