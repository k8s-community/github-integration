package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"github.com/k8s-community/cicd"
	githubWrap "github.com/k8s-community/github-integration/github"
	userManClient "github.com/k8s-community/user-manager/client"
	"github.com/takama/router"
	githubhook "gopkg.in/rjz/githubhook.v0"
	"github.com/AlekSi/pointer"
	"github.com/k8s-community/github-integration/models"
	"gopkg.in/reform.v1"
)

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
		err = h.runCiCdProcess(c, hook)
		if err != nil {
			h.Infolog.Printf("cannot run ci/cd process for hook (ID %s): %s", hook.Id, err)
			c.Code(http.StatusBadRequest).Body(nil)
			return
		}

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

	userManagerURL := h.Env["USERMAN_BASE_URL"]

	client, err := userManClient.NewClient(nil, userManagerURL)
	if err != nil {
		return err
	}

	h.Infolog.Print("Try to activate (sync) user in k8s system: ", *evt.Sender.Login)

	user := userManClient.NewUser(*evt.Installation.Account.Login)

	code, err := client.User.Sync(user)
	if err != nil {
		return err
	}

	h.Infolog.Printf("Service user-man, method sync, returned code: %d", code)

	return nil
}

// runCiCdProcess is used for start CI/CD process for some repository from push hook
func (h *Handler) runCiCdProcess(c *router.Control, hook *githubhook.Hook) error {
	evt := github.PushEvent{}

	err := hook.Extract(&evt)
	if err != nil {
		return err
	}

	h.setInstallationID(*evt.Repo.Owner.Name, *evt.Installation.ID)

	if !strings.HasPrefix(*evt.Ref, "refs/heads/" + h.Env["GITHUBINT_BRANCH"]) {
		return fmt.Errorf("incorrect branch %s for ci/cd process", *evt.Ref)
	}

	ciCdURL := h.Env["CICD_BASE_URL"]

	client := cicd.NewClient(ciCdURL)

	// set Pending sttatus to Github
	build := &githubWrap.BuildCallback{
		Username:   *evt.Repo.Owner.Name,
		Repository: *evt.Repo.Name,
		CommitHash: *evt.HeadCommit.ID,
		State:      "pending",
		BuildURL:   pointer.ToString("https://k8s.community"), // TODO fix it
		Context: pointer.ToString("k8s-community/cicd"), // move to constant!
		Description: pointer.ToString("Waiting for release..."),
	}
	err = h.updateCommitStatus(c, build)
	if err != nil {
		h.Errlog.Printf("cannot update commit status, build: %+v, err: %s", build, err)
	}

	// run CICD process
	req := &cicd.BuildRequest{
		Username:   *evt.Repo.Owner.Name,
		Repository: *evt.Repo.Name,
		CommitHash: *evt.HeadCommit.ID,
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

	// save installation for commit status update
	err = h.setInstallationID(*evt.Installation.Account.Login, *evt.Installation.ID)
	if err != nil {
		h.Errlog.Printf("Couldn't save installation: %+v", err)
	}

	return nil
}

// installationID gets installation from DB
func (h *Handler) installationID(username string) (*int, error) {
	st, err := h.DB.FindOneFrom(models.InstallationTable, "username", username)
	if err != nil {
		return nil, err
	}
	inst := st.(*models.Installation)

	return pointer.ToInt(inst.InstallationID), nil
}

func (h *Handler) setInstallationID(username string, instID int) error {
	var inst *models.Installation

	st, err := h.DB.FindOneFrom(models.InstallationTable, "username", username)
	if err != nil && err != reform.ErrNoRows {
		return err
	}

	if err == nil {
		inst = st.(*models.Installation)
	}

	inst.InstallationID = instID
	inst.Username = username

	err = h.DB.Save(inst)

	return err
}
