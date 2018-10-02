package loader

import (
	"context"
	"errors"
	"fmt"

	"github.com/graph-gophers/dataloader"

	"github.com/hatena/go-Intern-Bookmark/model"
	"github.com/hatena/go-Intern-Bookmark/service"
)

const bookmarkLoaderKey = "bookmarkLoader"

type bookmarkIDKey struct {
	id uint64
}

func (key bookmarkIDKey) String() string {
	return fmt.Sprint(key.id)
}

func (key bookmarkIDKey) Raw() interface{} {
	return key.id
}

func LoadBookmark(ctx context.Context, id uint64) (*model.Bookmark, error) {
	ldr, err := getLoader(ctx, bookmarkLoaderKey)
	if err != nil {
		return nil, err
	}
	data, err := ldr.Load(ctx, bookmarkIDKey{id: id})()
	if err != nil {
		return nil, err
	}
	return data.(*model.Bookmark), nil
}

func LoadBookmarksByEntryID(ctx context.Context, id uint64) ([]*model.Bookmark, error) {
	ldr, err := getLoader(ctx, bookmarkLoaderKey)
	if err != nil {
		return nil, err
	}
	data, err := ldr.Load(ctx, entryIDKey{id: id})()
	if err != nil {
		return nil, err
	}
	return data.([]*model.Bookmark), nil
}

func LoadBookmarksByUserID(ctx context.Context, id uint64) ([]*model.Bookmark, error) {
	ldr, err := getLoader(ctx, bookmarkLoaderKey)
	if err != nil {
		return nil, err
	}
	data, err := ldr.Load(ctx, userIDKey{id: id})()
	if err != nil {
		return nil, err
	}
	return data.([]*model.Bookmark), nil
}

func newBookmarkLoader(app service.BookmarkApp) dataloader.BatchFunc {
	return func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		results := make([]*dataloader.Result, len(keys))
		bookmarkIDs := make([]uint64, 0, len(keys))
		entryIDs := make([]uint64, 0, len(keys))
		userIDs := make([]uint64, 0, len(keys))
		for _, key := range keys {
			switch key := key.(type) {
			case bookmarkIDKey:
				bookmarkIDs = append(bookmarkIDs, key.id)
			case entryIDKey:
				entryIDs = append(entryIDs, key.id)
			case userIDKey:
				userIDs = append(userIDs, key.id)
			}
		}
		bookmarks, _ := app.ListBookmarksByIDs(bookmarkIDs)
		bookmarksByEntryIDs, _ := app.ListBookmarksByEntryIDs(entryIDs)
		bookmarksByUserIDs, _ := app.ListBookmarksByUserIDs(userIDs)
		for i, key := range keys {
			results[i] = &dataloader.Result{Data: nil, Error: nil}
			switch key := key.(type) {
			case bookmarkIDKey:
				for _, bookmark := range bookmarks {
					if key.id == bookmark.ID {
						results[i].Data = bookmark
						continue
					}
				}
				if results[i].Data == nil {
					results[i].Error = errors.New("bookmark not found")
				}
			case entryIDKey:
				results[i].Data = bookmarksByEntryIDs[key.id]
			case userIDKey:
				results[i].Data = bookmarksByUserIDs[key.id]
			}
		}
		return results
	}
}
