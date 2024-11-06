//go:build js && wasm
// +build js,wasm

package core

import (
	"embed"
	"fmt"
	"image"
	"io/fs"
	"syscall/js"
)

func drawFrame(frameBuffer *image.RGBA, width, height int, ctx js.Value) {
	// Convert Go's RGBA buffer to JS Uint8ClampedArray
	jsData := js.Global().Get("Uint8ClampedArray").New(len(frameBuffer.Pix))
	js.CopyBytesToJS(jsData, frameBuffer.Pix)

	// Create ImageData and put it on the canvas
	imgData := js.Global().Get("ImageData").New(jsData, width, height)
	ctx.Call("putImageData", imgData, 0, 0)
}

func MainWebAssembly() {
	// This is required to prevent the Go runtime from exiting.
	c := make(chan struct{}, 0)

	// Call the function to load and display the image.
	js.Global().Set("loadAndDisplayImage", js.FuncOf(loadAndDisplayImage))

	// Wait indefinitely.
	<-c
}

func loadAndDisplayImage(this js.Value, args []js.Value) interface{} {
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "canvas")
	ctx := canvas.Call("getContext", "2d")

	width := 512
	height := 512

	canvas.Set("width", width)
	canvas.Set("height", height)

	sceneImage := image.NewRGBA(image.Rect(0, 0, width, height))
	go RunEngine(sceneImage)

	// Create a new Image object in JavaScript.
	// img := document.Call("createElement", "img")
	// img.Set("src", "assets/redstone_block.png") // Replace with your actual image URL

	// Function to render the image continuously.
	var renderLoop js.Func
	renderLoop = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Request the next frame.
		// fmt.Println("renderLoop")
		drawFrame(sceneImage, width, height, ctx)
		js.Global().Call("requestAnimationFrame", renderLoop)
		return nil
	})

	// // On image load, draw it onto the canvas.
	// img.Set("onload", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 	width := img.Get("width").Int()
	// 	height := img.Get("height").Int()

	// 	// Set canvas dimensions to match the image.
	// 	canvas.Set("width", width)
	// 	canvas.Set("height", height)

	// 	// Draw the image onto the canvas.
	// 	ctx.Call("drawImage", img, 0, 0, width, height)

	// 	// Start the animation loop.
	// 	js.Global().Call("requestAnimationFrame", renderLoop)
	// 	return nil
	// }))

	js.Global().Call("requestAnimationFrame", renderLoop)

	return nil
}

func OutputSceneImage(img *image.RGBA) {
	// do nothing on wasm
}

func KeyboardEvents(scene *Scene) {
	// do nothing on wasm
	fmt.Println("WASM Keyboard Events")
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
	fmt.Println("WriteGameSame not implemented")
}
