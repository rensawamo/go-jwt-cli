package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"

	pkgjwt "jwt/pkg/jwt"
)

func main() {
	// 指定された公開鍵パスを使用してミドルウェアを作成する。
	middleware, err := pkgjwt.NewMiddleware(os.Args[1])
	if err != nil {
		panic(err)
	}
	// httpを終了する 省略

	// クライアントを作成し、クライアントインターセプター（ミドルウェア）を追加する
	// そうすることで、自動的にコンテキストからトークンが渡される
	conn, err := grpc.Dial("localhost:8083",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middleware.UnaryClientInterceptor),
	)
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}
	defer conn.Close()

	backendClient := pb.NewGreeterClient(conn)

	// フロントエンドがキーを必要としないことに注意
	frontend, err := NewFrontend(backendClient)
	if err != nil {
		panic(err)
	}


	mux := http.NewServeMux()
	mux.HandleFunc("/claims", frontend.ClaimsHandler)
	mux.HandleFunc("/hello", frontend.HelloHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok!!\n")) 
	})

	root := http.NewServeMux()

	root.Handle("/", middleware.HandleHTTP(mux))

	fmt.Println("Listening on :8082")
	err = http.ListenAndServe(":8082", root)
	if err != nil {
		panic(err)
	}
}

type Frontend struct {
	backendClient pb.GreeterClient
}

func NewFrontend(backendClient pb.GreeterClient) (*Frontend, error) {
	return &Frontend{
		backendClient: backendClient,
	}, nil
}

func (f *Frontend) ClaimsHandler(w http.ResponseWriter, r *http.Request) {
	token := pkgjwt.ShouldContextGetToken(r.Context())

	_, _ = w.Write([]byte(fmt.Sprint(token.Claims)))
}

func (f *Frontend) HelloHandler(w http.ResponseWriter, r *http.Request) {
	preferredName := r.URL.Query().Get("preferredName")
	if preferredName == "" {
		preferredName = "hello"
	}

	resp, err := f.backendClient.SayHello(r.Context(), &pb.HelloRequest{Name: preferredName})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not greet: %v", err))) 
		return
	}

	w.Write([]byte(fmt.Sprintf("Greeting: %s", resp.GetMessage()))) 
}
