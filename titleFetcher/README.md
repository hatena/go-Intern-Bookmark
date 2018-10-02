# Title fetcher
ブックマークエントリーを作るときに、URLにアクセスしてタイトルを取ってきます。

```go
type TitleFetcher interface {
	Fetch(url string) (string, error)
}
```
