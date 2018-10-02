package loader

import (
	"context"
	"errors"
	"fmt"

	"github.com/graph-gophers/dataloader"

	"github.com/hatena/go-Intern-Bookmark/model"
	"github.com/hatena/go-Intern-Bookmark/service"
)

const userLoaderKey = "userLoader"

type userIDKey struct {
	id uint64
}

func (key userIDKey) String() string {
	return fmt.Sprint(key.id)
}

func (key userIDKey) Raw() interface{} {
	return key.id
}

func LoadUser(ctx context.Context, id uint64) (*model.User, error) {
	ldr, err := getLoader(ctx, userLoaderKey)
	if err != nil {
		return nil, err
	}
	data, err := ldr.Load(ctx, userIDKey{id: id})()
	if err != nil {
		return nil, err
	}
	return data.(*model.User), nil
}

func newUserLoader(app service.BookmarkApp) dataloader.BatchFunc {
	return func(ctx context.Context, userIDKeys dataloader.Keys) []*dataloader.Result {
		results := make([]*dataloader.Result, len(userIDKeys))
		userIDs := make([]uint64, len(userIDKeys))
		for i, key := range userIDKeys {
			userIDs[i] = key.(userIDKey).id
		}
		users, _ := app.ListUsersByIDs(userIDs)
		for i, userID := range userIDs {
			results[i] = &dataloader.Result{Data: nil, Error: nil}
			for _, user := range users {
				if userID == user.ID {
					results[i].Data = user
					continue
				}
			}
			if results[i].Data == nil {
				results[i].Error = errors.New("user not found")
			}
		}
		return results
	}
}
