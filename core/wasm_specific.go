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

type State struct {
	scale       int
	scene       Scene
	sceneImage  *image.RGBA
	scaledImage *image.RGBA
	depthBuffer DepthBuffer
	quit        chan struct{}
}

func MainWebAssembly() {
	jsProgramWrapper := func(_ js.Value, _ []js.Value) interface{} {
		go runProgram()
		return nil
	}

	js.Global().Set("runProgram", js.FuncOf(jsProgramWrapper))

	// This is required to prevent the Go runtime from exiting.
	select {}
}

func drawFrame(frameBuffer *image.RGBA, width, height int, ctx js.Value) {
	// Convert Go's RGBA buffer to JS Uint8ClampedArray
	jsData := js.Global().Get("Uint8ClampedArray").New(len(frameBuffer.Pix))
	js.CopyBytesToJS(jsData, frameBuffer.Pix)

	// Create ImageData object
	imgData := js.Global().Get("ImageData").New(jsData, width, height)

	// Scale and draw the image data on the canvas
	ctx.Call("putImageData", imgData, 0, 0)
}

func sizeCanvas(state *State) (width int, height int) {
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "canvas")

	screenWidth := js.Global().Get("window").Get("innerWidth").Int()
	screenHeight := js.Global().Get("window").Get("innerHeight").Int()

	canvasWidth := screenWidth / state.scale * state.scale
	canvasHeight := screenHeight / state.scale * state.scale

	canvas.Set("width", canvasWidth)
	canvas.Set("height", canvasHeight)

	width = canvasWidth / state.scale
	height = canvasHeight / state.scale

	state.sceneImage = image.NewRGBA(image.Rect(0, 0, width, height))
	state.scaledImage = image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))
	state.depthBuffer = make(DepthBuffer, width*height)

	jsData := js.Global().Get("Uint8ClampedArray").New(canvasWidth * canvasHeight * 4)
	js.CopyBytesToJS(jsData, scaleImage(state.sceneImage, float64(state.scale), draw.NearestNeighbor).Pix)

	js.Global().Set("sceneImageData", js.Global().Get("ImageData").New(jsData, canvasWidth, canvasHeight))

	fmt.Println("Size [canvas, image]", canvasWidth, canvasHeight, width, height)
	state.scene.Camera.AspectRatio = float64(height) / float64(width)
	return
}

func createGameUpdate(state *State) func(event *GameEvent, gameLoopManager *GameLoopManager) {
	update := func(event *GameEvent, gameLoopManager *GameLoopManager) {
		// keyboardManager.Update()
		Update(&state.scene)
		if state.scene.GameState == Quit {
			close(state.quit)
		}
	}
	return update
}

func createGameRender(state *State) func(event *GameEvent, gameLoopManager *GameLoopManager) {
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "canvas")
	ctx := canvas.Call("getContext", "2d")

	outputSceneImage := func(img *image.RGBA) {
		sceneImageData := js.Global().Get("sceneImageData")
		jsData := sceneImageData.Get("data")
		if !jsData.Truthy() {
			fmt.Println("sceneImageData", sceneImageData)
			fmt.Println("jsData", jsData)
			return
		}
		js.CopyBytesToJS(jsData, img.Pix)
		ctx.Call("putImageData", sceneImageData, 0, 0)

	}

	render := func(event *GameEvent, gameLoopManager *GameLoopManager) {
		time.Sleep(1 * time.Millisecond)
		Render(&state.scene, state.sceneImage, state.scaledImage, state.scale, &state.depthBuffer, outputSceneImage)
	}
	return render
}

func createGameUpdateStatistics(state *State) func(event *GameEvent, gameLoopManager *GameLoopManager) {
	updateStatistics := func(event *GameEvent, gameLoopManager *GameLoopManager) {
		renderEvent := &gameLoopManager.Events[1]
		state.scene.RecordedFramesPerSecond = renderEvent.CallCount
		renderEvent.CallCount = 0

		updateEvent := &gameLoopManager.Events[0]
		state.scene.RecordedStepsPerSecond = updateEvent.CallCount
		updateEvent.CallCount = 0
	}
	return updateStatistics
}

func createOnResizeListener(state *State) {
	var debounceTimer *time.Timer
	onResize := func(this js.Value, args []js.Value) interface{} {
		// If a previous timer exists, stop it
		if debounceTimer != nil {
			debounceTimer.Stop()
		}

		// Create a new timer that waits for debounceDuration
		callback := func() {
			sizeCanvas(state)
		}

		debounceTimer = time.AfterFunc(1*time.Second, callback)
		return nil
	}

	// Set up the resize event listener
	resizeCallback := js.FuncOf(onResize)
	js.Global().Call("addEventListener", "resize", resizeCallback)
}

func runProgram() {
	state := State{}
	state.scale = 1
	state.scene = Scene{}
	state.quit = make(chan struct{})

	sizeCanvas(&state)
	InitialiseScene(&state.scene, state.sceneImage, state.scale)
	cleanupMouseListener := AddMouseListener(&state.scene)
	defer cleanupMouseListener()
	SetupMouseClickEvents(&state.scene)
	go KeyboardEvents(&state.scene)
	go RunGameSave(&state.scene) // consider moving inside update

	createOnResizeListener(&state)

	g := GameLoopManager{}

	events := [3]GameEvent{
		CreateGameEvent("Update", (1000/10)*time.Millisecond, createGameUpdate(&state)),
		CreateGameEvent("Render", (1_000_000/60)*time.Microsecond, createGameRender(&state)),
		CreateGameEvent("Statistics", (1000/1)*time.Millisecond, createGameUpdateStatistics(&state)),
	}

	sleepUndershoot := 5 * time.Millisecond
	g.Initialise(events, sleepUndershoot, state.quit)

	g.Run()
}

func OutputSceneImage(img *image.RGBA) {
	// do nothing on wasm
}

func HandleMouseEvents(scene *Scene, x, y float64) {
	camera := &scene.Camera

	// Update Yaw (Y axis rotation) based on horizontal mouse movement
	camera.Rotation.Y -= x

	// Calculate pitch (X axis rotation) based on vertical mouse movement
	dPitch := y

	// Update only the pitch (X axis rotation), without affecting roll (Z axis rotation)
	camera.Rotation.X -= dPitch

	// Clamp camera X rotation to avoid flipping (e.g., -90 to +90 degrees)
	if camera.Rotation.X > DegToRad(90) {
		camera.Rotation.X = DegToRad(90)
	}
	if camera.Rotation.X < DegToRad(-90) {
		camera.Rotation.X = DegToRad(-90)
	}
}

func KeyboardEvents(scene *Scene) {
	onKeyDownMC := func(this js.Value, p []js.Value) interface{} {
		key := p[0].Get("key").String()
		HandleKeyPress(scene, key, 0.3, DegToRad(5))
		return nil
	}

	canvas := js.Global().Get("document").Call("getElementById", "canvas")
	canvas.Call("addEventListener", "keydown", js.FuncOf(onKeyDownMC))
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

func OpenAsset(filename string) (fs.File, error) {
	file, err := assets.Open(fmt.Sprintf("assets/%s", filename))
	if err != nil {
		fmt.Printf("Failed to open asset '%s'\n", filename)
		return nil, err
	}
	return file, nil
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

func JSGetNow() float64 {
	return js.Global().Get("performance").Call("now").Float()
}

func buildOnMouseMoveCallback(scene *Scene) func(this js.Value, args []js.Value) any {
	const maxDelta float64 = 0.1
	const scaleMovement float64 = 0.00005
	var lastTime float64 = JSGetNow()
	callback := func(this js.Value, args []js.Value) any {
		event := args[0]
		currentTime := js.Global().Get("performance").Call("now").Float()
		// Calculate time delta (in milliseconds)
		timeDelta := min(currentTime-lastTime, 16.67) // clamp at 60fps
		lastTime = currentTime
		if js.Global().Get("document").Get("pointerLockElement").Truthy() {
			deltaX := event.Get("movementX").Float() * timeDelta * scaleMovement
			deltaY := event.Get("movementY").Float() * timeDelta * scaleMovement

			if math.Abs(deltaX) > maxDelta || math.Abs(deltaY) > maxDelta {
				return nil
			}

			// fmt.Println(deltaX, deltaY)

			HandleMouseEvents(scene, deltaX, deltaY)

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

type Mouse struct {
	Dx float64 // pixels
	Dy float64 // pixels
}

func AddMouseListener(scene *Scene) func() {
	var showWelcomePage bool
	err := ReadFromLocalStorage("ShowWelcomePage", &showWelcomePage)
	if err != nil {
		showWelcomePage = true
	}

	// Create JavaScript event listeners
	mouseMoveCallback := js.FuncOf(buildOnMouseMoveCallback(scene))
	// pointerLockChangeCallback := js.FuncOf(onPointerLockChange)

	// Add the event listeners to the document
	js.Global().Get("document").Call("addEventListener", "mousemove", mouseMoveCallback)
	// js.Global().Get("document").Call("addEventListener", "pointerlockchange", pointerLockChangeCallback)

	// Add click event listener to the canvas element
	canvas := js.Global().Get("document").Call("getElementById", "canvas")

	mouseOnClick := func(this js.Value, args []js.Value) any {
		canvas.Call("focus")
		js.Global().Get("document").Get("body").Call("requestPointerLock")
		return nil
	}
	mouseClickCallback := js.FuncOf(mouseOnClick)
	mouseIcon := js.Global().Get("document").Call("getElementById", "mouse-icon")
	mouseIcon.Call("addEventListener", "click", mouseClickCallback)

	controlsIcon := js.Global().Get("document").Call("getElementById", "controls-icon")
	controlsContainer := js.Global().Get("document").Call("getElementById", "controls-container")

	if showWelcomePage {
		controlsContainer.Get("classList").Call("remove", "hide")
	} else {
		controlsContainer.Get("classList").Call("add", "hide")
	}

	controlsOnClick := func(this js.Value, args []js.Value) any {
		fmt.Println("controls icon")
		// Toggle the 'hide' class using classList.toggle()
		controlsContainer.Get("classList").Call("toggle", "hide")
		if showWelcomePage {
			WriteToLocalStorage("ShowWelcomePage", false)
		}

		return nil
	}
	controlsClickCallback := js.FuncOf(controlsOnClick)
	controlsIcon.Call("addEventListener", "click", controlsClickCallback)

	// return cleanup function
	return func() {
		mouseMoveCallback.Release()
		// pointerLockChangeCallback.Release()
		controlsClickCallback.Release()
		mouseClickCallback.Release()
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

// MOUSE CLICK

func SetupMouseClickEvents(scene *Scene) {
	var selectedBlock Block = &WoolBlock{Cyan, None}
	// var x DirectionalBlock = &RedstoneTorch{Left, false}
	handleMouseClick := func(this js.Value, args []js.Value) any {
		if !js.Global().Get("document").Get("pointerLockElement").Truthy() {
			return nil
		}
		event := args[0]
		button := event.Get("button").Int()

		switch button {
		case 0:
			// fmt.Println("Left click [place]")
			// place block
			previousPos, selectedPos := GetRayCastPositions(scene)
			if selectedPos != nil {
				block := scene.World.GetBlock(*selectedPos)
				_, isLever := block.(Lever)
				if isLever {
					toggleLever(*selectedPos, &scene.World)
					return nil
				}
			}
			if previousPos != nil && selectedPos != nil {
				delta := selectedPos.Subtract(*previousPos)
				dir := delta.ToDirection().GetOppositeDirection()
				// fmt.Println(dir)
				// fmt.Printf("selectedBlock type: %T\n", selectedBlock)
				directionalBlock, isDirectionalBlock := selectedBlock.(DirectionalBlock)
				var block Block
				if isDirectionalBlock {
					block = directionalBlock.SetDirection(dir)
					// fmt.Println("isDirectionalBlock", block)
				} else {
					block = selectedBlock
					// fmt.Println("!isDirectionalBlock", block)
				}
				scene.World.SetBlock(*previousPos, block)
			}
			return nil
		case 2:
			// fmt.Println("Right click [destroy]")
			// destroy block
			_, selectedPos := GetRayCastPositions(scene)
			if selectedPos != nil {
				scene.World.SetBlock(*selectedPos, Air{})
			}
			return nil
		default:
			// fmt.Println("Other click [select]")
			// pick block
			_, selectedPos := GetRayCastPositions(scene)
			if selectedPos != nil {
				selectedBlock = scene.World.GetBlock(*selectedPos)
				fmt.Println("Selected Block", selectedBlock.Type(), selectedBlock)
			}

		}
		return nil
	}

	preventContextMenu := func(this js.Value, args []js.Value) any {
		if js.Global().Get("document").Get("pointerLockElement").Truthy() {
			args[0].Call("preventDefault")
		}
		return nil
	}

	// Create JavaScript event listeners
	js.Global().Get("document").Call("addEventListener", "mousedown", js.FuncOf(handleMouseClick))
	js.Global().Get("document").Call("addEventListener", "contextmenu", js.FuncOf(preventContextMenu))
}
