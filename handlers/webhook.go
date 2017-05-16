package handlers

import (
	"strings"

	"github.com/google/go-github/github"
	"github.com/takama/router"
	"github.com/vsaveliev/github-integration/jenkins"
	userManClient "github.com/vsaveliev/user-manager/client"
	"gopkg.in/rjz/githubhook.v0"
)

func (h *Handler) WebHookHandler(c *router.Control) {
	secret := []byte(h.Env["GITHUBINT_TOKEN"])
	hook, err := githubhook.Parse(secret, c.Request)
	if err != nil {
		h.Errlog.Printf("Cannot parse hook (ID %s): %s", hook.Id, err)
		return
	}

	switch hook.Event {
	case "integration_installation":
		// Triggered when an integration has been installed or uninstalled.
		h.Infolog.Printf("Initialization web hook (ID %s)", hook.Id)
		// to do nothing
	case "integration_installation_repositories":
		// Triggered when a repository is added or removed from an installation.
		h.Infolog.Printf("Initialization web hook for user repositories (ID %s)", hook.Id)
		err = h.initialUserManagement(hook)
	case "push":
		// Any Git push to a Repository, including editing tags or branches.
		// Commits via API actions that update references are also counted. This is the default event.
		h.Infolog.Printf("Push hook (ID %s)", hook.Id)
		err = h.runCiCdProcess(hook)
	default:
		h.Infolog.Printf("Not processed hook (ID %s), event = %s", hook.Id, hook.Event)
	}

	h.Infolog.Printf("Finish process hook (ID %s)", hook.Id)

	if err != nil {
		h.Errlog.Printf("Cannot process hook (ID %s): %s", hook.Id, err)
		return
	}
}

func (h *Handler) initialUserManagement(hook *githubhook.Hook) error {
	evt := github.IntegrationInstallationRepositoriesEvent{}

	err := hook.Extract(&evt)
	if err != nil {
		return err
	}

	userManagerURL := h.Env["USERMAN_SERVICE_HOST"]

	client, err := userManClient.NewClient(nil, userManagerURL)
	if err != nil {
		return err
	}

	h.Infolog.Print("Try to create user: ", *evt.Sender.Login)
	user := userManClient.NewUser(*evt.Sender.Login)
	err = client.User.Sync(*user)
	if err != nil {
		return err
	}

	for _, rep := range evt.RepositoriesAdded {
		arr := strings.Split(*rep.FullName, "/")
		repository := userManClient.NewRepository(arr[0], arr[1])

		h.Infolog.Print("Try to create repository: ", *rep.FullName)

		// TODO: handle error
		client.Repository.Create(*repository)
	}

	for _, rep := range evt.RepositoriesRemoved {
		arr := strings.Split(*rep.FullName, "/")
		repository := userManClient.NewRepository(arr[0], arr[1])

		h.Infolog.Print("Try to remove repository: ", *rep.FullName)

		// TODO: handle error
		client.Repository.Delete(*repository)
	}

	return nil
}

func (h *Handler) runCiCdProcess(hook *githubhook.Hook) error {
	evt := github.PushEvent{}

	err := hook.Extract(&evt)
	if err != nil {
		return err
	}

	jenkinsURL := h.Env["JENKINS_SERVICE_HOST"]
	token := h.Env["JENKINS_TOKEN"]

	client, err := jenkins.NewClient(nil, jenkinsURL)
	if err != nil {
		return err
	}

	// TODO: replace "/" --> "_"  , for example, "username_repName"
	jobName := *evt.Repo.FullName

	h.Infolog.Print("Try to proxy push hook for job: ", jobName)

	return client.RunBuild(jobName, token)
}
