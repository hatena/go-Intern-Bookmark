package resolver

import (
	"context"
	"fmt"

	"github.com/graph-gophers/graphql-go"

	"github.com/hatena/go-Intern-Bookmark/loader"
	"github.com/hatena/go-Intern-Bookmark/model"
)

type bookmarkResolver struct {
	bookmark *model.Bookmark
}

func (b *bookmarkResolver) ID(ctx context.Context) graphql.ID {
	return graphql.ID(fmt.Sprint(b.bookmark.ID))
}

func (b *bookmarkResolver) Comment(ctx context.Context) string {
	return b.bookmark.Comment
}

func (b *bookmarkResolver) User(ctx context.Context) (*userResolver, error) {
	user, err := loader.LoadUser(ctx, b.bookmark.UserID)
	if err != nil {
		return nil, err
	}
	return &userResolver{user: user}, nil
}

func (b *bookmarkResolver) Entry(ctx context.Context) (*entryResolver, error) {
	entry, err := loader.LoadEntry(ctx, b.bookmark.EntryID)
	if err != nil {
		return nil, err
	}
	return &entryResolver{entry: entry}, nil
}
