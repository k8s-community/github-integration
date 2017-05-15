package main

import (
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"github.com/vsaveliev/github-integration/jenkins"
	userManClient "github.com/vsaveliev/user-manager/client"
	"gopkg.in/rjz/githubhook.v0"
)

func (h *handler) webHookHandler(w http.ResponseWriter, r *http.Request) {
	secret := []byte(h.env["GITHUBINT_TOKEN"])
	hook, err := githubhook.Parse(secret, r)
	if err != nil {
		h.errlog.Printf("Cannot parse hook (ID %s): %s", hook.Id, err)
		return
	}

	switch hook.Event {
	case "integration_installation":
		// Triggered when an integration has been installed or uninstalled.
		h.infolog.Printf("Initialization web hook (ID %s)", hook.Id)
		// to do nothing
	case "integration_installation_repositories":
		// Triggered when a repository is added or removed from an installation.
		h.infolog.Printf("Initialization web hook for user repositories (ID %s)", hook.Id)
		err = h.initialUserManagement(hook)
	case "push":
		// Any Git push to a Repository, including editing tags or branches.
		// Commits via API actions that update references are also counted. This is the default event.
		h.infolog.Printf("Push hook (ID %s)", hook.Id)
		err = h.runCiCdProcess(hook)
	default:
		h.infolog.Printf("Not processed hook (ID %s), event = %s", hook.Id, hook.Event)
	}

	h.infolog.Printf("Finish process hook (ID %s)", hook.Id)

	if err != nil {
		h.errlog.Printf("Cannot process hook (ID %s): %s", hook.Id, err)
		return
	}
}

func (h *handler) initialUserManagement(hook *githubhook.Hook) error {
	evt := github.IntegrationInstallationRepositoriesEvent{}

	err := hook.Extract(&evt)
	if err != nil {
		return err
	}

	userManagerURL := h.env["USERMAN_SERVICE_HOST"]

	client, err := userManClient.NewClient(nil, userManagerURL)
	if err != nil {
		return err
	}

	h.infolog.Print("Try to create user: ", *evt.Sender.Login)
	user := userManClient.NewUser(*evt.Sender.Login)
	err = client.User.Sync(*user)
	if err != nil {
		return err
	}

	for _, rep := range evt.RepositoriesAdded {
		arr := strings.Split(*rep.FullName, "/")
		repository := userManClient.NewRepository(arr[0], arr[1])

		h.infolog.Print("Try to create repository: ", *rep.FullName)

		// TODO: handle error
		client.Repository.Create(*repository)
	}

	for _, rep := range evt.RepositoriesRemoved {
		arr := strings.Split(*rep.FullName, "/")
		repository := userManClient.NewRepository(arr[0], arr[1])

		h.infolog.Print("Try to remove repository: ", *rep.FullName)

		// TODO: handle error
		client.Repository.Delete(*repository)
	}

	return nil
}

func (h *handler) runCiCdProcess(hook *githubhook.Hook) error {
	evt := github.PushEvent{}

	err := hook.Extract(&evt)
	if err != nil {
		return err
	}

	jenkinsURL := h.env["JENKINS_SERVICE_HOST"]
	token := h.env["JENKINS_TOKEN"]

	client, err := jenkins.NewClient(nil, jenkinsURL)
	if err != nil {
		return err
	}

	// TODO: replace "/" --> "_"  , for example, "username_repName"
	jobName := *evt.Repo.FullName

	h.infolog.Print("Try to proxy push hook for job: ", jobName)

	return client.RunBuild(jobName, token)
}
