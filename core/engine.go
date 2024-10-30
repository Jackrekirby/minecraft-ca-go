package core

import (
	"fmt"
	"time"

	"golang.org/x/image/font"
)

type Scene struct {
	Iteration             int
	GameState             GameState
	Camera                Camera
	World                 World
	FramesPerSecond       int
	StepsPerSecond        int
	NumBlockUpdatesInStep int

	FontFace font.Face // should be a store of multiple fonts or internal handler for loaded assets
}

func runRenderLoop(scene *Scene) {
	var period time.Duration = time.Duration(1.0/float64(scene.FramesPerSecond)*1000.0) * time.Millisecond
	fmt.Println("Rendering every", period)
	for scene.GameState != Quit {
		DrawScene(scene)
		time.Sleep(period)
	}
}

func runGameLoop(scene *Scene) {
	var period time.Duration = time.Duration(1.0/float64(scene.StepsPerSecond)*1000.0) * time.Millisecond
	fmt.Println("Iterating every", period)

	scene.Iteration = 0
	for scene.GameState != Quit {
		if scene.GameState == Play {
			numUpdates := scene.World.UpdateWorld()
			if ProcessUserInputs(scene.Iteration, &scene.World) {
				numUpdates += 1
			}
			// fmt.Print("\033[F\033[K")
			// fmt.Println(scene.Iteration, hasAnyBlockUpdated)
			scene.NumBlockUpdatesInStep = numUpdates
			scene.Iteration = scene.Iteration + 1
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func RunEngine() {
	scene := Scene{}

	scene.GameState = Play
	scene.FramesPerSecond = 2
	scene.StepsPerSecond = 2

	scene.World = World{}
	scene.Camera = Camera{
		Position:    Point3D{X: 3.5, Y: 5.5, Z: -4},
		Rotation:    Point3D{X: DegToRad(0), Y: DegToRad(0), Z: DegToRad(0)},
		FOV:         90.0,
		AspectRatio: 1.0,
		Near:        0.1,
		Far:         100.0,
	}

	// load assets
	fontFace, err := LoadTrueTypeFont("assets/CascadiaMono.ttf", 24)
	if err != nil {
		panic(fmt.Sprintf("failed to load font: %v", err))
	}
	scene.FontFace = fontFace

	createWorld(&scene.World)

	go KeyboardEvents(&scene)
	go runRenderLoop(&scene)
	runGameLoop(&scene)
}
