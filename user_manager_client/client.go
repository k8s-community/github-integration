package user_manager_client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const (
	postMethod   = "POST"
	getMethod    = "GET"
	putMethod    = "PUT"
	deleteMethod = "DELETE"
)

type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// Services used for talking to different parts of the API.
	User       *UserService
	Repository *RepositoryService
}

func NewClient(httpClient *http.Client, baseUrl string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(baseUrl)

	c := &Client{
		client:  httpClient,
		BaseURL: baseURL,
	}

	c.User = &UserService{client: c}
	c.Repository = &RepositoryService{client: c}

	return c
}

func (c *Client) NewRequest(method string, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}
