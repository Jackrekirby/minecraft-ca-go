//go:build !js && !wasm
// +build !js,!wasm

package main

import "project_two/core"

func main() {
	core.RunEngineWrapper()
}
