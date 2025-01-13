package httpclientpool

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type HttpClient struct {
	Name    string
	Client  *http.Client
	Cookies []*http.Cookie
	Headers map[string]string
}

type HttpClientPool struct {
	clients map[string]*HttpClient
	mu      sync.Mutex
}

var httpClientPool = &HttpClientPool{
	clients: make(map[string]*HttpClient),
}

func NewHttpClient(name string) *HttpClient {
	if httpClientPool.clients[name] != nil {
		return httpClientPool.clients[name]
	}
	client := &HttpClient{
		Name:    name,
		Client:  http.DefaultClient,
		Cookies: []*http.Cookie{},
		Headers: map[string]string{},
	}
	httpClientPool.clients[name] = client
	return client
}

func GetHttpClient(name string) *HttpClient {
	return httpClientPool.clients[name]
}

func (c *HttpClient) Post(link string, encodedData string) (*http.Response, error) {
	req, err := http.NewRequest("POST", link, strings.NewReader(encodedData))
	if err != nil {
		return nil, err
	}
	c.Headers["Content-Length"] = strconv.Itoa(len(encodedData))

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}
	for _, cookie := range c.Cookies {
		req.AddCookie(cookie)
	}

	return c.Client.Do(req)
}

func (c *HttpClient) Get(link string, encodedData string) (*http.Response, error) {
	req, err := http.NewRequest("GET", link, strings.NewReader(encodedData))
	req.URL.RawQuery = encodedData
	if err != nil {
		return nil, err
	}
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}
	for _, cookie := range c.Cookies {
		req.AddCookie(cookie)
	}
	return c.Client.Do(req)
}

func (c *HttpClient) AddHeader(key string, value string) {
	c.Headers[key] = value
}

func (c *HttpClient) SetHeaders(headers map[string]string) {
	c.Headers = headers
}

func (c *HttpClient) SetCookies(cookies []*http.Cookie) {
	c.Cookies = cookies
}

func (c *HttpClient) AddCookie(cookie *http.Cookie) {
	c.Cookies = append(c.Cookies, cookie)
}
