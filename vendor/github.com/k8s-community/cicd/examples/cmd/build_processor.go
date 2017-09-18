package main

import (
	"flag"

	"github.com/Sirupsen/logrus"
	"github.com/k8s-community/cicd/builder"
	"github.com/satori/go.uuid"
)

var (
	fPrefix  = flag.String("prefix", "github.com", "Source code storage (to deal with 'go get')")
	fUser    = flag.String("user", "rumyantseva", "Username (part of path to repo)")
	fRepo    = flag.String("repo", "myapp", "Repository name")
	fCommit  = flag.String("commit", "develop", "Commit hash or branch name")
	fVersion = flag.String("commit", "", "Version to deploy")
)

// This example doesn't deal with API, it just calls processing
func main() {
	flag.Parse()
	log := logrus.New()
	task := builder.NewTask(
		func(state string, description string) { log.Info(state, description) },
		uuid.NewV4().String(),
		"test",
		*fPrefix,
		*fUser,
		*fRepo,
		*fCommit,
		*fVersion,
	)
	builder.Process(log, task)
}
