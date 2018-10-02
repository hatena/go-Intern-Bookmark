# config
アプリケーションを動かすのに必要な設定を、環境変数から読み込みます。

```go
type Config struct {
	Port  int
	DbDsn string
}

func Load() (*Config, error) {
	// ...
}
```
