package jwt

import (
	"fmt"
	"net/http"
	"strings"
)

// ミドルウェアは、jwtの解析と検証をすべて自動的に行う
type Middleware struct {
		// バリデータを埋め込んでトークン呼び出しをクリーンにする
	Validator
}

// NewMiddlewareは、指定された公開鍵ファイルを使用して検証を行う新しいミドルウェアを作成
// 与えられた公開鍵ファイル
func NewMiddleware(publicKeyPath string) (*Middleware, error) {
	validator, err := NewValidator(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to  validator: %w", err)
	}

	return &Middleware{
		Validator: *validator,
	}, nil
}

func (m *Middleware) HandleHTTP(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.Header.Get("Authorization"), " ")
		if len(parts) < 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("missing or invalid authorization header"))
			return
		}
		tokenString := parts[1]

		token, err := m.GetToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid token: " + err.Error())) 
			return
		}

		// 解析されたトークンで新しいコンテキストを取得する
		ctx := ContextWithToken(r.Context(), token)

		fmt.Println("middleware validated and set set token")

		// 更新されたコンテキストで次のハンドラを呼び出す
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}
