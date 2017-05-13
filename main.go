package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "The full URL to your integration's website.")
}

func authCallbackHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "The full URL to redirect to after a user authorizes an installation.")
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/webhook", webHookHandler)
	http.HandleFunc("/auth_callback", authCallbackHandler)

	http.ListenAndServe(":7788", nil)
}
