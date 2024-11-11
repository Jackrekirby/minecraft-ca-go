package core

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
)

var isCPUProfiling = false

// StartCPUProfile starts the CPU profiler and returns a defer function to stop it.
func StartCPUProfile() func() {
	isCPUProfiling = true
	filename := "cpu.pprof"
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("could not create CPU profile: %v", err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatalf("could not start CPU profile: %v", err)
	}

	log.Println("CPU profiling started")
	return func() {
		pprof.StopCPUProfile()
		f.Close()
		log.Println("CPU profiling stopped")
	}
}

// StartMemProfile writes a memory profile to the specified file and returns a defer function.
func StartMemProfile() func() {
	filename := "mem.pprof"
	return func() {
		f, err := os.Create(filename)
		if err != nil {
			log.Fatalf("could not create memory profile: %v", err)
		}
		defer f.Close()

		// Force a garbage collection to get up-to-date memory statistics.
		runtime.GC()

		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatalf("could not write memory profile: %v", err)
		}
		log.Println("Memory profile written to", filename)
	}
}

// StartTraceProfile starts the trace profiler and returns a defer function to stop it.
func StartTraceProfile(filename string) func() {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("could not create trace file: %v", err)
	}

	if err := trace.Start(f); err != nil {
		log.Fatalf("could not start trace: %v", err)
	}

	log.Println("Trace profiling started")
	return func() {
		trace.Stop()
		f.Close()
		log.Println("Trace profiling stopped")
	}
}
