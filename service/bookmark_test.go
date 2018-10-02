package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hatena/go-Intern-Bookmark/model"
)

func createTestUser(app BookmarkApp) *model.User {
	name := "test name " + randomString()
	password := randomString() + randomString()
	err := app.CreateNewUser(name, password)
	if err != nil {
		panic(err)
	}
	user, err := app.FindUserByName(name)
	if err != nil {
		panic(err)
	}
	return user
}

func TestBookmarkApp_CreateOrUpdateBookmark(t *testing.T) {
	app := newApp()
	defer closeApp(app)

	user := createTestUser(app)

	url := "https://example.com/" + randomString()
	comment := "ブックマークコメント " + randomString()
	bookmark, err := app.CreateOrUpdateBookmark(user.ID, url, comment)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, bookmark.UserID)
	assert.Equal(t, comment, bookmark.Comment)

	comment = "新しい ブックマークコメント " + randomString()
	bookmark, err = app.CreateOrUpdateBookmark(user.ID, url, comment)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, bookmark.UserID)
	assert.Equal(t, comment, bookmark.Comment)

	b, err := app.FindBookmarkByID(bookmark.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, b.UserID)
	assert.Equal(t, comment, b.Comment)
}

func TestBookmarkApp_ListBookmarksByUserID(t *testing.T) {
	app := newApp()
	defer closeApp(app)

	user := createTestUser(app)

	url := "https://example.com/" + randomString()
	comment := "ブックマークコメント " + randomString()
	bookmark, err := app.CreateOrUpdateBookmark(user.ID, url, comment)

	bookmarks, err := app.ListBookmarksByUserID(user.ID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, bookmarks, 1)

	err = app.DeleteBookmark(user.ID, bookmark.ID)
	assert.NoError(t, err)

	bookmarks, err = app.ListBookmarksByUserID(user.ID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, bookmarks, 0)
}
