package main

import (
	"flag"

	"github.com/Sirupsen/logrus"
	"github.com/k8s-community/cicd/builder/runners"
	"github.com/k8s-community/cicd/builder/task"
	"github.com/satori/go.uuid"
)

var (
	fPrefix    = flag.String("prefix", "github.com", "Source code storage (to deal with 'go get')")
	fRepo      = flag.String("repo", "myapp", "Repository name")
	fCommit    = flag.String("commit", "develop", "Commit hash or branch name")
	fVersion   = flag.String("commit", "", "Version to deploy")
	fNamespace = flag.String("namespace", "dev", "Namespace to release to")
)

// This example doesn't deal with API, it just calls processing
func main() {
	flag.Parse()
	log := logrus.New()
	taskItem := task.NewCICD(
		func(taskID, state, description string) { log.Info(taskID, state, description) },
		uuid.NewV4().String(),
		"test",
		*fPrefix,
		*fRepo,
		*fCommit,
		*fVersion,
		*fNamespace,
	)

	runner := runners.NewLocal(log)
	runner.Process(*taskItem)
}
