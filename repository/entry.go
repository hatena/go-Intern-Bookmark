package repository

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/hatena/go-Intern-Bookmark/model"
)

var entryNotFoundError = model.NotFoundError("entry")

func (r *repository) FindEntryByURL(url string) (*model.Entry, error) {
	var entry model.Entry
	err := r.db.Get(
		&entry,
		`SELECT id,url,title FROM entry
			WHERE url = ? LIMIT 1`, url,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entryNotFoundError
		}
		return nil, err
	}
	return &entry, nil
}

func (r *repository) FindEntryByID(id uint64) (*model.Entry, error) {
	var entry model.Entry
	err := r.db.Get(
		&entry,
		`SELECT id,url,title FROM entry
			WHERE id = ? LIMIT 1`, id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entryNotFoundError
		}
		return nil, err
	}
	return &entry, nil
}

func (r *repository) ListEntriesByIDs(entryIDs []uint64) ([]*model.Entry, error) {
	if len(entryIDs) == 0 {
		return nil, nil
	}
	entries := make([]*model.Entry, 0, len(entryIDs))
	query, args, err := sqlx.In(
		`SELECT id,url,title FROM entry
			WHERE id IN (?)`, entryIDs,
	)
	if err != nil {
		return nil, err
	}
	err = r.db.Select(&entries, query, args...)
	return entries, err
}

func (r *repository) CreateEntry(url string, title string) (*model.Entry, error) {
	id, err := r.generateID()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	_, err = r.db.Exec(
		`INSERT INTO entry
			(id, url, title, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?)`,
		id, url, title, now, now,
	)
	if err != nil {
		return nil, err
	}
	return &model.Entry{ID: id, URL: url, Title: title}, nil
}

func (r *repository) ListEntries(offset, limit uint64) ([]*model.Entry, error) {
	entries := make([]*model.Entry, 0, limit)
	err := r.db.Select(
		&entries,
		`SELECT id,url,title FROM entry
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?`,
		limit, offset,
	)
	return entries, err
}

func (r *repository) BookmarkCountsByEntryIds(entryIDs []uint64) (map[uint64]uint64, error) {
	if len(entryIDs) == 0 {
		return nil, nil
	}
	query, args, err := sqlx.In(
		`SELECT entry_id,COUNT(1) FROM bookmark
			WHERE entry_id IN (?)
			GROUP BY entry_id`, entryIDs,
	)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bookmarkCounts := make(map[uint64]uint64, len(entryIDs))
	for rows.Next() {
		var entryID, count uint64
		rows.Scan(&entryID, &count)
		bookmarkCounts[entryID] = count
	}
	return bookmarkCounts, nil
}
