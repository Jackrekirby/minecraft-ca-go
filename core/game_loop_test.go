package core

import (
	"fmt"
	"testing"
	"time"
)

// increasing sleep undershoot reduced overshoot but increases no sleep (high cpu)

// RunStatistics: [sleepUndershoot = 0]
//   Iterations: 125
//   RenderCount: 60
//   UpdateCount: 5
//   Avg Sleep Duration: 8.14836ms
//   Avg Sleep Count: 0.92
//   Avg No-Sleep Count: 0.00
//   Avg Update Overshoot: 8.90248ms
//   Avg Render Overshoot: 7.83861ms

// RunStatistics: [sleepUndershoot = 10 Millisecond]
//   Iterations: 16713090
//   RenderCount: 60
//   UpdateCount: 5
//   Avg Sleep Duration: 4.689738ms
//   Avg Sleep Count: 0.92
//   Avg No-Sleep Count: 257122.54
//   Avg Update Overshoot: 563.28µs
//   Avg Render Overshoot: 1.58545ms

func TestRunGameLoop(t *testing.T) {
	updateInterval := (1000 / 5) * time.Millisecond     // 5 FPS
	renderInterval := (1000000 / 60) * time.Microsecond // 60 FPS
	sleepUndershoot := 5 * time.Millisecond

	quit := make(chan struct{})

	sleepDuration := 1000 * time.Millisecond

	// Set up a timer to quit the loop after 1 second
	go func() {
		time.Sleep(sleepDuration)
		close(quit)
	}()

	update := func(delta time.Duration) {

	}

	render := func(delta time.Duration) {

	}

	// Capture the start time
	startTime := time.Now()
	// Run the game loop
	runStatistics := RunGameLoop(updateInterval, renderInterval, update, render, sleepUndershoot, quit, &Scene{})

	// Measure the elapsed time
	elapsed := time.Since(startTime)

	fmt.Println(runStatistics.String())

	// Check if the function exited in a reasonable time frame
	if elapsed < sleepDuration-200*time.Millisecond || elapsed > sleepDuration+300*time.Millisecond {
		t.Errorf("Expected the game loop to run for about 1 second, but it ran for %v", elapsed)
	}
}

func NanoToSeconds(nano int64) float64 {
	return float64(nano) / 1e9
}

func NanoToMilliSeconds(nano int64) float64 {
	return float64(nano) / 1e6
}

func TestRunGameLoopManager(t *testing.T) {
	sleepDuration := 1000 * time.Millisecond
	sleepUndershoot := 10 * time.Millisecond
	quit := make(chan struct{})

	g := GameLoopManager{}

	events := [3]GameEvent{
		CreateGameEvent("Update", (1000/5)*time.Millisecond, func(event *GameEvent, gameLoopManager *GameLoopManager) {}),       // update
		CreateGameEvent("Render", (1_000_000/60)*time.Microsecond, func(event *GameEvent, gameLoopManager *GameLoopManager) {}), // render
		CreateGameEvent("Statistics", (1000/1)*time.Millisecond, func(event *GameEvent, gameLoopManager *GameLoopManager) {}),   // statistics
	}

	g.Initialise(events, sleepUndershoot, quit)

	go func() {
		time.Sleep(sleepDuration)
		close(quit)
	}()

	startTime := time.Now()
	g.Run()
	elapsed := time.Since(startTime)

	fmt.Println("Iterations\t\t", g.Iterations)
	fmt.Println("SleepCount\t\t", g.SleepCount)
	fmt.Println("NoSleepCount\t\t", g.NoSleepCount)
	fmt.Println("AVG SleepDuration\t", time.Duration(g.ElapsedSleepDuration/int64(g.SleepCount)))
	fmt.Println("AVG TimeTillNextEvent\t", time.Duration(g.ElapsedTimeTillNextEvent/int64(g.SleepCount+g.NoSleepCount)))

	for _, event := range g.Events {
		fmt.Println(event.Name)
		fmt.Println("  CallCount\t\t", event.CallCount)
		fmt.Println("  AVG Overshoot\t\t", time.Duration(event.ElapsedOvershoot/int64(event.CallCount)))
	}

	// Check if the function exited in a reasonable time frame
	if elapsed < sleepDuration-200*time.Millisecond || elapsed > sleepDuration+300*time.Millisecond {
		t.Errorf("Expected the game loop to run for about 1 second, but it ran for %v", elapsed)
	}
}
