package web

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type httpClient struct {
	client *http.Client
	req    *http.Request
}

var client = &httpClient{client: &http.Client{
	CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	},
}}

func (c *httpClient) Get(path string) *httpClient {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		panic(err)
	}
	c.req = req
	return c
}

func (c *httpClient) Post(path string, data map[string]string) *httpClient {
	values := url.Values{}
	for k, v := range data {
		values.Add(k, v)
	}
	req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(values.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(values.Encode())))
	c.req = req
	return c
}

func (c *httpClient) PostJSON(path string, data interface{}) *httpClient {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(jsonBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(jsonBytes)))
	c.req = req
	return c
}

func (c *httpClient) WithCookie(cookie *http.Cookie) *httpClient {
	c.req.AddCookie(cookie)
	return c
}

func (c *httpClient) Do() (*http.Response, string) {
	resp, err := c.client.Do(c.req)
	if err != nil {
		panic(err)
	}
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	respBody := string(respBodyBytes)
	return resp, respBody
}
