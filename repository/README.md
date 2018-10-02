# リポジトリ層
このパッケージでは、データベースにアクセスするリポジトリ層を実装します。

`Repository` は、リポジトリ層が提供するインターフェースです。サービス層からタグジャンプするとここにたどり着くので、引数に名前をつけることをおすすめします。
```go
type Repository interface {
	FindUserByID(id uint64) (*User, error)
	ListBookmarksByUserID(userID uint64, offset, limit uint64) ([]*Bookmark, error)
}

func New(dsn string) (Repository, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Opening mysql failed: %v", err)
	}
	return &repository{db: db}, nil
}
```

[database/sqlパッケージのドキュメント](https://golang.org/pkg/database/sql/) や [sqlxパッケージのドキュメント](https://godoc.org/github.com/jmoiron/sqlx)を参考にして、クエリを発行するコードを実装します。
```go
var userNotFoundError = model.NotFoundError("user")

func (r *repository) FindUserByID(id uint64) (*User, error) {
	var user model.User
	err := r.db.Get(
		&user,
		`SELECT id,name FROM user
			WHERE id = ? LIMIT 1`, id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, userNotFoundError
		}
		return nil, err
	}
	return &user, nil
}
```
