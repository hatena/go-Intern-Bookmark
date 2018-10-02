package service

import (
	"errors"

	"github.com/hatena/go-Intern-Bookmark/model"
)

func (app *bookmarkApp) FindOrCreateEntry(url string) (*model.Entry, error) {
	entry, err := app.repo.FindEntryByURL(url)
	if err != nil {
		if model.IsNotFound(err) {
			title, err := app.titleFetcher.Fetch(url)
			if err != nil {
				title = url
			}
			return app.repo.CreateEntry(url, title)
		}
		return nil, err
	}
	return entry, err
}

func (app *bookmarkApp) ListEntries(page uint64, limit uint64) ([]*model.Entry, error) {
	if page < 1 || limit < 1 {
		return nil, errors.New("page and limit should be positive")
	}
	return app.repo.ListEntries((page-1)*limit, limit)
}

func (app *bookmarkApp) FindEntryByID(entryID uint64) (*model.Entry, error) {
	return app.repo.FindEntryByID(entryID)
}

func (app *bookmarkApp) ListEntriesByIDs(entryIDs []uint64) ([]*model.Entry, error) {
	return app.repo.ListEntriesByIDs(entryIDs)
}

func (app *bookmarkApp) BookmarkCountsByEntryIds(entryIDs []uint64) (map[uint64]uint64, error) {
	return app.repo.BookmarkCountsByEntryIds(entryIDs)
}
