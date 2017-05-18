package handlers

import (
	"fmt"
	"log"

	"net/http"

	"github.com/k8s-community/github-integration/version"
	"github.com/takama/router"
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

func (h *Handler) HealthzHandler(c *router.Control) {
	c.Code(http.StatusOK).Body("Ok")
}

func (h *Handler) InfoHandler(c *router.Control) {
	c.Code(http.StatusOK).Body(
		map[string]string{
			"version": version.RELEASE,
			"commit":  version.COMMIT,
			"repo":    version.REPO,
		},
	)
}
