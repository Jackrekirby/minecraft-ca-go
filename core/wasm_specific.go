//go:build js && wasm
// +build js,wasm

package core

import (
	"embed"
	"fmt"
	"image"
	"io/fs"
	"syscall/js"
	"time"

	"golang.org/x/image/draw"
)

func drawFrame(frameBuffer *image.RGBA, width, height int, ctx js.Value) {
	// Convert Go's RGBA buffer to JS Uint8ClampedArray
	jsData := js.Global().Get("Uint8ClampedArray").New(len(frameBuffer.Pix))
	js.CopyBytesToJS(jsData, frameBuffer.Pix)

	// Create ImageData object
	imgData := js.Global().Get("ImageData").New(jsData, width, height)

	// Scale and draw the image data on the canvas
	ctx.Call("putImageData", imgData, 0, 0)
}

func MainWebAssembly() {
	// This is required to prevent the Go runtime from exiting.
	c := make(chan struct{}, 0)

	js.Global().Set("runProgram", js.FuncOf(runProgram2))

	// Wait indefinitely.
	<-c
}

func runProgram(this js.Value, args []js.Value) interface{} {
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "canvas")
	ctx := canvas.Call("getContext", "2d")

	scale := 2

	// Get the current screen size
	screenWidth := js.Global().Get("window").Get("innerWidth").Int() / scale * scale
	screenHeight := js.Global().Get("window").Get("innerHeight").Int() / scale * scale

	// Choose the shorter dimension to maintain aspect ratio
	// screenShortDim := min(screenWidth, screenHeight)

	// Set canvas size to fill the smaller screen dimension
	canvas.Set("width", screenWidth)
	canvas.Set("height", screenHeight)

	// Create a new RGBA image buffer with the default resolution (can be adjusted if needed)
	width := screenWidth / scale
	height := screenHeight / scale

	// Create a new RGBA image buffer
	sceneImage := image.NewRGBA(image.Rect(0, 0, width, height))

	// Run the engine or any other processes to fill the framebuffer
	go RunEngine(sceneImage, scale)

	// Rendering loop to continuously update the canvas
	var renderLoop js.Func
	renderLoop = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Call the drawFrame function, passing the current canvas size
		drawFrame(sceneImage, screenWidth, screenHeight, ctx)
		js.Global().Call("requestAnimationFrame", renderLoop)
		return nil
	})

	// Start the render loop
	js.Global().Call("requestAnimationFrame", renderLoop)
	return nil
}

func runProgram2(_ js.Value, _ []js.Value) interface{} {
	go runProgram2Inner()
	return nil
}

func runProgram2Inner() {
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "canvas")
	ctx := canvas.Call("getContext", "2d")

	scale := 1

	// Get the current screen size
	screenWidth := js.Global().Get("window").Get("innerWidth").Int() / scale * scale
	screenHeight := js.Global().Get("window").Get("innerHeight").Int() / scale * scale

	// Choose the shorter dimension to maintain aspect ratio
	// screenShortDim := min(screenWidth, screenHeight)

	// Set canvas size to fill the smaller screen dimension
	canvas.Set("width", screenWidth)
	canvas.Set("height", screenHeight)

	// Create a new RGBA image buffer with the default resolution (can be adjusted if needed)
	width := screenWidth / scale
	height := screenHeight / scale

	// Create a new RGBA image buffer
	sceneImage := image.NewRGBA(image.Rect(0, 0, width, height))

	fmt.Println("sw, sh, w, h", screenWidth, screenHeight, width, height)

	// Run the engine or any other processes to fill the framebuffer
	// updateInterval := (1000000 / 10) * time.Microsecond
	// renderInterval := (1000000 / 10) * time.Microsecond
	sleepUndershoot := 5 * time.Millisecond

	quit := make(chan struct{})

	scene := Scene{}
	InitialiseScene(&scene, sceneImage, scale)
	go KeyboardEvents(&scene)
	// keyboardManager := KeyboardManager{}
	// keyboardManager.Initialise(&scene)

	update := func(event *GameEvent, gameLoopManager *GameLoopManager) {
		// keyboardManager.Update()
		Update(&scene)
		if scene.GameState == Quit {
			close(quit)
		}
	}

	jsData := js.Global().Get("Uint8ClampedArray").New(screenWidth * screenHeight * 4)
	js.CopyBytesToJS(jsData, scaleImage(*sceneImage, float64(scale), draw.NearestNeighbor).Pix)

	sceneImageData := js.Global().Get("ImageData").New(jsData, screenWidth, screenHeight)

	outputSceneImage := func(img *image.RGBA) {
		// drawFrame(img, screenWidth, screenHeight, ctx)
		jsData := sceneImageData.Get("data")
		js.CopyBytesToJS(jsData, img.Pix)
		ctx.Call("putImageData", sceneImageData, 0, 0)
	}

	depthBuffer := make(DepthBuffer, width*height)

	render := func(event *GameEvent, gameLoopManager *GameLoopManager) {
		time.Sleep(1 * time.Millisecond)
		Render(&scene, sceneImage, scale, &depthBuffer, outputSceneImage)
	}

	updateStatistics := func(event *GameEvent, gameLoopManager *GameLoopManager) {
		renderEvent := &gameLoopManager.Events[1]
		scene.RecordedFramesPerSecond = renderEvent.CallCount
		renderEvent.CallCount = 0

		updateEvent := &gameLoopManager.Events[0]
		scene.RecordedStepsPerSecond = updateEvent.CallCount
		updateEvent.CallCount = 0
	}

	g := GameLoopManager{}

	events := [3]GameEvent{
		CreateGameEvent("Update", (1000/10)*time.Millisecond, update),
		CreateGameEvent("Render", (1_000_000/60)*time.Microsecond, render),
		CreateGameEvent("Statistics", (1000/1)*time.Millisecond, updateStatistics),
	}

	g.Initialise(events, sleepUndershoot, quit)

	g.Run()

	// runStatistics := RunGameLoop(updateInterval, renderInterval, update, render, sleepUndershoot, quit, &scene)

	// fmt.Println(runStatistics.String())
}

func OutputSceneImage(img *image.RGBA) {
	// do nothing on wasm
}

func KeyboardEvents(scene *Scene) {
	// do nothing on wasm
	fmt.Println("WASM Keyboard Events")

	onKeyDownMC := func(this js.Value, p []js.Value) interface{} {
		fmt.Println("Key Down")
		// Get the key code or key name
		key := p[0].Get("key").String()
		fmt.Printf("Key Down: %s\n", key)
		HandleKeyPress(scene, key)
		return nil
	}

	onKeyUpMC := func(this js.Value, p []js.Value) interface{} {
		// Get the key code or key name
		key := p[0].Get("key").String()
		fmt.Printf("Key Up: %s\n", key)

		return nil
	}

	// JavaScript function to capture keyboard events
	js.Global().Set("onKeyDownMC", js.FuncOf(onKeyDownMC))
	js.Global().Set("onKeyUpMC", js.FuncOf(onKeyUpMC))
}

type KeyboardManager struct {
	quitKey string
}

// Initialise the keyboard manager
func (km *KeyboardManager) Initialise(scene *Scene) {
	KeyboardEvents(scene)
}

// Update method to handle keyboard events
func (km *KeyboardManager) Update() {
}

// Clean up resources (optional method)
func (km *KeyboardManager) Destroy() {
}

//go:embed assets
var assets embed.FS

func LoadAsset(filename string) ([]byte, error) {
	bytes, err := assets.ReadFile(fmt.Sprintf("assets/%s", filename))
	if err != nil {
		fmt.Printf("Failed to read asset '%s'\n", filename)
		return nil, err
	}
	return bytes, nil
}

func LoadAssets() ([]fs.DirEntry, error) {
	files, err := assets.ReadDir("assets")
	if err != nil {
		fmt.Printf("Failed to read asset directory")
		return nil, err
	}

	return files, nil
}

// Define a struct that matches the JSON structure
type GameSave struct {
	CameraPosition Point3D `json:"CameraPosition"`
	CameraRotation Point3D `json:"CameraRotation"`
}

func LoadGameSave() (GameSave, error) {
	var gameSave GameSave = GameSave{
		CameraPosition: Point3D{X: 1.1988887394336163, Y: 5.5, Z: -4.806844720508948},
		CameraRotation: Point3D{X: 0, Y: 6.021385919380435, Z: 0},
	}

	return gameSave, nil
}

func WriteGameSame(gameSave GameSave) {
	// fmt.Println("WriteGameSame not implemented")
}
