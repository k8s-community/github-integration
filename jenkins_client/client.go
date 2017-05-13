package jenkins_client

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	runBuildURL = "/buildWithParameters"
)

type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL
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

	return c
}

// https://www.nczonline.net/blog/2015/10/triggering-jenkins-builds-by-url/
func (c *Client) RunBuild(jobName string, token string) error {
	URL := runBuildURL + "/" + jobName + "?token=" + token

	resp, err := c.client.Get(URL)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("request error, code = %s", resp.StatusCode)
	}

	return nil
}
