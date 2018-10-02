package web

//go:generate go-assets-builder --package=web --output=./templates-gen.go --strip-prefix="/templates/" --variable=Templates ../templates

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/justinas/nosurf"

	"github.com/hatena/go-Intern-Bookmark/loader"
	"github.com/hatena/go-Intern-Bookmark/model"
	"github.com/hatena/go-Intern-Bookmark/resolver"
	"github.com/hatena/go-Intern-Bookmark/service"
)

type Server interface {
	Handler() http.Handler
}

const sessionKey = "BOOKMARK_SESSION"

var templates map[string]*template.Template

func init() {
	var err error
	templates, err = loadTemplates()
	if err != nil {
		panic(err)
	}
}

func loadTemplates() (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)
	bs, err := ioutil.ReadAll(Templates.Files["main.tmpl"])
	if err != nil {
		return nil, err
	}
	mainTmpl := template.Must(template.New("main.tmpl").Parse(string(bs)))
	for fileName, file := range Templates.Files {
		bs, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		mainTmpl := template.Must(mainTmpl.Clone())
		templates[fileName] = template.Must(mainTmpl.New(fileName).Parse(string(bs)))
	}
	return templates, nil
}

func NewServer(app service.BookmarkApp) Server {
	return &server{app: app}
}

type server struct {
	app service.BookmarkApp
}

func (s *server) Handler() http.Handler {
	router := httptreemux.New()
	handle := func(method, path string, handler http.Handler) {
		router.UsingContext().Handler(method, path,
			csrfMiddleware(loggingMiddleware(headerMiddleware(handler))),
		)
	}

	handle("GET", "/", s.indexHandler())
	handle("GET", "/signup", s.willSignupHandler())
	handle("POST", "/signup", s.signupHandler())
	handle("GET", "/signin", s.willSigninHandler())
	handle("POST", "/signin", s.signinHandler())
	handle("POST", "/signout", s.signoutHandler())
	handle("GET", "/entries/:id", s.entryHandler())
	handle("GET", "/bookmarks", s.bookmarksHandler())
	handle("GET", "/bookmarks/add", s.willAddBookmarkHandler())
	handle("POST", "/bookmarks/add", s.addBookmarkHandler())
	handle("POST", "/bookmarks/:id/delete", s.deleteBookmarksHandler())

	handle("GET", "/spa/", s.spaHandler())
	handle("GET", "/spa/*", s.spaHandler())

	handle("GET", "/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	handle("GET", "/graphiql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templates["graphiql.tmpl"].ExecuteTemplate(w, "graphiql.tmpl", nil)
	}))
	router.UsingContext().Handler("POST", "/query",
		s.attatchLoaderMiddleware(s.resolveUserMiddleware(
			loggingMiddleware(headerMiddleware(resolver.NewHandler(s.app))),
		)),
	)

	return router
}

func (s *server) getParams(r *http.Request, name string) string {
	return httptreemux.ContextParams(r.Context())[name]
}

var csrfMiddleware = func(next http.Handler) http.Handler {
	return nosurf.New(next)
}

var csrfToken = func(r *http.Request) string {
	return nosurf.Token(r)
}

func (s *server) renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["CSRFToken"] = csrfToken(r)
	err := templates[tmpl].ExecuteTemplate(w, "main.tmpl", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *server) findUser(r *http.Request) (user *model.User) {
	cookie, err := r.Cookie(sessionKey)
	if err == nil && cookie.Value != "" {
		user, _ = s.app.FindUserByToken(cookie.Value)
	}
	return
}

// Middleware for fetch user from session key for GraphQL
func (s *server) resolveUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.findUser(r)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", user)))
	})
}

// Middleware for attaching data loaders for GraphQL
func (s *server) attatchLoaderMiddleware(next http.Handler) http.Handler {
	loaders := loader.New(s.app)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(loaders.Attach(r.Context())))
	})
}

func (s *server) indexHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.findUser(r)
		entries, err := s.app.ListEntries(1, 10)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		entryIDs := make([]uint64, len(entries))
		for i, entry := range entries {
			entryIDs[i] = entry.ID
		}
		bookmarkCounts, err := s.app.BookmarkCountsByEntryIds(entryIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		type EntryWithCount struct {
			*model.Entry
			Count uint64
		}
		entryWithCount := make([]EntryWithCount, len(entries))
		for i, e := range entries {
			entryWithCount[i] = EntryWithCount{Entry: e, Count: bookmarkCounts[e.ID]}
		}
		s.renderTemplate(w, r, "index.tmpl", map[string]interface{}{
			"User":    user,
			"Entries": entryWithCount,
		})
	})
}

func (s *server) willSignupHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.renderTemplate(w, r, "signup.tmpl", nil)
	})
}

func (s *server) signupHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, password := r.FormValue("name"), r.FormValue("password")
		if err := s.app.CreateNewUser(name, password); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user, err := s.app.FindUserByName(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		expiresAt := time.Now().Add(24 * time.Hour)
		token, err := s.app.CreateNewToken(user.ID, expiresAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    sessionKey,
			Value:   token,
			Expires: expiresAt,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}

func (s *server) willSigninHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.renderTemplate(w, r, "signin.tmpl", nil)
	})
}

func (s *server) signinHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, password := r.FormValue("name"), r.FormValue("password")
		if ok, err := s.app.LoginUser(name, password); err != nil || !ok {
			http.Error(w, "user not found or invalid password", http.StatusBadRequest)
			return
		}
		user, err := s.app.FindUserByName(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		expiresAt := time.Now().Add(24 * time.Hour)
		token, err := s.app.CreateNewToken(user.ID, expiresAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    sessionKey,
			Value:   token,
			Expires: expiresAt,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}

func (s *server) signoutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    sessionKey,
			Value:   "",
			Expires: time.Unix(0, 0),
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}

func (s *server) entryHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entryID, err := strconv.ParseUint(s.getParams(r, "id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid entry id", http.StatusBadRequest)
			return
		}
		user := s.findUser(r)
		entry, err := s.app.FindEntryByID(entryID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		bookmarks, err := s.app.ListBookmarksByEntryID(entryID, 1, 10)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userIDs := make([]uint64, len(bookmarks))
		for i, bookmark := range bookmarks {
			userIDs[i] = bookmark.UserID
		}
		users, err := s.app.ListUsersByIDs(userIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		type BookmarkWithUser struct {
			Bookmark *model.Bookmark
			User     *model.User
		}
		bookmarkWithUsers := make([]BookmarkWithUser, len(bookmarks))
		for i, bookmark := range bookmarks {
			bookmarkWithUsers[i].Bookmark = bookmark
			for _, user := range users {
				if bookmark.UserID == user.ID {
					bookmarkWithUsers[i].User = user
					break
				}
			}
		}
		s.renderTemplate(w, r, "entry.tmpl", map[string]interface{}{
			"User":      user,
			"Entry":     entry,
			"Bookmarks": bookmarkWithUsers,
		})
	})
}

func (s *server) bookmarksHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.findUser(r)
		if user == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		bookmarks, err := s.app.ListBookmarksByUserID(user.ID, 1, 10)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		entryIDs := make([]uint64, len(bookmarks))
		for i, bookmark := range bookmarks {
			entryIDs[i] = bookmark.EntryID
		}
		entries, err := s.app.ListEntriesByIDs(entryIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		type BookmarkWithEntry struct {
			Bookmark *model.Bookmark
			Entry    *model.Entry
		}
		bookmarkWithEntries := make([]BookmarkWithEntry, len(bookmarks))
		for i, bookmark := range bookmarks {
			bookmarkWithEntries[i].Bookmark = bookmark
			for _, entry := range entries {
				if bookmark.EntryID == entry.ID {
					bookmarkWithEntries[i].Entry = entry
					break
				}
			}
		}
		s.renderTemplate(w, r, "bookmarks.tmpl", map[string]interface{}{
			"User":      user,
			"Bookmarks": bookmarkWithEntries,
		})
	})
}

func (s *server) willAddBookmarkHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.findUser(r)
		if user == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		s.renderTemplate(w, r, "add.tmpl", map[string]interface{}{
			"User": user,
		})
	})
}

func (s *server) addBookmarkHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.findUser(r)
		if user == nil {
			http.Error(w, "please login", http.StatusBadRequest)
			return
		}
		url, comment := r.FormValue("url"), r.FormValue("comment")
		if _, err := s.app.CreateOrUpdateBookmark(user.ID, url, comment); err != nil {
			http.Error(w, "failed to create bookmark", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/bookmarks", http.StatusSeeOther)
	})
}

func (s *server) deleteBookmarksHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.findUser(r)
		if user == nil {
			http.Error(w, "please login", http.StatusBadRequest)
			return
		}
		bookmarkID, err := strconv.ParseUint(s.getParams(r, "id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid bookmark id", http.StatusBadRequest)
			return
		}
		if err := s.app.DeleteBookmark(user.ID, bookmarkID); err != nil {
			http.Error(w, fmt.Sprintf("failed to delete bookmark: %+v", err), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/bookmarks", http.StatusSeeOther)
	})
}

func (s *server) spaHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templates["spa.tmpl"].ExecuteTemplate(w, "spa.tmpl", nil)
	})
}
