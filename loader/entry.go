package loader

import (
	"context"
	"errors"
	"fmt"

	"github.com/graph-gophers/dataloader"

	"github.com/hatena/go-Intern-Bookmark/model"
	"github.com/hatena/go-Intern-Bookmark/service"
)

const entryLoaderKey = "entryLoader"

type entryIDKey struct {
	id uint64
}

func (key entryIDKey) String() string {
	return fmt.Sprint(key.id)
}

func (key entryIDKey) Raw() interface{} {
	return key.id
}

func LoadEntry(ctx context.Context, id uint64) (*model.Entry, error) {
	ldr, err := getLoader(ctx, entryLoaderKey)
	if err != nil {
		return nil, err
	}
	data, err := ldr.Load(ctx, entryIDKey{id: id})()
	if err != nil {
		return nil, err
	}
	return data.(*model.Entry), nil
}

func newEntryLoader(app service.BookmarkApp) dataloader.BatchFunc {
	return func(ctx context.Context, entryIDKeys dataloader.Keys) []*dataloader.Result {
		results := make([]*dataloader.Result, len(entryIDKeys))
		entryIDs := make([]uint64, len(entryIDKeys))
		for i, key := range entryIDKeys {
			entryIDs[i] = key.(entryIDKey).id
		}
		entrys, _ := app.ListEntriesByIDs(entryIDs)
		for i, entryID := range entryIDs {
			results[i] = &dataloader.Result{Data: nil, Error: nil}
			for _, entry := range entrys {
				if entryID == entry.ID {
					results[i].Data = entry
					continue
				}
			}
			if results[i].Data == nil {
				results[i].Error = errors.New("entry not found")
			}
		}
		return results
	}
}
