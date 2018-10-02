package service

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/hatena/go-Intern-Bookmark/config"
	"github.com/hatena/go-Intern-Bookmark/repository"
	"github.com/hatena/go-Intern-Bookmark/titleFetcher"
)

func newApp() BookmarkApp {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}
	repo, err := repository.New(conf.DbDsn)
	if err != nil {
		panic(err)
	}
	return NewApp(repo, titleFetcher.New())
}

func closeApp(app BookmarkApp) {
	err := app.Close()
	if err != nil {
		panic(err)
	}
}

func randomString() string {
	return strconv.FormatInt(time.Now().Unix()^rand.Int63(), 16)
}
