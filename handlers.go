package main

import (
	"fmt"
	"log"
	"net/http"
)

// Handler defines
type handler struct {
	infolog *log.Logger
	errlog  *log.Logger
	env     map[string]string
}

func (h *handler) homeHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "The full URL to your integration's website.")
}

func (h *handler) authCallbackHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "The full URL to redirect to after a user authorizes an installation.")
}
