package web

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hatena/go-Intern-Bookmark/config"
	"github.com/hatena/go-Intern-Bookmark/model"
	"github.com/hatena/go-Intern-Bookmark/repository"
	"github.com/hatena/go-Intern-Bookmark/service"
	"github.com/hatena/go-Intern-Bookmark/titleFetcher"
)

func init() {
	csrfMiddleware = func(next http.Handler) http.Handler {
		return next
	}
	csrfToken = func(r *http.Request) string {
		return ""
	}
}

func newAppServer() (service.BookmarkApp, *httptest.Server) {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}
	repo, err := repository.New(conf.DbDsn)
	if err != nil {
		panic(err)
	}
	app := service.NewApp(repo, titleFetcher.New())
	handler := NewServer(app).Handler()
	return app, httptest.NewServer(handler)
}

func randomString() string {
	return strconv.FormatInt(time.Now().Unix()^rand.Int63(), 16)
}

func createTestUser(app service.BookmarkApp) *model.User {
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

func TestServer_Index(t *testing.T) {
	_, testServer := newAppServer()
	defer testServer.Close()

	resp, respBody := client.Get(testServer.URL + "/").Do()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, respBody, `<h1>ブックマーク</h1>`)
	assert.Contains(t, respBody, `<a href="/signup">ユーザー登録</a>`)
}

func TestServer_Signup(t *testing.T) {
	app, testServer := newAppServer()
	defer testServer.Close()

	resp, respBody := client.Get(testServer.URL + "/signup").Do()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, respBody, `<h1>ユーザー登録</h1>`)

	name, password := "test name "+randomString(), randomString()
	resp, _ = client.Post(testServer.URL+"/signup", map[string]string{
		"name":     name,
		"password": password,
	}).Do()
	location := resp.Header.Get("Location")

	assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
	assert.Equal(t, "/", location)

	loginSuccess, _ := app.LoginUser(name, password)
	assert.Equal(t, true, loginSuccess)
}

func TestServer_Signin(t *testing.T) {
	app, testServer := newAppServer()
	defer testServer.Close()

	resp, respBody := client.Get(testServer.URL + "/signin").Do()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, respBody, `<h1>ログイン</h1>`)

	name, password := "test name "+randomString(), randomString()
	err := app.CreateNewUser(name, password)
	assert.NoError(t, err)
	resp, _ = client.Post(testServer.URL+"/signin", map[string]string{
		"name":     name,
		"password": password,
	}).Do()
	location := resp.Header.Get("Location")

	assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
	assert.Equal(t, "/", location)
}

func TestServer_Signout(t *testing.T) {
	app, testServer := newAppServer()
	defer testServer.Close()

	user := createTestUser(app)
	expiresAt := time.Now().Add(24 * time.Hour)
	token, _ := app.CreateNewToken(user.ID, expiresAt)
	sessionCookie := &http.Cookie{Name: sessionKey, Value: token, Expires: expiresAt}

	resp, respBody := client.Get(testServer.URL + "/").WithCookie(sessionCookie).Do()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, respBody, "ユーザー名: "+user.Name)
	assert.Contains(t, respBody, `<input type="submit" value="ログアウト"/>`)

	resp, _ = client.Post(testServer.URL+"/signout", nil).WithCookie(sessionCookie).Do()
	location := resp.Header.Get("Location")

	assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
	assert.Equal(t, "/", location)
	var cookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == sessionKey {
			cookie = c
		}
	}
	assert.Equal(t, "", cookie.Value)
}

func TestServer_Entry(t *testing.T) {
	app, testServer := newAppServer()
	defer testServer.Close()

	user := createTestUser(app)
	url := "https://example.com/" + randomString()
	comment := "ブックマークのコメント" + randomString()
	app.CreateOrUpdateBookmark(user.ID, url, comment)
	entry, _ := app.FindOrCreateEntry(url)

	resp, respBody := client.Get(
		testServer.URL + "/entries/" + fmt.Sprint(entry.ID),
	).Do()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, respBody, user.Name)
	assert.Contains(t, respBody, comment)
}

func TestServer_AddDeleteBookmark(t *testing.T) {
	app, testServer := newAppServer()
	defer testServer.Close()

	user := createTestUser(app)
	expiresAt := time.Now().Add(24 * time.Hour)
	token, _ := app.CreateNewToken(user.ID, expiresAt)
	sessionCookie := &http.Cookie{Name: sessionKey, Value: token, Expires: expiresAt}

	resp, respBody := client.Get(
		testServer.URL + "/bookmarks/add",
	).WithCookie(sessionCookie).Do()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, respBody, `<h1>ブックマーク追加</h1>`)

	comment := "ブックマークのコメント" + randomString()
	resp, _ = client.Post(testServer.URL+"/bookmarks/add", map[string]string{
		"url":     "https://example.com/" + randomString(),
		"comment": comment,
	}).WithCookie(sessionCookie).Do()
	location := resp.Header.Get("Location")

	assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
	assert.Equal(t, "/bookmarks", location)

	bookmarks, _ := app.ListBookmarksByUserID(user.ID, 1, 10)
	assert.Equal(t, comment, bookmarks[0].Comment)
	bookmark := bookmarks[0]

	resp, respBody = client.Get(
		testServer.URL + "/bookmarks",
	).WithCookie(sessionCookie).Do()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, respBody, comment)
	assert.Contains(t, respBody, `<input type="submit" value="削除"/>`)

	resp, _ = client.Post(
		testServer.URL+"/bookmarks/"+fmt.Sprint(bookmark.ID)+"/delete", nil,
	).WithCookie(sessionCookie).Do()
	location = resp.Header.Get("Location")

	assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
	assert.Equal(t, "/bookmarks", location)

	bookmarks, _ = app.ListBookmarksByUserID(user.ID, 1, 10)
	assert.Len(t, bookmarks, 0)
}
