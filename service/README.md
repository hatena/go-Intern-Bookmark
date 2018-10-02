# サービス層
サービス層が公開している `BookmarkApp` は、アプリケーションロジックを表現するinterfaceです。
```go
type BookmarkApp interface {
	FindUserByID(userID uint64) (*model.User, error)
	CreateOrUpdateBookmark(userID uint64, url string, comment string) (*model.Bookmark, error)
	DeleteBookmark(userID uint64, bookmarkID uint64) error
}
```

リポジトリを保持する `bookmarkApp` 構造体に対して実装を行います。
```go
type bookmarkApp struct {
	repo repository.Repository
}

func (app *bookmarkApp) FindUserByID(userID uint64) (*model.User, error) {
	return app.repo.FindUserByID(userID)
}

func (app *bookmarkApp) DeleteBookmark(userID uint64, bookmarkID uint64) error {
	return app.repo.DeleteBookmark(userID, bookmarkID)
}
```
