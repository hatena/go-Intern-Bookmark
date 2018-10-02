package repository

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/hatena/go-Intern-Bookmark/model"
)

type Repository interface {
	CreateNewUser(name string, passwordHash string) error
	FindUserByName(name string) (*model.User, error)
	FindUserByID(id uint64) (*model.User, error)
	ListUsersByIDs(userIDs []uint64) ([]*model.User, error)
	FindPasswordHashByName(name string) (string, error)
	CreateNewToken(userID uint64, token string, expiresAt time.Time) error
	FindUserByToken(token string) (*model.User, error)

	FindEntryByURL(url string) (*model.Entry, error)
	FindEntryByID(id uint64) (*model.Entry, error)
	ListEntriesByIDs(entryIDs []uint64) ([]*model.Entry, error)
	CreateEntry(url string, title string) (*model.Entry, error)
	ListEntries(offset, limit uint64) ([]*model.Entry, error)
	BookmarkCountsByEntryIds(entryIDs []uint64) (map[uint64]uint64, error)

	FindBookmark(userID uint64, entryID uint64) (*model.Bookmark, error)
	FindBookmarkByID(bookmarkID uint64) (*model.Bookmark, error)
	CreateBookmark(userID uint64, entryID uint64, comment string) (*model.Bookmark, error)
	UpdateBookmark(id uint64, comment string) error
	ListBookmarksByIDs(bookmarkIDs []uint64) ([]*model.Bookmark, error)
	ListBookmarksByUserID(userID uint64, offset, limit uint64) ([]*model.Bookmark, error)
	ListBookmarksByUserIDs(userIDs []uint64) (map[uint64][]*model.Bookmark, error)
	ListBookmarksByEntryID(entryID uint64, offset, limit uint64) ([]*model.Bookmark, error)
	ListBookmarksByEntryIDs(entryIDs []uint64) (map[uint64][]*model.Bookmark, error)
	DeleteBookmark(userID uint64, id uint64) error
	Close() error
}

func New(dsn string) (Repository, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Opening mysql failed: %v", err)
	}
	return &repository{db: db}, nil
}

type repository struct {
	db *sqlx.DB
}

func (r *repository) generateID() (uint64, error) {
	var id uint64
	err := r.db.Get(&id, "SELECT UUID_SHORT()")
	return id, err
}

func (r *repository) Close() error {
	return r.db.Close()
}
