package cicd_test

import (
	"fmt"

	"github.com/k8s-community/cicd"
)

func ExampleBuild() {
	client := cicd.NewClient("http://127.0.0.1:8080")

	request := &cicd.BuildRequest{
		Username:   "rumyantseva",
		Repository: "myapp",
		CommitHash: "fc6d3deecc2d9b09d69f26dcadcc8dacc6e663dc",
	}

	resp, err := client.Build(request)
	if err != nil {
		fmt.Printf("Server error: %s", err)
	} else if resp.Error != nil {
		fmt.Printf("Client's errors %d: %s", resp.Error.Code, resp.Error.Message)
	} else if resp.Data != nil {
		fmt.Printf("Build was created, request ID is: %s", resp.Data.RequestID)
	} else {
		fmt.Printf("Unknown error")
	}
}
