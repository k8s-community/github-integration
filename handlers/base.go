package handlers

import (
	"fmt"
	"log"

	"net/http"

	"github.com/k8s-community/github-integration/github"
	"github.com/k8s-community/github-integration/version"
	"github.com/takama/router"
	"strconv"
)

// Handler defines
type Handler struct {
	Infolog *log.Logger
	Errlog  *log.Logger
	Env     map[string]string
}

// HomeHandler is default handler for home page
// TODO: redirect to landing page
func (h *Handler) HomeHandler(c *router.Control) {
	fmt.Fprint(c.Writer, "The full URL to your integration's website.")
}

// AuthCallbackHandler is handler for auth callback
func (h *Handler) AuthCallbackHandler(c *router.Control) {
	fmt.Fprint(c.Writer, "The full URL to redirect to after a user authorizes an installation.")
}

// HealthzHandler todo: add description
func (h *Handler) HealthzHandler(c *router.Control) {
	c.Code(http.StatusOK).Body("Ok")
}

// InfoHandler todo: add description
func (h *Handler) InfoHandler(c *router.Control) {
	c.Code(http.StatusOK).Body(
		map[string]string{
			"version": version.RELEASE,
			"commit":  version.COMMIT,
			"repo":    version.REPO,
		},
	)
}

func (h *Handler) updateCommitStatus(c *router.Control, build *github.BuildCallback) error {
	installationID, ok := h.getInstallationID(build.Username)
	if !ok {
		c.Code(http.StatusNotFound).Body(nil)
		return fmt.Errorf("Couldn't find installation for %s", build.Username)
	}

	privKey := []byte(h.Env["GITHUBINT_PRIV_KEY"])
	integrationID, err := strconv.Atoi(h.Env["GITHUBINT_INTEGRATION_ID"])

	client, err := github.NewClient(nil, integrationID, installationID, privKey)
	if err != nil {
		c.Code(http.StatusInternalServerError).Body(nil)
		return fmt.Errorf("Couldn't init client for github: %s", err)
	}

	err = client.UpdateCommitStatus(build)
	if err != nil {
		c.Code(http.StatusInternalServerError).Body(nil)
		return fmt.Errorf("Couldn't update commit status: %s", err)
	}

	return nil
}
