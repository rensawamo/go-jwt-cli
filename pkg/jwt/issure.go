package jwt

import (
	"crypto"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// Issuerは、JWTトークンを発行するための構造体
type Issuer struct {
	key crypto.PrivateKey
}

// NewIssuerは、与えられたパスをed25519秘密鍵として解析することで、新しい発行者を作成する
func NewIssuer(privateKeyPath string) (*Issuer, error) {
	keyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		panic(fmt.Errorf("unable to read private key file: %w", err))
	}

	key, err := jwt.ParseEdPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse as ed private key: %w", err)
	}

	return &Issuer{
		key: key,
	}, nil
}

// IssueTokenは、指定されたロールを持つ指定されたユーザのための新しいトークンを発行する
func (i *Issuer) IssueToken(user string, roles []string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, jwt.MapClaims{
		  // 標準的なclaimsを設定する
			"exp": now.Add(time.Minute).Unix(),  // トークンの有効期限を設定（現在時刻から1分後に設定）
			"aud": "api",                        // トークンの受け手を指定（このトークンはapi用）
			"iat": now.Unix(),                   // トークンの発行時刻
			"nbf": now.Unix(),                   // トークンが有効になる時刻（発行時刻と同じ）
			"iss": "http://localhost:8081",      // トークンの発行者の識別子

			// ユーザーに関するカスタムclaimsを追加
			"user": user,                        // 認証されたユーザーの名前

			// ユーザーの権限リストを含める
			// 複雑なデータタイプを持つclaimsを示す
			"roles": roles,
	})

	// シークレットを使用して、エンコードされたトークンを文字列として取得
	tokenString, err := token.SignedString(i.key)
	if err != nil {
		return "", fmt.Errorf("unable sign token: %w", err)
	}

	return tokenString, nil
}
