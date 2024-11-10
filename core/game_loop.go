package core

import (
	"fmt"
	"image"
	"time"

	"golang.org/x/image/draw"
)

type RunStatistics struct {
	Iterations      int
	UpdateCount     int
	RenderCount     int
	SleepDuration   time.Duration
	SleepCount      int
	NoSleepCount    int
	UpdateOvershoot time.Duration
	RenderOvershoot time.Duration
}

func (r *RunStatistics) EventCount() int {
	return r.UpdateCount + r.RenderCount
}

func (r *RunStatistics) AvgSleepDuration() time.Duration {
	c := r.EventCount()
	if c == 0 {
		return 0
	}
	return r.SleepDuration / time.Duration(c)
}

func (r *RunStatistics) AvgUpdateOvershoot() time.Duration {
	if r.UpdateCount == 0 {
		return 0
	}
	return r.UpdateOvershoot / time.Duration(r.UpdateCount)
}

func (r *RunStatistics) AvgRenderOvershoot() time.Duration {
	if r.RenderCount == 0 {
		return 0
	}
	return r.RenderOvershoot / time.Duration(r.RenderCount)
}

func (r *RunStatistics) AvgNoSleepCount() float64 {
	c := r.EventCount()
	if c == 0 {
		return 0
	}
	return float64(r.NoSleepCount) / float64(c)
}

func (r *RunStatistics) AvgSleepCount() float64 {
	c := r.EventCount()
	if c == 0 {
		return 0
	}
	return float64(r.SleepCount) / float64(c)
}

// String method for custom formatting
func (r *RunStatistics) String() string {
	return fmt.Sprintf(
		"RunStatistics:\n"+
			"  Iterations: %d\n"+
			"  RenderCount: %d\n"+
			"  UpdateCount: %d\n"+
			"  Avg Sleep Duration: %s\n"+
			"  Avg Sleep Count: %.2f\n"+
			"  Avg No-Sleep Count: %.2f\n"+
			"  Avg Update Overshoot: %s\n"+
			"  Avg Render Overshoot: %s",
		r.Iterations,
		r.RenderCount,
		r.UpdateCount,
		r.AvgSleepDuration(),
		r.AvgSleepCount(),
		r.AvgNoSleepCount(),
		r.AvgUpdateOvershoot(),
		r.AvgRenderOvershoot(),
	)
}

func RunGameLoop(
	updateInterval time.Duration,
	renderInterval time.Duration,
	update func(time.Duration),
	render func(time.Duration),
	sleepUndershoot time.Duration,
	quit chan struct{},
	scene *Scene,
) RunStatistics {
	// we want update immediately
	updateIntervalNs := updateInterval.Nanoseconds()
	renderIntervalNs := renderInterval.Nanoseconds()
	sleepUndershootNs := sleepUndershoot.Nanoseconds()
	lastUpdate := nanotime() - updateIntervalNs
	lastRender := nanotime() - renderIntervalNs
	runStatistics := RunStatistics{}
	// Game loop
	for {
		select {
		case <-quit:
			return runStatistics
		default:
			now := nanotime()

			updateDelta := now - lastUpdate
			renderDelta := now - lastRender
			updateOvershoot := updateDelta - updateIntervalNs
			renderOvershoot := renderDelta - renderIntervalNs

			// Choose to update or render based on which one is most overshot
			// To avoid overshoot accumulating per update, correct for overshoot next update
			// Clamp overshoot that is longer than one update interval
			if updateOvershoot > renderOvershoot && updateOvershoot >= 0 { // Update
				clampedOvershoot := min(updateOvershoot, updateIntervalNs)
				lastUpdate = now - clampedOvershoot

				runStatistics.UpdateOvershoot += time.Duration(clampedOvershoot)
				runStatistics.UpdateCount++
				update(time.Duration(updateDelta))
			} else if renderOvershoot >= 0 { // Render
				clampedOvershoot := min(renderOvershoot, updateIntervalNs)
				lastRender = now - clampedOvershoot

				runStatistics.RenderOvershoot += time.Duration(clampedOvershoot)
				runStatistics.RenderCount++
				render(time.Duration(renderDelta))
				// we must sleep after render for UI events
				time.Sleep(1 * time.Millisecond)
			} else { // Wait
				timeTillNextUpdate := -updateOvershoot
				timeTillNextRender := -renderOvershoot
				timeTillNextEvent := min(timeTillNextUpdate, timeTillNextRender)

				sleepDuration := time.Duration(timeTillNextEvent - sleepUndershootNs)
				if sleepDuration > 0 {
					// fmt.Println("Sleep [now, duration]\t", now.Format("15:04:05.000"), sleepDuration)
					runStatistics.SleepDuration += sleepDuration
					runStatistics.SleepCount++
					time.Sleep(sleepDuration)
				} else {
					// fmt.Println("NoSleep [now, duration]\t", now.Format("15:04:05.000"), timeTillNextUpdate, sleepDuration)
					runStatistics.NoSleepCount++
				}
			}
			runStatistics.Iterations++

			// // Update logic at a fixed rate
			// if updateDelta >= updateInterval {
			// 	// To avoid overshoot accumulating per update, correct for overshoot next frame
			// 	// Clamp overshoot that is longer than one frame
			// 	overshoot := min(updateDelta-updateInterval, updateInterval)
			// 	lastUpdate = now.Add(-overshoot)
			// 	// fmt.Println("Update [now, overshoot]\t", now.Format("15:04:05.000"), overshoot)
			// 	runStatistics.Overshoot += overshoot
			// 	runStatistics.Iterations++

			// 	update(updateDelta)
			// } else {
			// 	timeTillNextUpdate := min(updateInterval - updateDelta)
			// 	if timeTillNextUpdate > sleepUndershoot {
			// 		sleepDuration := timeTillNextUpdate - sleepUndershoot
			// 		// fmt.Println("Sleep [now, duration]\t", now.Format("15:04:05.000"), sleepDuration)
			// 		runStatistics.SleepDuration += sleepDuration
			// 		runStatistics.SleepCount++
			// 		time.Sleep(sleepDuration)
			// 	} else {
			// 		//fmt.Println("NoSleep [now, duration]\t", now.Format("15:04:05.000"), timeTillNextUpdate)
			// 		runStatistics.NoSleepCount++
			// 	}
			// }
		}
	}
}

func Render(scene *Scene, img *image.RGBA, scale int, depthBuffer *DepthBuffer, outputSceneImage func(*image.RGBA)) {
	// startTime := NowInSeconds()
	DrawScene(scene, img, depthBuffer)

	var scaledImage *image.RGBA
	if scale > 1 {
		scaledImage = scaleImage(*img, float64(scale), draw.NearestNeighbor)
	} else {
		scaledImage = img
	}

	outputSceneImage(scaledImage)
	// elapsedTime := NowInSeconds() - startTime
	// scene.RecordedFramesPerSecond = int(1.0 / elapsedTime)
}

func Update(scene *Scene) {
	if scene.GameState != Playing && scene.GameState != Pausing {
		return
	}
	maxSubUpdateIterations := 50
	// startTime := NowInSeconds()

	numUpdates := 0
	// Process User Inputs
	if ProcessUserInputs(scene.Iteration, &scene.World) {
		numUpdates += 1
	}
	// Process Sub Updates
	totalSubUpdates := 0
	i := 0
	for i < maxSubUpdateIterations {
		numSubUpdates := scene.World.SubUpdateWorld()
		totalSubUpdates += numSubUpdates
		if numSubUpdates == 0 {
			break
		}
		i++
	}
	scene.NumBlockSubUpdateIterationsInStep = i
	scene.NumBlockSubUpdatesInStep = totalSubUpdates

	// Process Updates
	numUpdates += scene.World.UpdateWorld()

	scene.NumBlockUpdatesInStep = numUpdates
	scene.Iteration = scene.Iteration + 1
	if scene.GameState == Pausing {
		scene.GameState = Paused
	}
	// elapsedTime := NowInSeconds() - startTime
	// if elapsedTime < (1.0 / 10_000.0) {
	// 	scene.RecordedStepsPerSecond = 10_000
	// } else {
	// 	scene.RecordedStepsPerSecond = int(1.0 / elapsedTime)
	// }
}

func RunEngine3(sceneImage *image.RGBA, scale int) {
	updateInterval := (1000 / 5) * time.Millisecond     // 5 FPS
	renderInterval := (1000000 / 10) * time.Microsecond // 60 FPS
	sleepUndershoot := 5 * time.Millisecond

	quit := make(chan struct{})

	scene := Scene{}
	InitialiseScene(&scene, sceneImage, scale)
	go KeyboardEvents(&scene)
	// keyboardManager := KeyboardManager{}
	// keyboardManager.Initialise(&scene)

	update := func(delta time.Duration) {
		// keyboardManager.Update()
		Update(&scene)
		if scene.GameState == Quit {
			close(quit)
		}
	}

	width, height := sceneImage.Bounds().Dx(), sceneImage.Bounds().Dy()
	depthBuffer := make(DepthBuffer, width*height)

	render := func(delta time.Duration) {
		Render(&scene, sceneImage, scale, &depthBuffer, OutputSceneImage)
	}

	runStatistics := RunGameLoop(updateInterval, renderInterval, update, render, sleepUndershoot, quit, &scene)

	// keyboardManager.Destroy()

	fmt.Println(runStatistics.String())
}

func RunEngine2(sceneImage *image.RGBA, scale int) {
	sleepUndershoot := 10 * time.Millisecond
	quit := make(chan struct{})

	scene := Scene{}

	update := func(event *GameEvent, gameLoopManager *GameLoopManager) {
		// keyboardManager.Update()
		Update(&scene)
		if scene.GameState == Quit {
			close(quit)
		}
	}

	width, height := sceneImage.Bounds().Dx(), sceneImage.Bounds().Dy()
	depthBuffer := make(DepthBuffer, width*height)

	render := func(event *GameEvent, gameLoopManager *GameLoopManager) {
		Render(&scene, sceneImage, scale, &depthBuffer, OutputSceneImage)
	}

	updateStatistics := func(event *GameEvent, gameLoopManager *GameLoopManager) {
		renderEvent := &gameLoopManager.Events[1]
		scene.RecordedFramesPerSecond = renderEvent.CallCount
		renderEvent.CallCount = 0

		updateEvent := &gameLoopManager.Events[0]
		scene.RecordedStepsPerSecond = updateEvent.CallCount
		updateEvent.CallCount = 0
	}

	events := [3]GameEvent{
		CreateGameEvent("Update", (1000/2)*time.Millisecond, update),
		CreateGameEvent("Render", (1_000_000/2)*time.Microsecond, render),
		CreateGameEvent("Statistics", (1000/1)*time.Millisecond, updateStatistics),
	}

	g := GameLoopManager{}
	g.Initialise(events, sleepUndershoot, quit)

	InitialiseScene(&scene, sceneImage, scale)
	go KeyboardEvents(&scene)

	g.Run()
}
