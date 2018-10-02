package model

type Entry struct {
	ID    uint64 `db:"id"`
	URL   string `db:"url"`
	Title string `db:"title"`
}
