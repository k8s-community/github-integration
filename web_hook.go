package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/vsaveliev/github-integration/jenkins_client"
	"github.com/vsaveliev/github-integration/user_manager_client"
	"gopkg.in/rjz/githubhook.v0"
	"net/http"
	"strings"
)

func webHookHandler(w http.ResponseWriter, r *http.Request) {
	secret := []byte("asdasd")
	hook, err := githubhook.Parse(secret, r)
	if err != nil {
		fmt.Println("handler: cannot parse hook - ", err)
		return
	}

	switch hook.Event {
	case "integration_installation":
		// Triggered when an integration has been installed or uninstalled.
		fmt.Print("Initialization web hook")
		// to do nothing
	case "integration_installation_repositories":
		// Triggered when a repository is added or removed from an installation.
		fmt.Print("User repositories initialization web hook")
		err = initialUserManagement(hook)
	case "push":
		// Any Git push to a Repository, including editing tags or branches.
		// Commits via API actions that update references are also counted. This is the default event.
		fmt.Print("Push web hook")
		err = runCiCdProcess(hook)
	}

	if err != nil {
		fmt.Println("cannot process web hook", err)
		return
	}
}

func initialUserManagement(hook *githubhook.Hook) error {
	// need to choose event
	evt := github.IntegrationInstallationRepositoriesEvent{}

	err := hook.Extract(&evt)
	if err != nil {
		return err
	}

	fmt.Print("Login: ", *evt.Sender.Login)
	fmt.Printf("Added repositories: %+v", evt.RepositoriesAdded)
	fmt.Printf("Removed repositories: %+v", evt.RepositoriesRemoved)

	// TODO: move to env
	userManagerURL := "http://user.vsaveliev.com"

	client, err := user_manager_client.NewClient(nil, userManagerURL)
	if err != nil {
		return err
	}

	user := user_manager_client.NewUser(*evt.Sender.Login)
	err = client.User.Create(*user)
	if err != nil {
		return err
	}

	for _, rep := range evt.RepositoriesAdded {
		arr := strings.Split(*rep.FullName, "/")
		repository := user_manager_client.NewRepository(arr[0], arr[1])

		// TODO: handle error
		client.Repository.Create(*repository)
	}

	for _, rep := range evt.RepositoriesRemoved {
		arr := strings.Split(*rep.FullName, "/")
		repository := user_manager_client.NewRepository(arr[0], arr[1])

		// TODO: handle error
		client.Repository.Delete(*repository)
	}

	return nil
}

func runCiCdProcess(hook *githubhook.Hook) error {
	// need to choose event
	evt := github.PushEvent{}

	err := hook.Extract(&evt)
	if err != nil {
		return err
	}

	fmt.Print("Login: ", *evt.Sender.Login)

	// TODO: move to env
	jenkinsURL := "http://user.vsaveliev.com"
	token := "k8s-community"

	client, err := jenkins_client.NewClient(nil, jenkinsURL)
	if err != nil {
		return err
	}

	// TODO: replace "/" --> "_"  , for example, "username_repName"
	jobName := *evt.Repo.FullName

	return client.RunBuild(jobName, token)
}
