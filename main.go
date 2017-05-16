package main

import (
	"log"
	"net/http"
	"os"

	"github.com/takama/router"
	"github.com/vsaveliev/github-integration/handlers"
)

const (
	apiPrefix = "/api/v1"
)

func main() {
	keys := []string{
		"GITHUBINT_SERVICE_PORT", "GITHUBINT_TOKEN", "GITHUBINT_PRIV_KEY", "GITHUBINT_INTEGRATION_ID",
		"USERMAN_SERVICE_HOST", "USERMAN_SERVICE_PORT",
		"JENKINS_SERVICE_HOST", "JENKINS_SERVICE_PORT", "JENKINS_TOKEN",
	}

	h := &handlers.Handler{
		Infolog: log.New(os.Stdout, "[GITHUBINT:INFO]: ", log.LstdFlags),
		Errlog:  log.New(os.Stderr, "[GITHUBINT:ERROR]: ", log.LstdFlags),
		Env:     make(map[string]string, len(keys)),
	}

	for _, key := range keys {
		value := os.Getenv(key)
		if value == "" {
			h.Errlog.Fatalf("%s environment variable was not set", key)
		}
		h.Env[key] = value
	}

	r := router.New()

	r.GET(apiPrefix+"/", h.HomeHandler)
	r.GET(apiPrefix+"/home", h.HomeHandler)
	r.POST(apiPrefix+"/webhook", h.WebHookHandler)
	r.POST(apiPrefix+"/auth_callback", h.AuthCallbackHandler)
	r.POST(apiPrefix+"/build-cb", h.BuildCallbackHandler)

	http.ListenAndServe(":"+h.Env["GITHUBINT_SERVICE_PORT"], nil)
}
