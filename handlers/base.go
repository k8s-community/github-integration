package handlers

import (
	"fmt"
	"log"

	"github.com/takama/router"
)

// Handler defines
type Handler struct {
	Infolog *log.Logger
	Errlog  *log.Logger
	Env     map[string]string
}

func (h *Handler) HomeHandler(c *router.Control) {
	fmt.Fprint(c.Writer, "The full URL to your integration's website.")
}

func (h *Handler) AuthCallbackHandler(c *router.Control) {
	fmt.Fprint(c.Writer, "The full URL to redirect to after a user authorizes an installation.")
}
