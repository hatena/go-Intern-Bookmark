package resolver

import (
	"context"
	"fmt"

	"github.com/graph-gophers/graphql-go"

	"github.com/hatena/go-Intern-Bookmark/loader"
	"github.com/hatena/go-Intern-Bookmark/model"
)

type entryResolver struct {
	entry *model.Entry
}

func (e *entryResolver) ID(ctx context.Context) graphql.ID {
	return graphql.ID(fmt.Sprint(e.entry.ID))
}

func (e *entryResolver) URL(ctx context.Context) string {
	return e.entry.URL
}

func (e *entryResolver) Title(ctx context.Context) string {
	return e.entry.Title
}

func (e *entryResolver) Bookmarks(ctx context.Context) ([]*bookmarkResolver, error) {
	bookmarks, err := loader.LoadBookmarksByEntryID(ctx, e.entry.ID)
	if err != nil {
		return nil, err
	}
	brs := make([]*bookmarkResolver, len(bookmarks))
	for i, bookmark := range bookmarks {
		brs[i] = &bookmarkResolver{bookmark: bookmark}
	}
	return brs, nil
}
