package main

import (
	"log"
	"os"

	"os/signal"
	"syscall"

	"github.com/k8s-community/github-integration/handlers"
	"github.com/takama/router"
)

const (
	apiPrefix = "/api/v1"
)

// main function
func main() {
	keys := []string{
		"GITHUBINT_SERVICE_PORT",
		"GITHUBINT_TOKEN", "GITHUBINT_PRIV_KEY", "GITHUBINT_INTEGRATION_ID",
		"USERMAN_SERVICE_HOST", "USERMAN_SERVICE_PORT",
		"CICD_SERVICE_HOST", "CICD_SERVICE_PORT",
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

	r.GET("/healthz", h.HealthzHandler)
	r.GET("/info", h.InfoHandler)

	r.GET(apiPrefix+"/home", h.HomeHandler)
	r.POST(apiPrefix+"/webhook", h.WebHookHandler)
	r.POST(apiPrefix+"/auth_callback", h.AuthCallbackHandler)
	r.POST(apiPrefix+"/build-cb", h.BuildCallbackHandler)
	h.Infolog.Printf("start listening port %s", h.Env["GITHUBINT_SERVICE_PORT"])

	go r.Listen(":" + h.Env["GITHUBINT_SERVICE_PORT"])

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	killSignal := <-interrupt
	h.Infolog.Println("Got signal:", killSignal)

	if killSignal == os.Kill {
		h.Infolog.Println("Service was killed")
	} else {
		h.Infolog.Println("Service was terminated by system signal")
	}

	h.Infolog.Println("shutdown")
}
