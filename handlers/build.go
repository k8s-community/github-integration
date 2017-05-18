package handlers

import (
	"net/http"

	"strconv"

	"encoding/json"
	"io/ioutil"

	"github.com/k8s-community/github-integration/github"
	"github.com/takama/router"
)

// BuildCallbackHandler is handler for callback from build service (system)
func (h *Handler) BuildCallbackHandler(c *router.Control) {
	var build github.BuildCallback

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		h.Errlog.Printf("couldn't read request body: %s", err)
		c.Code(http.StatusBadRequest).Body(nil)
		return
	}

	err = json.Unmarshal(body, &build)
	if err != nil {
		h.Errlog.Printf("couldn't validate request body: %s", err)
		c.Code(http.StatusBadRequest).Body(nil)
		return
	}

	installationID, ok := h.getInstallationID(build.Username)
	if !ok {
		h.Errlog.Printf("cannot find installation for user %s in memory", build.Username)
		c.Code(http.StatusNotFound).Body(nil)
		return
	}

	privKey := []byte(h.Env["GITHUBINT_PRIV_KEY"])
	integrationID, _ := strconv.Atoi(h.Env["GITHUBINT_INTEGRATION_ID"])

	client, err := github.NewClient(nil, integrationID, installationID, privKey)

	err = client.UpdateCommitStatus(build)
	if err != nil {
		h.Errlog.Printf("cannot update commit status, build: %+v", build)

		c.Code(http.StatusInternalServerError).Body(nil)
		return
	}

	c.Code(http.StatusOK).Body(nil)
}
