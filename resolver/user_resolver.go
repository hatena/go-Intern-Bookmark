package resolver

import (
	"context"
	"fmt"

	"github.com/graph-gophers/graphql-go"

	"github.com/hatena/go-Intern-Bookmark/loader"
	"github.com/hatena/go-Intern-Bookmark/model"
)

type userResolver struct {
	user *model.User
}

func (u *userResolver) ID(ctx context.Context) graphql.ID {
	return graphql.ID(fmt.Sprint(u.user.ID))
}

func (u *userResolver) Name(ctx context.Context) string {
	return u.user.Name
}

func (u *userResolver) Bookmarks(ctx context.Context) ([]*bookmarkResolver, error) {
	bookmarks, err := loader.LoadBookmarksByUserID(ctx, u.user.ID)
	if err != nil {
		return nil, err
	}
	bs := make([]*bookmarkResolver, len(bookmarks))
	for i, bookmark := range bookmarks {
		bs[i] = &bookmarkResolver{bookmark: bookmark}
	}
	return bs, nil
}
