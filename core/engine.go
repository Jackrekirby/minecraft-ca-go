package core

import (
	"fmt"
	"image"
	"time"

	"golang.org/x/image/font"
)

type Scene struct {
	Iteration int
	GameState GameState
	Camera    Camera
	World     World
	// metrics
	FramesPerSecond                   int
	StepsPerSecond                    int
	SubStepsPerSecond                 int
	NumBlockUpdatesInStep             int
	NumBlockSubUpdateIterationsInStep int
	NumBlockSubUpdatesInStep          int
	RecordedFramesPerSecond           int
	RecordedStepsPerSecond            int

	FontFace font.Face // should be a store of multiple fonts or internal handler for loaded assets
	Tilemap  Tilemap
}

func ratePerSecondToDuration(rate int) time.Duration {
	return time.Duration(1.0/float64(rate)*1000.0) * time.Millisecond
}

func runRenderLoop(scene *Scene, img *image.RGBA) {
	period := ratePerSecondToDuration(scene.FramesPerSecond)
	for scene.GameState != Quit {
		startTime := time.Now()
		DrawScene(scene, img)
		OutputSceneImage(img)
		elapsedTime := time.Since(startTime)
		scene.RecordedFramesPerSecond = int(1.0 / elapsedTime.Seconds())
		sleepTime := period - elapsedTime
		if sleepTime < 0 {
			fmt.Println("Render loop cannot meet target rate")
		}
		time.Sleep(sleepTime)
	}
}

func runGameLoop(scene *Scene) {
	period := ratePerSecondToDuration(scene.StepsPerSecond)
	subperiod := ratePerSecondToDuration(scene.SubStepsPerSecond)
	scene.Iteration = 0
	maxSubUpdateIterations := 50
	for scene.GameState != Quit {
		startTime := time.Now()
		if scene.GameState == Playing || scene.GameState == Pausing {
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
				time.Sleep(subperiod)
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
		}
		elapsedTime := time.Since(startTime)
		if elapsedTime.Seconds() < (1.0 / 10000.0) {
			scene.RecordedStepsPerSecond = 10000.0
		} else {
			scene.RecordedStepsPerSecond = int(1.0 / elapsedTime.Seconds())
		}
		sleepTime := period - elapsedTime
		if sleepTime < 0 {
			fmt.Println("Game loop cannot meet target rate")
		}
		time.Sleep(sleepTime)
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

func RunEngine(sceneImage *image.RGBA) {
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
	scene.SubStepsPerSecond = 0

	scene.World = World{}
	scene.Camera = Camera{
		Position:    gameSave.CameraPosition,
		Rotation:    gameSave.CameraRotation,
		FOV:         90.0,
		AspectRatio: float64(sceneImage.Bounds().Dy()) / float64(sceneImage.Bounds().Dx()),
		Near:        0.1,
		Far:         100.0,
	}

	// load assets
	fontFace, err := LoadTrueTypeFont("CascadiaMono.ttf", 18)
	if err != nil {
		panic(fmt.Sprintf("failed to load font: %v", err))
	}
	scene.FontFace = fontFace

	tilemap, err3 := GenerateTilemap("assets", 16)
	if err3 != nil {
		panic(fmt.Sprintf("failed to load img: %v", err3))
	}
	SaveImage(&tilemap.Image, "output/tilemap.png")
	// fmt.Println(tilemap.Metas)
	scene.Tilemap = *tilemap

	// createWorld(&scene.World)

	go KeyboardEvents(&scene)
	go runRenderLoop(&scene, sceneImage)
	go runGameSave(&scene)
	runGameLoop(&scene)
}
