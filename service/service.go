package service

import (
	"math/rand"
	"time"

	"github.com/hatena/go-Intern-Bookmark/model"
	"github.com/hatena/go-Intern-Bookmark/repository"
	"github.com/hatena/go-Intern-Bookmark/titleFetcher"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type BookmarkApp interface {
	Close() error

	CreateNewUser(name string, passwordHash string) error
	FindUserByName(name string) (*model.User, error)
	FindUserByID(userID uint64) (*model.User, error)
	ListUsersByIDs(userIDs []uint64) ([]*model.User, error)
	LoginUser(name string, password string) (bool, error)
	CreateNewToken(userID uint64, expiresAt time.Time) (string, error)
	FindUserByToken(token string) (*model.User, error)

	FindOrCreateEntry(url string) (*model.Entry, error)
	ListEntries(page uint64, limit uint64) ([]*model.Entry, error)
	FindEntryByID(entryID uint64) (*model.Entry, error)
	ListEntriesByIDs(entryIDs []uint64) ([]*model.Entry, error)
	BookmarkCountsByEntryIds(entryIDs []uint64) (map[uint64]uint64, error)

	CreateOrUpdateBookmark(userID uint64, url string, comment string) (*model.Bookmark, error)
	FindBookmarkByID(bookmarkIDs uint64) (*model.Bookmark, error)
	ListBookmarksByIDs(bookmarkIDs []uint64) ([]*model.Bookmark, error)
	ListBookmarksByUserID(userID uint64, page uint64, limit uint64) ([]*model.Bookmark, error)
	ListBookmarksByUserIDs(userIDs []uint64) (map[uint64][]*model.Bookmark, error)
	ListBookmarksByEntryID(entryID uint64, page uint64, limit uint64) ([]*model.Bookmark, error)
	ListBookmarksByEntryIDs(entryIDs []uint64) (map[uint64][]*model.Bookmark, error)
	DeleteBookmark(userID uint64, bookmarkID uint64) error
}

func NewApp(repo repository.Repository, titleFetcher titleFetcher.TitleFetcher) BookmarkApp {
	return &bookmarkApp{repo, titleFetcher}
}

type bookmarkApp struct {
	repo         repository.Repository
	titleFetcher titleFetcher.TitleFetcher
}

func (app *bookmarkApp) Close() error {
	return app.repo.Close()
}
