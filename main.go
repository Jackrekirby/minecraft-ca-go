//go:build !js && !wasm
// +build !js,!wasm

package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"project_two/core"
)

func profile() {
	// Start a pprof server in a separate goroutine
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

func main() {
	go profile()
	core.RunEngineWrapper()

}
