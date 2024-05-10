
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

