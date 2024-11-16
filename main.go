//go:build !js && !wasm
// +build !js,!wasm

package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"project_two/core"
	"time"
)

func profile() {
	// Start the pprof server in a separate goroutine
	go func() {
		fmt.Println("Starting pprof server at http://localhost:6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			fmt.Printf("Error starting pprof server: %v\n", err)
		}
	}()

	// Example workload to keep the application running
	for {
		fmt.Println("Application running...")
		time.Sleep(2 * time.Second)
	}
}

func main() {
	//profile()
	// defer core.StartCPUProfile()()
	// core.RunEngineWrapper()
	core.RunOBJTest()
}
