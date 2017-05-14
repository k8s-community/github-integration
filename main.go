package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	keys := []string{
		"GITHUBINT_SERVICE_PORT", "GITHUBINT_TOKEN",
		"USERMAN_SERVICE_HOST", "USERMAN_SERVICE_PORT",
		"JENKINS_SERVICE_HOST", "JENKINS_SERVICE_PORT", "JENKINS_TOKEN",
	}

	h := &handler{
		infolog: log.New(os.Stdout, "[GITHUBINT:INFO]: ", log.LstdFlags),
		errlog:  log.New(os.Stderr, "[GITHUBINT:ERROR]: ", log.LstdFlags),
		env:     make(map[string]string, len(keys)),
	}

	for _, key := range keys {
		value := os.Getenv(key)
		if value == "" {
			h.errlog.Fatalf("%s environment variable was not set", key)
		}
		h.env[key] = value
	}

	http.HandleFunc("/", h.homeHandler)
	http.HandleFunc("/home", h.homeHandler)
	http.HandleFunc("/webhook", h.webHookHandler)
	http.HandleFunc("/auth_callback", h.authCallbackHandler)

	http.ListenAndServe(":"+h.env["GITHUBINT_SERVICE_PORT"], nil)
}
