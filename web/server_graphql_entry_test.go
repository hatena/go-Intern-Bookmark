package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_Graphql_Entry(t *testing.T) {
	app, testServer := newAppServer()
	defer testServer.Close()

	url := "https://example.com/" + randomString()
	entry, _ := app.FindOrCreateEntry(url)

	user1 := createTestUser(app)
	comment1 := "ブックマークのコメント" + randomString()
	app.CreateOrUpdateBookmark(user1.ID, url, comment1)

	user2 := createTestUser(app)
	comment2 := "ブックマークのコメント" + randomString()
	app.CreateOrUpdateBookmark(user2.ID, url, comment2)

	resp, respBody := client.PostJSON(testServer.URL+"/query", map[string]interface{}{
		"query": fmt.Sprintf(
			`query {
				getEntry(entryId: %q) {
					id,
					url,
					title,
					bookmarks {
						id,
						comment,
						user {
							id,
							name
						}
					}
				}
			}`,
			fmt.Sprint(entry.ID),
		),
	}).Do()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	type EntryData struct {
		ID        string
		URL       string
		Title     string
		Bookmarks []struct {
			ID      string
			Comment string
			User    struct {
				ID   string
				Name string
			}
		}
	}
	var getEntryData struct {
		Data struct {
			GetEntry EntryData
		}
	}
	assert.NoError(t, json.Unmarshal([]byte(respBody), &getEntryData))
	assert.Equal(t, fmt.Sprint(entry.ID), getEntryData.Data.GetEntry.ID)
	assert.Equal(t, url, getEntryData.Data.GetEntry.URL)
	assert.Equal(t, "Example Domain", getEntryData.Data.GetEntry.Title)
	assert.Len(t, getEntryData.Data.GetEntry.Bookmarks, 2)
	assert.Equal(t, comment2, getEntryData.Data.GetEntry.Bookmarks[0].Comment)
	assert.Equal(t, comment1, getEntryData.Data.GetEntry.Bookmarks[1].Comment)
	assert.Equal(t, user2.Name, getEntryData.Data.GetEntry.Bookmarks[0].User.Name)
	assert.Equal(t, user1.Name, getEntryData.Data.GetEntry.Bookmarks[1].User.Name)

	resp, respBody = client.PostJSON(testServer.URL+"/query", map[string]interface{}{
		"query": fmt.Sprintf(
			`query {
				listEntries() {
					id,
					url,
					title,
					bookmarks {
						id,
						comment
					}
				}
			}`,
		),
	}).Do()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var listEntriesData struct {
		Data struct {
			ListEntries []EntryData
		}
	}
	assert.NoError(t, json.Unmarshal([]byte(respBody), &listEntriesData))
}
