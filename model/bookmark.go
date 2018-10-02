package model

type Bookmark struct {
	ID      uint64 `db:"id"`
	UserID  uint64 `db:"user_id"`
	EntryID uint64 `db:"entry_id"`
	Comment string `db:"comment"`
}
