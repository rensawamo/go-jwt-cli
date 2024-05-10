
### 秘密鍵と 公開鍵の作成
```sh
$ openssl genpkey -out auth.ed
openssl pkey -in auth.ed -pubout > auth.ed.pub
```

### auth 
```sh
# serverのセットアップ
$ go run ./cmd/auth auth.ed


# ログイン
$ curl -X POST http://localhost:8080/login -H "Authorization: Basic $(echo -n 'admin:pass' | base64)"
eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhcGkiLCJleHAiOjE3MTUzMDg0NDQsImlhdCI6MTcxNTMwODM4NCwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgxIiwibmJmIjoxNzE1MzA4Mzg0LCJyb2xlcyI6WyJhZG1pbiIsImJhc2ljIl0sInVzZX.....
``` 

### frontend 
curlからtokenを受け取り、新しいtokenの認証ありのmidrewareをもつ新しいクライアントに変更して
tokenの検証を行う。そしてtokenが有効な場合、grpcのサービス関数にctx の中に認証の情報を含めて responseを返す

```sh
# serverの セットアップ
$ go run ./cmd/front


# tokenをセットして grpcのサービス関数に tokenを入れ込んで実行
$ token=$(curl admin:pass@localhost:8081/login); echo $token;curl -H "Authorization: Bearer $token" localhost:8082/hello;echo
```


### fronend の ミドルウェア
grpc クライアントを tokenの認証ありのmidrewareをもつ新しいクライアントに変更して
tokenの検証を行う。
そしてそのトークンを mildreware が検証を行い、有効なら ハンドラー関数を実行する

```sh
# serverの セットアップ
$ go run ./cmd/frontend-midleware auth.ed.pub 


# ログインして サービス関数の実行
$  token=$(curl admin:pass@localhost:8081/login); echo $token; curl -H "Authorization: Bearer $token" localhost:8082/hello; echo


# /hello で grpcのtokenがセットされ mildlrewareで認証されて 権限情報を取得
$ curl -H "Authorization: Bearer $token" localhost:8082/hello; echo
```
loginして tokenをうけとり midlewareにセットする。
