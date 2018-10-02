package resolver

import (
	"context"
	"errors"
	"strconv"

	"github.com/hatena/go-Intern-Bookmark/model"
	"github.com/hatena/go-Intern-Bookmark/service"
)

type Resolver interface {
	Visitor(context.Context) (*userResolver, error)
	GetUser(context.Context, struct{ UserID string }) (*userResolver, error)

	CreateBookmark(context.Context, struct{ URL, Comment string }) (*bookmarkResolver, error)
	GetBookmark(context.Context, struct{ BookmarkID string }) (*bookmarkResolver, error)
	DeleteBookmark(context.Context, struct{ BookmarkID string }) (bool, error)

	GetEntry(context.Context, struct{ EntryID string }) (*entryResolver, error)
	ListEntries(context.Context) ([]*entryResolver, error)
}

func newResolver(app service.BookmarkApp) Resolver {
	return &resolver{app: app}
}

type resolver struct {
	app service.BookmarkApp
}

func currentUser(ctx context.Context) *model.User {
	return ctx.Value("user").(*model.User)
}

func (r *resolver) Visitor(ctx context.Context) (*userResolver, error) {
	return &userResolver{currentUser(ctx)}, nil
}

func (r *resolver) GetUser(ctx context.Context, args struct{ UserID string }) (*userResolver, error) {
	userID, err := strconv.ParseUint(args.UserID, 10, 64)
	if err != nil {
		return nil, err
	}
	user, err := r.app.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return &userResolver{user}, nil
}

func (r *resolver) CreateBookmark(ctx context.Context, args struct{ URL, Comment string }) (*bookmarkResolver, error) {
	user := currentUser(ctx)
	if user == nil {
		return nil, errors.New("user not found")
	}
	bookmark, err := r.app.CreateOrUpdateBookmark(user.ID, args.URL, args.Comment)
	if err != nil {
		return nil, err
	}
	return &bookmarkResolver{bookmark: bookmark}, nil
}

func (r *resolver) GetBookmark(ctx context.Context, args struct{ BookmarkID string }) (*bookmarkResolver, error) {
	bookmarkID, err := strconv.ParseUint(args.BookmarkID, 10, 64)
	if err != nil {
		return nil, err
	}
	bookmark, err := r.app.FindBookmarkByID(bookmarkID)
	if err != nil {
		return nil, err
	}
	if bookmark == nil {
		return nil, errors.New("bookmark not found")
	}
	return &bookmarkResolver{bookmark: bookmark}, nil
}

func (r *resolver) DeleteBookmark(ctx context.Context, args struct{ BookmarkID string }) (bool, error) {
	user := currentUser(ctx)
	if user == nil {
		return false, errors.New("user not found")
	}
	bookmarkID, err := strconv.ParseUint(args.BookmarkID, 10, 64)
	if err != nil {
		return false, err
	}
	err = r.app.DeleteBookmark(user.ID, bookmarkID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *resolver) GetEntry(ctx context.Context, args struct{ EntryID string }) (*entryResolver, error) {
	entryID, err := strconv.ParseUint(args.EntryID, 10, 64)
	if err != nil {
		return nil, err
	}
	entry, err := r.app.FindEntryByID(entryID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, errors.New("entry not found")
	}
	return &entryResolver{entry: entry}, nil
}

func (r *resolver) ListEntries(ctx context.Context) ([]*entryResolver, error) {
	entries, err := r.app.ListEntries(1, 10)
	if err != nil {
		return nil, err
	}
	ers := make([]*entryResolver, len(entries))
	for i, entry := range entries {
		ers[i] = &entryResolver{entry: entry}
	}
	return ers, nil
}
