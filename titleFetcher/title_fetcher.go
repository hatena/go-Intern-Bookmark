package titleFetcher

import (
	"bufio"
	"errors"
	"net/http"
	"strings"
)

type TitleFetcher interface {
	Fetch(url string) (string, error)
}

type titleFetcher struct {
}

func New() TitleFetcher {
	return &titleFetcher{}
}

func (tf *titleFetcher) Fetch(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.Index(line, "<title>")
		if i == -1 {
			continue
		}
		line = line[i+7:]
		j := strings.Index(line, "</title>")
		if j == -1 {
			return "", errors.New("unmatching title tag")
		}
		return line[:j], nil
	}
	return "", errors.New("title not found")
}
