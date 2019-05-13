package service

import (
	"errors"

	"github.com/hatena/go-Intern-Bookmark/model"
)

func (app *bookmarkApp) CreateOrUpdateBookmark(userID uint64, url string, comment string) (*model.Bookmark, error) {
	if url == "" {
		return nil, errors.New("empty url")
	}
	user, err := app.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	entry, err := app.FindOrCreateEntry(url)
	if err != nil {
		return nil, err
	}
	bookmark, err := app.repo.FindBookmark(user.ID, entry.ID)
	if bookmark != nil {
		if err := app.repo.UpdateBookmark(bookmark.ID, comment); err != nil {
			return nil, err
		}
		return app.repo.FindBookmark(user.ID, entry.ID)
	}
	return app.repo.CreateBookmark(user.ID, entry.ID, comment)
}

func (app *bookmarkApp) FindBookmarkByID(bookmarkID uint64) (*model.Bookmark, error) {
	return app.repo.FindBookmarkByID(bookmarkID)
}

func (app *bookmarkApp) ListBookmarksByIDs(bookmarkIDs []uint64) ([]*model.Bookmark, error) {
	return app.repo.ListBookmarksByIDs(bookmarkIDs)
}

func (app *bookmarkApp) ListBookmarksByUserID(userID uint64, page uint64, limit uint64) ([]*model.Bookmark, error) {
	return app.repo.ListBookmarksByUserID(userID, (page-1)*limit, limit)
}

func (app *bookmarkApp) ListBookmarksByUserIDs(userIDs []uint64) (map[uint64][]*model.Bookmark, error) {
	return app.repo.ListBookmarksByUserIDs(userIDs)
}

func (app *bookmarkApp) ListBookmarksByEntryID(entryID uint64, page uint64, limit uint64) ([]*model.Bookmark, error) {
	return app.repo.ListBookmarksByEntryID(entryID, (page-1)*limit, limit)
}

func (app *bookmarkApp) ListBookmarksByEntryIDs(entryIDs []uint64) (map[uint64][]*model.Bookmark, error) {
	return app.repo.ListBookmarksByEntryIDs(entryIDs)
}

func (app *bookmarkApp) DeleteBookmark(userID uint64, bookmarkID uint64) error {
	return app.repo.DeleteBookmark(userID, bookmarkID)
}
