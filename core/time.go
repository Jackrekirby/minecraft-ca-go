package core

import (
	_ "unsafe" // This is required for go:linkname to work
)

//go:linkname nanotime runtime.nanotime1
func nanotime() int64

func NowInSeconds() float64 {
	return float64(nanotime()) / float64(1e9)
}
