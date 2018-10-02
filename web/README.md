# web層
このパッケージでは、webサーバーのルーティングや認証、テンプレートのレンダリングなどを行います。

リクエストを受け取ってレスポンスを返すのがwebサーバーです。
Goの `http` パッケージは、webサーバーの動作を表す `Handler` interfaceを提供しています。
```go
type Handler interface {
  ServeHTTP(ResponseWriter, *Request)
}
```

この `http.Handler` を返すサーバーを実装します。
```go
type Server interface {
	Handler() http.Handler
}

func NewServer(app service.BookmarkApp) Server {
	return &server{app: app}
}

type server struct {
	app service.BookmarkApp
}

func (s *server) Handler() http.Handler {
	// ...
}
```

HTTPミドルウェアは、 `http.Handler` を引数にとり `http.Handler` を返す形で実装できます。

例えばリクエストのログを出力するミドルウェアは、次のように実装できます (middleware.goの実装は、ステータスコードを表示するためにもう少し複雑になっています)。
```go
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s took %.2fmsec", r.Method, r.URL.Path,
			float64(time.Now().Sub(start).Nanoseconds())/1e6)
	})
}
```
