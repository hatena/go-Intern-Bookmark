package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServer_Graphql_User(t *testing.T) {
	app, testServer := newAppServer()
	defer testServer.Close()

	user := createTestUser(app)
	expiresAt := time.Now().Add(24 * time.Hour)
	token, _ := app.CreateNewToken(user.ID, expiresAt)
	sessionCookie := &http.Cookie{Name: sessionKey, Value: token, Expires: expiresAt}

	resp, respBody := client.PostJSON(testServer.URL+"/query", map[string]interface{}{
		"query": fmt.Sprintf(
			`query {
				visitor() {
					id,
					name
				}
			}`,
		),
	}).WithCookie(sessionCookie).Do()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var visitorData struct {
		Data struct {
			Visitor struct {
				ID   string
				Name string
			}
		}
	}
	assert.NoError(t, json.Unmarshal([]byte(respBody), &visitorData))
	assert.Equal(t, fmt.Sprint(user.ID), visitorData.Data.Visitor.ID)
	assert.Equal(t, user.Name, visitorData.Data.Visitor.Name)
}
