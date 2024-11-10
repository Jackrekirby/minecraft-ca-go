package core

import (
	"testing"
	"time"
)

func BenchmarkTimeNow(b *testing.B) {
	var totalDuration time.Duration = 0
	for i := 0; i < b.N; i++ {
		startTime := time.Now()
		elapsed := time.Since(startTime)
		totalDuration += elapsed
	}
}

func BenchmarkNanotime(b *testing.B) {
	var totalDuration int64 = 0
	for i := 0; i < b.N; i++ {
		startTime := nanotime()
		elapsed := nanotime() - startTime
		totalDuration += elapsed
	}
}
