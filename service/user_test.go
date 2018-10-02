package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBookmarkApp_FindUserByName(t *testing.T) {
	app := newApp()
	defer closeApp(app)

	name := "test name " + randomString()
	password := randomString() + randomString()
	err := app.CreateNewUser(name, password)
	assert.NoError(t, err)

	user, err := app.FindUserByName(name)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, name)

	user2, err := app.FindUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user2.Name, name)
}

func TestBookmarkApp_ListUsersByIDs(t *testing.T) {
	app := newApp()
	defer closeApp(app)

	userIDs := make([]uint64, 5)
	for i := 0; i < 5; i++ {
		name := "test name " + randomString()
		password := randomString() + randomString()
		_ = app.CreateNewUser(name, password)
		user, _ := app.FindUserByName(name)
		userIDs[i] = user.ID
	}

	users, err := app.ListUsersByIDs(userIDs)
	assert.NoError(t, err)
	assert.Len(t, users, 5)
}

func TestBookmarkApp_LoginUser(t *testing.T) {
	app := newApp()
	defer closeApp(app)

	name := "test name " + randomString()
	password := randomString() + randomString()
	err := app.CreateNewUser(name, password)
	assert.NoError(t, err)

	login, err := app.LoginUser(name, password)
	assert.NoError(t, err)
	assert.True(t, login)

	login, err = app.LoginUser(name, password+".")
	assert.NoError(t, err)
	assert.False(t, login)
}

func TestBookmarkApp_CreateNewToken(t *testing.T) {
	app := newApp()
	defer closeApp(app)

	name := "test name " + randomString()
	password := randomString() + randomString()
	err := app.CreateNewUser(name, password)
	assert.NoError(t, err)
	user, _ := app.FindUserByName(name)

	token, err := app.CreateNewToken(user.ID, time.Now().Add(1*time.Hour))
	assert.NoError(t, err)
	assert.NotEqual(t, "", token)

	u, err := app.FindUserByToken(token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, u.ID)
}
