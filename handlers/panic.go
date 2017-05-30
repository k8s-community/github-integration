package handlers

import (
	"log"
	"runtime/debug"

	"github.com/takama/router"
)

func Panic(c *router.Control) {
	log.Println(debug.Stack())
}
