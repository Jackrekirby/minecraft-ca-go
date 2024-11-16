//go:build js && wasm
// +build js,wasm

package core

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"io/fs"
	"math"
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
	// c := make(chan struct{}, 0)

	js.Global().Set("runProgram", js.FuncOf(runProgram2))

	// Wait indefinitely.
	// <-c
	select {}
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

// var sceneImageData js.Value

// type State struct {
// 	sceneImageData *js.Value
// }

func sizeCanvas(
	sceneImage **image.RGBA,
	depthBuffer **DepthBuffer,
	scale int,
	scene *Scene,
) (width int, height int) {
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "canvas")

	screenWidth := js.Global().Get("window").Get("innerWidth").Int()
	screenHeight := js.Global().Get("window").Get("innerHeight").Int()

	canvasWidth := screenWidth / scale * scale
	canvasHeight := screenHeight / scale * scale

	canvas.Set("width", canvasWidth)
	canvas.Set("height", canvasHeight)

	width = canvasWidth / scale
	height = canvasHeight / scale

	// newSceneImage := image.NewRGBA(image.Rect(0, 0, width, height))
	*sceneImage = image.NewRGBA(image.Rect(0, 0, width, height))
	// fmt.Println("sceneImage1", sceneImage)
	newDepthBuffer := make(DepthBuffer, width*height)
	*depthBuffer = &newDepthBuffer

	jsData := js.Global().Get("Uint8ClampedArray").New(canvasWidth * canvasHeight * 4)
	js.CopyBytesToJS(jsData, scaleImage(*sceneImage, float64(scale), draw.NearestNeighbor).Pix)

	js.Global().Set("sceneImageData", js.Global().Get("ImageData").New(jsData, canvasWidth, canvasHeight))
	// state.sceneImageData = &x
	// x := &js.Global().Get("ImageData").New(jsData, canvasWidth, canvasHeight)

	fmt.Println("Size [canvas, image]", canvasWidth, canvasHeight, width, height)
	// fmt.Println("sceneImage1", sceneImage)
	scene.Camera.AspectRatio = float64(height) / float64(width)
	return
}

func runProgram2Inner() {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("Recovered from panic:", r)
	// 	}
	// }()

	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "canvas")
	ctx := canvas.Call("getContext", "2d")

	scale := 1

	scene := Scene{}

	var sceneImage *image.RGBA
	var depthBuffer *DepthBuffer
	// var state *State = &State{}
	sizeCanvas(&sceneImage, &depthBuffer, scale, &scene)
	// fmt.Println("sceneImageData", js.Global().Get("sceneImageData"))
	// fmt.Println("depthBuffer", depthBuffer != nil)
	// // Get the current screen size
	// screenWidth := js.Global().Get("window").Get("innerWidth").Int() / scale * scale
	// screenHeight := js.Global().Get("window").Get("innerHeight").Int() / scale * scale

	// // Choose the shorter dimension to maintain aspect ratio
	// // screenShortDim := min(screenWidth, screenHeight)

	// // Set canvas size to fill the smaller screen dimension
	// canvas.Set("width", screenWidth)
	// canvas.Set("height", screenHeight)

	// // Create a new RGBA image buffer with the default resolution (can be adjusted if needed)
	// width := screenWidth / scale
	// height := screenHeight / scale

	// // Create a new RGBA image buffer
	// sceneImage := image.NewRGBA(image.Rect(0, 0, width, height))

	// fmt.Println("sw, sh, w, h", screenWidth, screenHeight, width, height)

	// Run the engine or any other processes to fill the framebuffer
	// updateInterval := (1000000 / 10) * time.Microsecond
	// renderInterval := (1000000 / 10) * time.Microsecond
	sleepUndershoot := 5 * time.Millisecond

	quit := make(chan struct{})

	InitialiseScene(&scene, sceneImage, scale)

	mouse := Mouse{}
	cleanupMouseListener := AddMouseListener(&scene, &mouse)
	defer cleanupMouseListener()
	// HandleMouseEvents(&scene, &mouse)
	go KeyboardEvents(&scene)
	go runGameSave(&scene) // consider moving inside update
	// keyboardManager := KeyboardManager{}
	// keyboardManager.Initialise(&scene)

	update := func(event *GameEvent, gameLoopManager *GameLoopManager) {
		// keyboardManager.Update()
		Update(&scene)
		if scene.GameState == Quit {
			close(quit)
		}
	}

	// jsData := js.Global().Get("Uint8ClampedArray").New(screenWidth * screenHeight * 4)
	// js.CopyBytesToJS(jsData, scaleImage(*sceneImage, float64(scale), draw.NearestNeighbor).Pix)

	// sceneImageData := js.Global().Get("ImageData").New(jsData, screenWidth, screenHeight)

	var debounceTimer *time.Timer
	onResize := func(this js.Value, args []js.Value) interface{} {
		// If a previous timer exists, stop it
		if debounceTimer != nil {
			debounceTimer.Stop()
		}

		// Create a new timer that waits for debounceDuration
		callback := func() {
			sizeCanvas(&sceneImage, &depthBuffer, scale, &scene)
		}

		debounceTimer = time.AfterFunc(1*time.Second, callback)
		return nil
	}

	// Set up the resize event listener
	resizeCallback := js.FuncOf(onResize)
	js.Global().Call("addEventListener", "resize", resizeCallback)

	// sceneImageData := js.Global().Get("ImageData")
	outputSceneImage := func(img *image.RGBA) {
		sceneImageData := js.Global().Get("sceneImageData")
		// drawFrame(img, screenWidth, screenHeight, ctx)
		// fmt.Println("sceneImageData", sceneImageData)
		jsData := sceneImageData.Get("data")
		// fmt.Println("jsData", jsData)
		if !jsData.Truthy() {
			fmt.Println("sceneImageData", sceneImageData)
			fmt.Println("jsData", jsData)
			return
		}
		// fmt.Println("Length of sceneImageData:", len(img.Pix))
		// fmt.Println("Length of jsData:", jsData.Length())
		js.CopyBytesToJS(jsData, img.Pix)
		ctx.Call("putImageData", sceneImageData, 0, 0)
	}

	render := func(event *GameEvent, gameLoopManager *GameLoopManager) {
		time.Sleep(1 * time.Millisecond)
		Render(&scene, sceneImage, scale, depthBuffer, outputSceneImage)
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

func HandleMouseEvents(scene *Scene, mouse *Mouse) {

	rotationSpeed := 0.005 // Adjust this value for sensitivity
	camera := &scene.Camera

	// Update Yaw (Y axis rotation) based on horizontal mouse movement
	camera.Rotation.Y -= float64(mouse.Dx) * rotationSpeed

	// Calculate pitch (X axis rotation) based on vertical mouse movement
	dPitch := float64(mouse.Dy) * rotationSpeed

	// Update only the pitch (X axis rotation), without affecting roll (Z axis rotation)
	camera.Rotation.X -= dPitch

	// Ensure roll (Z axis rotation) remains unchanged
	// camera.Rotation.Z = 0.0

	// rotationSpeed := 0.005 // Adjust this value for sensitivity
	// camera := &scene.Camera

	// // Update Yaw (Y axis rotation) based on horizontal mouse movement
	// camera.Rotation.Y -= float64(mouse.Dx) * rotationSpeed

	// // Calculate pitch (X and Z axis rotation) based on vertical mouse movement
	// dPitch := float64(mouse.Dy) * rotationSpeed

	// // Distribute dPitch between X and Z axes based on current Y rotation (yaw)
	// camera.Rotation.X -= dPitch * math.Cos(camera.Rotation.Y)
	// camera.Rotation.Z -= dPitch * math.Sin(camera.Rotation.Y)

	// Clamp camera X rotation to avoid flipping (e.g., -90 to +90 degrees)
	if camera.Rotation.X > DegToRad(90) {
		camera.Rotation.X = DegToRad(90)
	}
	if camera.Rotation.X < DegToRad(-90) {
		camera.Rotation.X = DegToRad(-90)
	}

	// // Optional: Clamp Z rotation (if needed)
	// if camera.Rotation.Z > DegToRad(90) {
	// 	camera.Rotation.Z = DegToRad(90)
	// }
	// if camera.Rotation.Z < DegToRad(-90) {
	// 	camera.Rotation.Z = DegToRad(-90)
	// }

	// Reset mouse deltas
	mouse.Dx = 0
	mouse.Dy = 0
}

func KeyboardEvents(scene *Scene) {
	// do nothing on wasm
	// fmt.Println("WASM Keyboard Events")

	onKeyDownMC := func(this js.Value, p []js.Value) interface{} {
		// fmt.Println("Key Down")
		// Get the key code or key name
		key := p[0].Get("key").String()
		// fmt.Printf("Key Down: %s\n", key)
		HandleKeyPress(scene, key, 0.3, DegToRad(5))
		return nil
	}

	// onKeyUpMC := func(this js.Value, p []js.Value) interface{} {
	// 	// Get the key code or key name
	// 	key := p[0].Get("key").String()
	// 	fmt.Printf("Key Up: %s\n", key)

	// 	return nil
	// }

	// JavaScript function to capture keyboard events
	// js.Global().Set("onKeyDownMC", js.FuncOf(onKeyDownMC))

	canvas := js.Global().Get("document").Call("getElementById", "canvas")
	canvas.Call("addEventListener", "keydown", js.FuncOf(onKeyDownMC))
	// js.Global().Set("onKeyUpMC", js.FuncOf(onKeyUpMC))
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
	var gameSave GameSave
	err := ReadFromLocalStorage("game", &gameSave)
	return gameSave, err
}

func WriteGameSame(gameSave GameSave) {
	WriteToLocalStorage("game", gameSave)
}

// MOUSE CONTROL

func buildOnMouseMoveCallback(scene *Scene, mouse *Mouse) func(this js.Value, args []js.Value) any {
	const maxDelta float64 = 30
	callback := func(this js.Value, args []js.Value) any {
		event := args[0]
		if js.Global().Get("document").Get("pointerLockElement").Truthy() {
			deltaX := event.Get("movementX").Int()
			deltaY := event.Get("movementY").Int()

			if math.Abs(float64(deltaX)) > maxDelta || math.Abs(float64(deltaY)) > maxDelta {
				return nil
			}

			mouse.Dx = deltaX
			mouse.Dy = deltaY

			HandleMouseEvents(scene, mouse)

			// Log the camera movement
			// js.Global().Get("console").Call("log", "Camera move: deltaX =", deltaX, ", deltaY =", deltaY)
		}
		return nil
	}
	return callback
}

func onPointerLockChange(this js.Value, args []js.Value) any {
	if js.Global().Get("document").Get("pointerLockElement").IsNull() {
		js.Global().Get("console").Call("log", "Pointer lock exited")
	}
	return nil
}

func onClick(this js.Value, args []js.Value) any {
	js.Global().Get("document").Get("body").Call("requestPointerLock")
	return nil
}

type Mouse struct {
	Dx int // pixels
	Dy int // pixels
}

func AddMouseListener(scene *Scene, mouse *Mouse) func() {
	// Create JavaScript event listeners
	mouseMoveCallback := js.FuncOf(buildOnMouseMoveCallback(scene, mouse))
	// pointerLockChangeCallback := js.FuncOf(onPointerLockChange)
	clickCallback := js.FuncOf(onClick)

	// Add the event listeners to the document
	js.Global().Get("document").Call("addEventListener", "mousemove", mouseMoveCallback)
	// js.Global().Get("document").Call("addEventListener", "pointerlockchange", pointerLockChangeCallback)

	// Add click event listener to the canvas element
	canvas := js.Global().Get("document").Call("getElementById", "canvas")
	canvas.Call("addEventListener", "click", clickCallback)

	// return cleanup function
	return func() {
		mouseMoveCallback.Release()
		// pointerLockChangeCallback.Release()
		clickCallback.Release()
	}
}

// Session Storage

// WriteLocalStorage writes a byte array to localStorage using a given key.
func WriteBytesToLocalStorage(key string, data []byte) {
	// Convert byte array to a Base64-encoded string (JavaScript's localStorage can only store strings).
	encoded := base64.StdEncoding.EncodeToString(data)

	// Get the localStorage object.
	localStorage := js.Global().Get("localStorage")
	if !localStorage.Truthy() {
		fmt.Println("localStorage is not available")
		return
	}

	// Store the encoded data under the given key.
	localStorage.Call("setItem", key, encoded)
	// fmt.Println("Data written to localStorage:", key)
}

// ReadLocalStorage reads a byte array from localStorage using a given key.
func ReadBytesFromLocalStorage(key string) ([]byte, error) {
	// Get the localStorage object.
	localStorage := js.Global().Get("localStorage")
	if !localStorage.Truthy() {
		return nil, fmt.Errorf("localStorage is not available")
	}

	// Retrieve the Base64-encoded string from localStorage.
	item := localStorage.Call("getItem", key)
	if item.IsNull() {
		return nil, fmt.Errorf("no data found for key: %s", key)
	}

	// Decode the Base64 string to a byte array.
	data, err := base64.StdEncoding.DecodeString(item.String())
	if err != nil {
		return nil, fmt.Errorf("failed to decode data: %w", err)
	}

	return data, nil
}

func WriteToLocalStorage(key string, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("Error writing struct to JSON: %s", err))
	}

	WriteBytesToLocalStorage(key, jsonData)
}

func ReadFromLocalStorage(key string, data any) error {
	rawData, err := ReadBytesFromLocalStorage(key)
	if err != nil {
		return fmt.Errorf("error reading session storage: %w", err)
	}

	err = json.Unmarshal(rawData, data)
	if err != nil {
		return fmt.Errorf("error reading JSON data: %w", err)
	}
	return nil
}
