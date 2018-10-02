package repository

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/hatena/go-Intern-Bookmark/model"
)

var bookmarkNotFoundError = model.NotFoundError("bookmark")

func (r *repository) FindBookmark(userID uint64, entryID uint64) (*model.Bookmark, error) {
	var bookmark model.Bookmark
	err := r.db.Get(
		&bookmark,
		`SELECT id,user_id,entry_id,comment FROM bookmark
			WHERE user_id = ? AND entry_id = ? LIMIT 1`,
		userID, entryID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, bookmarkNotFoundError
		}
		return nil, err
	}
	return &bookmark, nil
}

func (r *repository) FindBookmarkByID(bookmarkID uint64) (*model.Bookmark, error) {
	var bookmark model.Bookmark
	err := r.db.Get(
		&bookmark,
		`SELECT id,user_id,entry_id,comment FROM bookmark
			WHERE id = ? LIMIT 1`,
		bookmarkID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, bookmarkNotFoundError
		}
		return nil, err
	}
	return &bookmark, nil
}

func (r *repository) CreateBookmark(userID uint64, entryID uint64, comment string) (*model.Bookmark, error) {
	id, err := r.generateID()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	_, err = r.db.Exec(
		`INSERT INTO bookmark
			(id, user_id, entry_id, comment, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?)`,
		id, userID, entryID, comment, now, now,
	)
	if err != nil {
		return nil, err
	}
	return &model.Bookmark{ID: id, UserID: userID, EntryID: entryID, Comment: comment}, nil
}

func (r *repository) UpdateBookmark(id uint64, comment string) error {
	now := time.Now()
	_, err := r.db.Exec(
		`UPDATE bookmark SET comment = ?, updated_at = ?
			WHERE id = ?`,
		comment, now, id,
	)
	return err
}

func (r *repository) ListBookmarksByIDs(bookmarkIDs []uint64) ([]*model.Bookmark, error) {
	if len(bookmarkIDs) == 0 {
		return nil, nil
	}
	bookmarks := make([]*model.Bookmark, 0, len(bookmarkIDs))
	query, args, err := sqlx.In(
		`SELECT id,user_id,entry_id,comment FROM bookmark
			WHERE id IN (?)`, bookmarkIDs,
	)
	if err != nil {
		return nil, err
	}
	err = r.db.Select(&bookmarks, query, args...)
	if err != nil {
		return nil, err
	}
	return bookmarks, err
}

func (r *repository) ListBookmarksByUserID(userID uint64, offset, limit uint64) ([]*model.Bookmark, error) {
	bookmarks := make([]*model.Bookmark, 0, limit)
	err := r.db.Select(
		&bookmarks,
		`SELECT id,user_id,entry_id,comment FROM bookmark
			WHERE user_id = ?
			ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		userID, limit, offset,
	)
	return bookmarks, err
}

func (r *repository) ListBookmarksByUserIDs(userIDs []uint64) (map[uint64][]*model.Bookmark, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	query, args, err := sqlx.In(
		`SELECT id,user_id,entry_id,comment FROM bookmark
			WHERE user_id IN (?)
			ORDER BY created_at DESC`,
		userIDs,
	)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bookmarks := make(map[uint64][]*model.Bookmark)
	for rows.Next() {
		var bookmark model.Bookmark
		rows.Scan(&bookmark.ID, &bookmark.UserID, &bookmark.EntryID, &bookmark.Comment)
		bookmarks[bookmark.UserID] = append(bookmarks[bookmark.UserID], &bookmark)
	}
	return bookmarks, nil
}

func (r *repository) ListBookmarksByEntryID(entryID uint64, offset, limit uint64) ([]*model.Bookmark, error) {
	bookmarks := make([]*model.Bookmark, 0, limit)
	err := r.db.Select(
		&bookmarks,
		`SELECT id,user_id,entry_id,comment FROM bookmark
			WHERE entry_id = ?
			ORDER BY updated_at DESC LIMIT ? OFFSET ?`,
		entryID, limit, offset,
	)
	return bookmarks, err
}

func (r *repository) ListBookmarksByEntryIDs(entryIDs []uint64) (map[uint64][]*model.Bookmark, error) {
	if len(entryIDs) == 0 {
		return nil, nil
	}
	query, args, err := sqlx.In(
		`SELECT id,user_id,entry_id,comment FROM bookmark
			WHERE entry_id IN (?)
			ORDER BY updated_at DESC`,
		entryIDs,
	)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bookmarks := make(map[uint64][]*model.Bookmark)
	for _, entryID := range entryIDs {
		bookmarks[entryID] = make([]*model.Bookmark, 0)
	}
	for rows.Next() {
		var bookmark model.Bookmark
		rows.Scan(&bookmark.ID, &bookmark.UserID, &bookmark.EntryID, &bookmark.Comment)
		bookmarks[bookmark.EntryID] = append(bookmarks[bookmark.EntryID], &bookmark)
	}
	return bookmarks, nil
}

func (r *repository) DeleteBookmark(userID uint64, id uint64) error {
	_, err := r.db.Exec(
		`DELETE FROM bookmark WHERE user_id = ? AND id = ?`,
		userID, id,
	)
	return err
}
