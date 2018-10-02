package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServer_Graphql_Bookmark(t *testing.T) {
	app, testServer := newAppServer()
	defer testServer.Close()

	user := createTestUser(app)
	expiresAt := time.Now().Add(24 * time.Hour)
	token, _ := app.CreateNewToken(user.ID, expiresAt)
	sessionCookie := &http.Cookie{Name: sessionKey, Value: token, Expires: expiresAt}
	url := "https://example.com/" + randomString()
	comment := "ブックマークのコメント" + randomString()
	resp, respBody := client.PostJSON(testServer.URL+"/query", map[string]interface{}{
		"query": fmt.Sprintf(
			`mutation {
				createBookmark(url: %q, comment: %q) {
					id,
					comment
				}
			}`,
			url, comment,
		),
	}).WithCookie(sessionCookie).Do()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var createBookmarkData struct {
		Data struct {
			CreateBookmark struct {
				ID      string
				Comment string
			}
		}
	}
	assert.NoError(t, json.Unmarshal([]byte(respBody), &createBookmarkData))
	assert.Equal(t, comment, createBookmarkData.Data.CreateBookmark.Comment)

	resp, respBody = client.PostJSON(testServer.URL+"/query", map[string]interface{}{
		"query": fmt.Sprintf(
			`query {
				getBookmark(bookmarkId: %q) {
					id,
					comment,
					user {
						id,
						name
					}
				}
			}`,
			createBookmarkData.Data.CreateBookmark.ID,
		),
	}).Do()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var getBookmarkData struct {
		Data struct {
			GetBookmark struct {
				ID      string
				Comment string
				User    struct {
					ID   string
					Name string
				}
			}
		}
	}
	assert.NoError(t, json.Unmarshal([]byte(respBody), &getBookmarkData))
	assert.Equal(t, comment, getBookmarkData.Data.GetBookmark.Comment)
	assert.Equal(t, fmt.Sprint(user.ID), getBookmarkData.Data.GetBookmark.User.ID)
	assert.Equal(t, user.Name, getBookmarkData.Data.GetBookmark.User.Name)

	resp, respBody = client.PostJSON(testServer.URL+"/query", map[string]interface{}{
		"query": fmt.Sprintf(
			`mutation {
				deleteBookmark(bookmarkId: %q) {}
			}`,
			createBookmarkData.Data.CreateBookmark.ID,
		),
	}).WithCookie(sessionCookie).Do()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var deleteBookmarkData struct{ Data struct{ DeleteBookmark bool } }
	assert.NoError(t, json.Unmarshal([]byte(respBody), &deleteBookmarkData))
	assert.True(t, deleteBookmarkData.Data.DeleteBookmark)
}
