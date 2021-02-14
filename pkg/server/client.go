package server

import (
	"io"
	"net/http"
	"time"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36"

type HttpClient struct {
	Client *http.Client
}
var httpClient *HttpClient

// Wrapping http.Client
func NewClient() *HttpClient {
	httpClient = &HttpClient{http.DefaultClient}
	httpClient.Client.Timeout = time.Second * 10

	return httpClient
}

// Wrapping http.Requests
func NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}

func GetHttpClient() *HttpClient {
	c := *httpClient
	return &c
}

func (c *HttpClient) Get(url string) (resp *http.Response, err error) {
	req, err := NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("User-Agent", UserAgent)
	return c.Client.Do(req)
}

func (c *HttpClient) Post(url string, body io.Reader) (resp *http.Response, err error) {
	req, err := NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("User-Agent", UserAgent)
	return c.Client.Do(req)
}
