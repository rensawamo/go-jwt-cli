package jwt

import (
	"crypto"
	"fmt"
	"io/ioutil"

	"github.com/golang-jwt/jwt"
)

// バリデータはJWTトークンの解析と検証を行う
type Validator struct {
	key crypto.PublicKey
}

// 公開鍵で Validator を作成する
func NewValidator(publicKeyPath string) (*Validator, error) {
	keyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("public key error : %w", err)
	}

	key, err := jwt.ParseEdPublicKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("private key error : %w", err)
	}

	return &Validator{
		key: key,
	}, nil
}

// 、トークン文字列を受け取り、その署名の検証とクレームのバリデーションを行います。無効なトークンの場合はエラーを返す
func (v *Validator) GetToken(tokenString string) (*jwt.Token, error) {

	// key func のおかげで
  // jwt.Parseは署名検証やクレーム検証も行う。
	token, err := jwt.Parse(
		tokenString,
		// トークンが期待される署名メソッドを使用しているかどうかを確認する
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
				return nil, fmt.Errorf("dont match alglizm: %v", token.Header["alglizum"])
			}
			// 署名検証のため公開鍵を返す
			return v.key, nil
		})
	if err != nil {
		return nil, fmt.Errorf("unable to parse: %w", err)
	}

	return token, nil
}
