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

func ratePerSecondToDuration(rate int) time.Duration {
	return time.Duration(1.0/float64(rate)*1000.0) * time.Millisecond
}

func runRenderLoop(scene *Scene) {
	period := ratePerSecondToDuration(scene.FramesPerSecond)
	for scene.GameState != Quit {
		DrawScene(scene)
		time.Sleep(period)
	}
}

func runGameLoop(scene *Scene) {
	period := ratePerSecondToDuration(scene.StepsPerSecond)
	scene.Iteration = 0
	for scene.GameState != Quit {
		if scene.GameState == Playing || scene.GameState == Pausing {
			numUpdates := scene.World.UpdateWorld()
			if ProcessUserInputs(scene.Iteration, &scene.World) {
				numUpdates += 1
			}
			scene.NumBlockUpdatesInStep = numUpdates
			scene.Iteration = scene.Iteration + 1
			if scene.GameState == Pausing {
				scene.GameState = Paused
			}
		}
		time.Sleep(period)
	}
}

func runGameSave(scene *Scene) {
	period := ratePerSecondToDuration(5)
	scene.Iteration = 0
	for scene.GameState != Quit {
		gameSave := GameSave{CameraPosition: scene.Camera.Position, CameraRotation: scene.Camera.Rotation}
		WriteGameSame(gameSave)
		time.Sleep(period)
	}
}

func RunEngine() {
	fmt.Println("Minecraft 3D Celluar Automata in Go")

	gameSave, err := LoadGameSave()
	if err != nil {
		gameSave = GameSave{
			CameraPosition: Point3D{X: 3.5, Y: 5.5, Z: -4},
			CameraRotation: Point3D{X: DegToRad(0), Y: DegToRad(0), Z: DegToRad(0)},
		}
	}
	scene := Scene{}

	scene.GameState = Playing
	scene.FramesPerSecond = 2
	scene.StepsPerSecond = 2

	scene.World = World{}
	scene.Camera = Camera{
		Position:    gameSave.CameraPosition,
		Rotation:    gameSave.CameraRotation,
		FOV:         90.0,
		AspectRatio: 1.0,
		Near:        0.1,
		Far:         100.0,
	}

	// load assets
	fontFace, err := LoadTrueTypeFont("assets/CascadiaMono.ttf", 18)
	if err != nil {
		panic(fmt.Sprintf("failed to load font: %v", err))
	}
	scene.FontFace = fontFace

	createWorld(&scene.World)

	go KeyboardEvents(&scene)
	go runRenderLoop(&scene)
	go runGameSave(&scene)
	runGameLoop(&scene)
}
