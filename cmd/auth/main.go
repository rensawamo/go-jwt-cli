package main

import (
	_ "embed"
	"errors"
	"jwt/pkg/jwt"
	"net/http"
	"fmt"
	"os"
)

func main() {
	issuer, err := jwt.NewIssuer(os.Args[1])
	if err != nil {
		panic(err)
	}

	auth, err := NewAuthenticationService(issuer)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/login", auth.HandleLogin)

	fmt.Println("Listening on :8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

// AuthService は認証を処理し、トークンを発行する。
type AuthService struct {
	issuer *jwt.Issuer
}

// NewAuthService は、指定された発行者を使用して新しいサービスを作成します。
func NewAuthenticationService(issuer *jwt.Issuer) (*AuthService, error) {
	if issuer == nil {
		return nil, errors.New("required issue")
	}
	return &AuthService{
		issuer: issuer,
	}, nil
}

func (a *AuthService) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// ユーザー名とパスワードを取得
	user, pass, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("missing auth")) 
		return
	}

	if user != "admin" || pass != "pass" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid credentials")) 
		return
	}

	tokenString, err := a.issuer.IssueToken("admin", []string{"admin", "basic"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable token:" + err.Error())) 
		return
	}
	// shellのtokenへ
	_, _ = w.Write([]byte(tokenString + "\n"))
}
