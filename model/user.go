package model

type User struct {
	ID   uint64 `db:"id"`
	Name string `db:"name"`
}
