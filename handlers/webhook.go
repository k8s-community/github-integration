package handlers

import (
	"net/http"

	"fmt"

	"github.com/google/go-github/github"
	"github.com/k8s-community/cicd"
	userManClient "github.com/k8s-community/user-manager/client"
	"github.com/takama/router"
	githubhook "gopkg.in/rjz/githubhook.v0"
)

// installations is used for installations storing
// TODO: move this data in persistent key-value storage
var installations map[string]int = make(map[string]int)

// WebHookHandler is common handler for web hooks (installation, repositories installation, push)
func (h *Handler) WebHookHandler(c *router.Control) {
	secret := []byte(h.Env["GITHUBINT_TOKEN"])

	hook, err := githubhook.Parse(secret, c.Request)
	if err != nil {
		h.Errlog.Printf("cannot parse hook (ID %s): %s", hook.Id, err)
		return
	}

	switch hook.Event {
	case "integration_installation":
		// Triggered when an integration has been installed or uninstalled by user.
		h.Infolog.Printf("initialization web hook (ID %s)", hook.Id)
		err = h.saveInstallation(hook)

	case "integration_installation_repositories":
		// Triggered when a repository is added or removed from an installation.
		h.Infolog.Printf("initialization web hook for user repositories (ID %s)", hook.Id)
		err = h.initialUserManagement(hook)

	case "push":
		// Any Git push to a Repository, including editing tags or branches.
		// Commits via API actions that update references are also counted. This is the default event.
		h.Infolog.Printf("push hook (ID %s)", hook.Id)
		err = h.runCiCdProcess(hook)

	default:
		h.Infolog.Printf("not processed hook (ID %s), event = %s", hook.Id, hook.Event)
		c.Code(http.StatusNotFound).Body(nil)
		return
	}

	if err != nil {
		h.Errlog.Printf("cannot process hook (ID %s): %s", hook.Id, err)
		c.Code(http.StatusInternalServerError).Body(nil)
		return
	}

	h.Infolog.Printf("finished to process hook (ID %s)", hook.Id)
	c.Code(http.StatusOK).Body(nil)
}

// initialUserManagement is used for user activation in k8s system
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

	h.Infolog.Print("Try to activate (sync) user in k8s system: ", *evt.Sender.Login)

	user := userManClient.NewUser(*evt.Sender.Login)

	err = client.User.Sync(*user)
	if err != nil {
		return err
	}

	return nil
}

// runCiCdProcess is used for start CI/CD process for some repository from push hook
func (h *Handler) runCiCdProcess(hook *githubhook.Hook) error {
	evt := github.PushEvent{}

	err := hook.Extract(&evt)
	if err != nil {
		return err
	}

	ciCdURL := h.Env["CICD_SERVICE_HOST"]
	client := cicd.NewClient(ciCdURL)

	req := cicd.BuildRequest{
		Username:   evt.Repo.Owner.Name,
		Repository: evt.Repo.Name,
		CommitHash: evt.HeadCommit.ID,
	}

	_, err = client.Build(req)
	if err != nil {
		return fmt.Errorf("cannot run ci/cd process for hook (ID %s): %s", hook.Id, err)
	}

	return nil
}

// saveInstallation saves installation in memory
func (h *Handler) saveInstallation(hook *githubhook.Hook) error {
	evt := github.IntegrationInstallationEvent{}

	err := hook.Extract(&evt)
	if err != nil {
		return err
	}

	h.Infolog.Printf("save installation for user %s (installation ID = %d)", *evt.Sender.Login, *evt.Installation.ID)

	// save installation in memory for commit status update
	installations[*evt.Sender.Login] = *evt.Installation.ID

	return nil
}

// getInstallationID gets installation from memory
func (h *Handler) getInstallationID(username string) (int, bool) {
	id, ok := installations[username]

	return id, ok
}
