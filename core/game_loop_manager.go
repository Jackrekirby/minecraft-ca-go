package core

import (
	"time"
)

type GameEvent struct {
	Name             string
	Period           int64 // nanoseconds
	Callback         func(event *GameEvent, gameLoopManager *GameLoopManager)
	LastCallTime     int64
	Delta            int64
	Overshoot        int64
	CallCount        int
	ElapsedOvershoot int64
}

type GameLoopManager struct {
	// statistics
	Iterations               int
	SleepCount               int
	NoSleepCount             int
	ElapsedSleepDuration     int64
	ElapsedTimeTillNextEvent int64
	// inputs
	Events          [3]GameEvent
	SleepUndershoot int64
	Quit            chan struct{}
}

func CreateGameEvent(name string, period time.Duration, callback func(event *GameEvent, gameLoopManager *GameLoopManager)) GameEvent {
	return GameEvent{
		Name:     name,
		Period:   period.Nanoseconds(),
		Callback: callback,
	}
}

func (g *GameLoopManager) Initialise(
	events [3]GameEvent,
	sleepUndershoot time.Duration,
	quit chan struct{},
) {
	g.Events = events
	g.SleepUndershoot = sleepUndershoot.Nanoseconds()
	g.Quit = quit
}

func (g *GameLoopManager) Run() {
	// we want all events to update immediately on first call
	now := nanotime()
	for i := range g.Events {
		g.Events[i].LastCallTime = now - g.Events[i].Period
	}

	for {
		select {
		case <-g.Quit:
			return
		default:
			g.Iterate()
		}
	}
}

func (g *GameLoopManager) Iterate() {
	now := nanotime()

	var maxOvershoot int64 = -1
	var maxOvershootIndex int = -1

	for i := range g.Events {
		delta := now - g.Events[i].LastCallTime
		overshoot := delta - g.Events[i].Period
		g.Events[i].Delta = delta
		g.Events[i].Overshoot = overshoot
		if overshoot > maxOvershoot {
			maxOvershoot = overshoot
			maxOvershootIndex = i
		}
		// fmt.Println("a", g.Events[i].Name, overshoot)
	}

	if maxOvershootIndex >= 0 { // call event
		// fmt.Println("c")
		// To avoid overshoot accumulating per update, correct for overshoot next update
		// Clamp overshoot that is longer than one update interval
		event := &g.Events[maxOvershootIndex]
		// fmt.Println(event.Name, event.Delta, event.Overshoot)
		clampedOvershoot := min(event.Overshoot, event.Period)
		event.LastCallTime = now - clampedOvershoot
		event.CallCount++
		event.ElapsedOvershoot += clampedOvershoot
		event.Callback(event, g)

	} else { // sleep
		var minTimeTillNextEvent int64 = 1e12 // overshoot will be negative
		for i := range g.Events {
			overshoot := -g.Events[i].Overshoot
			// fmt.Println("b", g.Events[i].Name, overshoot, minTimeTillNextEvent)
			if minTimeTillNextEvent > overshoot {
				minTimeTillNextEvent = overshoot
			}
		}
		sleepDuration := minTimeTillNextEvent - g.SleepUndershoot
		g.ElapsedTimeTillNextEvent += minTimeTillNextEvent
		if sleepDuration > 0 {
			g.SleepCount++
			g.ElapsedSleepDuration += sleepDuration
			time.Sleep(time.Duration(sleepDuration))
		} else {
			g.NoSleepCount++
		}
	}

	g.Iterations++
}
