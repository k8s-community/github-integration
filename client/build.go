package client

import (
	"encoding/json"
	"fmt"
)

const (
	buildCallbackURLStr = "/build-cb"
)

// Possible GitHub Build states
const (
	StatePending = "pending"
	StateSuccess = "success"
	StateError   = "error"
	StateFailure = "failure"
)

const (
	ContextCICD = "k8s-community/cicd"
)

// BuildService defines
type BuildService struct {
	client *Client
}

// BuildCallback defines
type BuildCallback struct {
	Username    string `json:"username"`
	Repository  string `json:"repository"`
	CommitHash  string `json:"commitHash"`
	State       string `json:"state"`
	BuildURL    string `json:"buildURL"`
	Description string `json:"description"`
	Context     string `json:"context"`
}

// BuildCallback sends request for update commit status on github side
func (u *BuildService) BuildCallback(build BuildCallback) error {
	req, err := u.client.NewRequest(postMethod, buildCallbackURLStr, build)
	if err != nil {
		return err
	}

	_, err = u.client.Do(req, nil)
	if err != nil {
		requestBody, _ := json.Marshal(build)
		return fmt.Errorf("couldn't process request: %v, request body:'%s'", err, requestBody)
	}

	return nil
}
