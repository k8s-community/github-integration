package cicd

import (
	"fmt"
	"net/http"

	"github.com/k8s-community/cicd/utils/rest"
)

const (
	buildURL = "/api/v1/build"
)

// Client defines REST client
type Client struct {
	client *rest.Client
}

// NewClient creates an instance of the Client
func NewClient(baseURL string) *Client {
	return &Client{
		client: rest.NewClient(nil, baseURL),
	}
}

// Build runs CICD-build. Please, see an ExampleBuild.
func (c *Client) Build(request *BuildRequest) (*BuildResponse, error) {
	req, err := c.client.NewRequest("POST", buildURL, request)
	if err != nil {
		return nil, err
	}

	var response = new(BuildResponse)

	resp, err := c.client.Do(req, response)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		if response.Error != nil {
			return nil, fmt.Errorf("Code %d, %s", response.Error.Code, response.Error.Message)
		}

		return nil, fmt.Errorf("Unknown error from CICD service.")
	}

	return response, nil
}
