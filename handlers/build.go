package handlers

import (
	"net/http"

	"strconv"

	"github.com/takama/router"
	"github.com/vsaveliev/github-integration/github"
)

func (h *Handler) BuildCallbackHandler(c *router.Control) {
	build := github.BuildCallback{
		Username:    c.Get(":username"),
		Repository:  c.Get(":repository"),
		CommitHash:  c.Get(":commitHash"),
		State:       c.Get(":state"),
		BuildURL:    c.Get(":buildURL"),
		Description: c.Get(":description"),
		Context:     c.Get(":context"),
	}

	installationID, ok := h.getInstallationID(build.Username)
	if !ok {
		h.Errlog.Printf("cannot find installation for %s in memory", build.Username)
		c.Code(http.StatusNotFound)
		return
	}

	privKey := []byte(h.Env["GITHUBINT_PRIV_KEY"])
	integrationID, _ := strconv.Atoi(h.Env["GITHUBINT_INTEGRATION_ID"])

	client, err := github.NewClient(nil, integrationID, installationID, privKey)

	err = client.UpdateCommitStatus(build)
	if err != nil {
		h.Errlog.Printf("cannot update commit status, build: %+v")

		c.Code(http.StatusInternalServerError)
		return
	}

	c.Code(http.StatusOK)
}
