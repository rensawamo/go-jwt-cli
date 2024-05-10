package main

import (
	_ "embed"
	"errors"
	"fmt"
	pkgjwt "jwt/pkg/jwt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
)

func main() {

	validator, err := pkgjwt.NewValidator(os.Args[1])
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial("localhost:8083",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	backendClient := pb.NewGreeterClient(conn)

	frontend, err := NewFrontend(validator, backendClient)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", frontend.RootHandler)
	mux.HandleFunc("/hello", frontend.HelloHandler)
	mux.HandleFunc("/claims", frontend.ClaimsHandler)

	fmt.Println("Listening on :8082")
	err = http.ListenAndServe(":8082", mux)
	if err != nil {
		panic(err)
	}
}

type Frontend struct {
	validator     *pkgjwt.Validator
	backendClient pb.GreeterClient
}

func NewFrontend(validator *pkgjwt.Validator, backendClient pb.GreeterClient) (*Frontend, error) {
	return &Frontend{
		validator:     validator,
		backendClient: backendClient,
	}, nil
}

func (f *Frontend) ClaimsHandler(w http.ResponseWriter, r *http.Request) {
		// トークンを取得して claimsを出力
	token, err := f.getHeaderToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("auth error:" + err.Error())) //nolint
		return
	}

	_, _ = w.Write([]byte(fmt.Sprint(token.Claims)))
}

func (f *Frontend) HelloHandler(w http.ResponseWriter, r *http.Request) {
		// リクエストから好みの名前を取得
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "ren"
	}

	// トークンを取得して、grpcのコンテキストに追加
	token, err := f.getHeaderToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("auth error:" + err.Error())) //nolint
		return
	}
	ctx  := metadata.NewOutgoingContext(
		r.Context(),
		metadata.New(
			map[string]string{
				"jwtTokens": token.Raw,
			},
		),
	)
	// セッター end omit
	// 新しいコンテキストで呼び出しを行う
	resp, err := f.backendClient.SayHello(ctx, &pb.HelloRequest{
		Name: name,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not greet: %v", err))) 
		return
	}

	w.Write([]byte(fmt.Sprintf("Greeting: %s", resp.GetMessage()))) 
}

func (f *Frontend) RootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := f.getHeaderToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("auth error:" + err.Error())) 
		return
	}
	w.Write([]byte("ok\n")) 
}

// getHeaderTokenは、Authorizationヘッダーに有効なJWTトークンがあるかどうかをチェックする。
func (f *Frontend) getHeaderToken(h http.Header) (*jwt.Token, error) {
	auth := strings.Split(h.Get("Authorization"), " ")
	if len(auth) < 2 || auth[0] != "Bearer" {
		return nil, errors.New("invalid header")
	}
	tokenString := auth[1]

	// Parse the token from the header
	token, err := f.validator.GetToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return token, nil
}
