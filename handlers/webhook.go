package handlers

import (
	"github.com/google/go-github/github"
	"github.com/takama/router"
	userManClient "github.com/vsaveliev/user-manager/client"
	githubhook "gopkg.in/rjz/githubhook.v0"
)

// TODO: move this data in persistent key-value storage
var installations map[string]int = make(map[string]int, 1)

func (h *Handler) WebHookHandler(c *router.Control) {
	secret := []byte(h.Env["GITHUBINT_TOKEN"])

	hook, err := githubhook.Parse(secret, c.Request)
	if err != nil {
		h.Errlog.Printf("Cannot parse hook (ID %s): %s", hook.Id, err)
		return
	}

	switch hook.Event {
	case "integration_installation":
		// Triggered when an integration has been installed or uninstalled by user.
		h.Infolog.Printf("Initialization web hook (ID %s)", hook.Id)
		err = h.saveInstallation(hook)

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

	if err != nil {
		h.Errlog.Printf("Cannot process hook (ID %s): %s", hook.Id, err)
	}

	h.Infolog.Printf("Finished process hook (ID %s)", hook.Id)
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

	h.Infolog.Print("Try to activate (sync) user with k8s: ", *evt.Sender.Login)

	user := userManClient.NewUser(*evt.Sender.Login)
	err = client.User.Sync(*user)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) runCiCdProcess(hook *githubhook.Hook) error {
	evt := github.PushEvent{}

	err := hook.Extract(&evt)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) saveInstallation(hook *githubhook.Hook) error {
	evt := github.IntegrationInstallationEvent{}

	err := hook.Extract(&evt)
	if err != nil {
		return err
	}

	// save installation in memory for commit status update
	installations[*evt.Sender.Name] = *evt.Installation.ID

	return nil
}

func (h *Handler) getInstallationID(username string) (int, bool) {
	id, ok := installations[username]

	return id, ok
}
