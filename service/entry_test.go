package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hatena/go-Intern-Bookmark/model"
)

func TestBookmarkApp_FindOrCreateEntry(t *testing.T) {
	app := newApp()
	defer closeApp(app)

	url := "https://example.com/" + randomString()
	entry, err := app.FindOrCreateEntry(url)
	assert.NoError(t, err)
	assert.Equal(t, url, entry.URL)
	assert.Equal(t, "Example Domain", entry.Title)

	entry, err = app.FindOrCreateEntry(url)
	assert.NoError(t, err)
	assert.Equal(t, url, entry.URL)
	assert.Equal(t, "Example Domain", entry.Title)

	url = "http://b.hatena.ne.jp"
	entry, err = app.FindOrCreateEntry(url)
	assert.NoError(t, err)
	assert.Equal(t, url, entry.URL)
	assert.Equal(t, "はてなブックマーク", entry.Title)
}

func TestBookmarkApp_FindEntryByID(t *testing.T) {
	app := newApp()
	defer closeApp(app)

	url := "https://example.com/" + randomString()
	entry, err := app.FindOrCreateEntry(url)
	assert.NoError(t, err)
	assert.Equal(t, url, entry.URL)
	assert.Equal(t, "Example Domain", entry.Title)

	e, err := app.FindEntryByID(entry.ID)
	assert.NoError(t, err)
	assert.Equal(t, e.URL, entry.URL)
	assert.Equal(t, e.Title, entry.Title)
}

func TestBookmarkApp_ListEntries(t *testing.T) {
	app := newApp()
	defer closeApp(app)

	url := "https://example.com/" + randomString()

	for i := 0; i < 10; i++ {
		_, err := app.FindOrCreateEntry(url + fmt.Sprint(i))
		assert.NoError(t, err)
	}

	entries, err := app.ListEntries(1, 3)
	assert.NoError(t, err)
	assert.Len(t, entries, 3)
	assert.Equal(t, url+"9", entries[0].URL)
	assert.Equal(t, url+"8", entries[1].URL)
	assert.Equal(t, url+"7", entries[2].URL)

	entries, err = app.ListEntries(2, 4)
	assert.NoError(t, err)
	assert.Len(t, entries, 4)
	assert.Equal(t, url+"5", entries[0].URL)
	assert.Equal(t, url+"4", entries[1].URL)
}

func TestBookmarkApp_ListEntriesByIDs(t *testing.T) {
	app := newApp()
	defer closeApp(app)

	url := "https://example.com/"
	entry1, err := app.FindOrCreateEntry(url)
	assert.NoError(t, err)

	url = "http://b.hatena.ne.jp"
	entry2, err := app.FindOrCreateEntry(url)
	assert.NoError(t, err)

	url = "http://hatena.ne.jp"
	entry3, err := app.FindOrCreateEntry(url)
	assert.NoError(t, err)

	entries, err := app.ListEntriesByIDs([]uint64{entry1.ID, entry2.ID, entry3.ID})
	assert.NoError(t, err)
	assert.Len(t, entries, 3)
	for _, entry := range []*model.Entry{entry1, entry2, entry3} {
		var contained bool
		for _, e := range entries {
			contained = contained || entry.ID == e.ID
		}
		assert.True(t, contained)
	}

	entries, err = app.ListEntriesByIDs([]uint64{})
	assert.NoError(t, err)
	assert.Len(t, entries, 0)
}
