package core

import (
	"fmt"
	"image"
	"time"

	"golang.org/x/image/font"
)

type Player struct {
	Position Point3D
	Rotation Point3D
}

type Scene struct {
	Iteration int
	GameState GameState
	Camera    Camera
	World     World
	Player    Player
	// metrics
	FramesPerSecond                   int // not being used anymore to set frame rate along with other vars
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

func RunGameSave(scene *Scene) {
	period := ratePerSecondToDuration(1)
	// scene.Iteration = 0
	for scene.GameState != Quit {
		gameSave := GameSave{CameraPosition: scene.Camera.Position, CameraRotation: scene.Camera.Rotation}
		WriteGameSame(gameSave)
		time.Sleep(period)
	}
}

type SceneEvent interface {
	Initialise(scene *Scene)
	Update()
	Destroy()
}

func InitialiseScene(scene *Scene, sceneImage *image.RGBA, scale int) {
	gameSave, err := LoadGameSave()
	if err != nil {
		gameSave = GameSave{
			CameraPosition: Point3D{X: 3.5, Y: 5.5, Z: -4},
			CameraRotation: Point3D{X: DegToRad(0), Y: DegToRad(0), Z: DegToRad(0)},
		}
	}
	// gameSave = GameSave{
	// 	CameraPosition: Point3D{0.1, 9.7, -2.7},
	// 	CameraRotation: Point3D{DegToRad(332.5), DegToRad(349.5), 0},
	// }

	scene.GameState = Playing
	scene.FramesPerSecond = 2
	scene.StepsPerSecond = 2
	scene.SubStepsPerSecond = 0

	width := sceneImage.Bounds().Dx()
	height := sceneImage.Bounds().Dy()

	scene.World = World{}
	scene.Camera = Camera{
		Position:    gameSave.CameraPosition,
		Rotation:    gameSave.CameraRotation,
		FOV:         90.0,
		AspectRatio: float64(height) / float64(width),
		Near:        0.1,
		Far:         100.0,
	}
	scene.Player = Player{
		Position: gameSave.CameraPosition,
		Rotation: gameSave.CameraRotation,
	}

	// load assets
	fontFace, err := LoadTrueTypeFont("CascadiaMono.ttf", min(18.0, (18.0/512.0)*float64(width)))
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
}
