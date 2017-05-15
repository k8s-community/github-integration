package jenkins

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

func NewClient(httpClient *http.Client, baseUrl string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, err := url.Parse(baseUrl)
	if err != nil {
		fmt.Printf("jenkins client: cannot parse url %s: %s", baseUrl, err)
		return nil, err
	}

	c := &Client{
		client:  httpClient,
		BaseURL: baseURL,
	}

	return c, nil
}

// https://www.nczonline.net/blog/2015/10/triggering-jenkins-builds-by-url/
func (c *Client) RunBuild(jobName string, token string) error {
	URL := c.BaseURL.String() + runBuildURL + "/" + jobName + "?token=" + token

	resp, err := c.client.Get(URL)
	if err != nil {
		return fmt.Errorf("jenkins client: cannot run jenkins job, error: %s", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("jenkins client: cannot run jenkins job, http error code = %d", resp.StatusCode)
	}

	return nil
}
