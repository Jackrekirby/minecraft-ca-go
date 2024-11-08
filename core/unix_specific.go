//go:build !js && !wasm
// +build !js,!wasm

package core

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/fs"
	"os"

	"github.com/eiannone/keyboard"
)

func KeyboardEvents(scene *Scene) {
	// Open the keyboard
	err := keyboard.Open()
	if err != nil {
		fmt.Println("Error opening keyboard:", err)
		return
	}
	defer keyboard.Close()

	fmt.Println("Listening for keyboard inputs. Press 'q' to quit.")

	for {
		// Read key press
		key, _, err := keyboard.GetKey()
		if err != nil {
			fmt.Println("Error reading key:", err)
			break
		}
		HandleKeyPress(scene, string(key))
	}
}

func OutputSceneImage(img *image.RGBA) {
	file, err := os.Create("output/scene.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}

func RunEngineWrapper() {
	imageSize := 512
	sceneImage := image.NewRGBA(image.Rect(0, 0, imageSize, imageSize))
	RunEngine(sceneImage)
}

func LoadAsset(filename string) ([]byte, error) {
	bytes, err := os.ReadFile(fmt.Sprintf("core/assets/%s", filename))
	if err != nil {
		fmt.Printf("Failed to read asset '%s'\n", filename)
		return nil, err
	}
	return bytes, nil
}

func LoadAssets() ([]fs.DirEntry, error) {
	files, err := os.ReadDir("core/assets")
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

const fileName = "save/game.json"

func LoadGameSave() (GameSave, error) {
	// Initialize an empty GameSave struct to return in case of error
	var gameSave GameSave

	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return gameSave, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Read the file contents
	jsonData, err := io.ReadAll(file)
	if err != nil {
		return gameSave, fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal JSON data into GameSave struct
	err = json.Unmarshal(jsonData, &gameSave)
	if err != nil {
		return gameSave, fmt.Errorf("error reading JSON data: %w", err)
	}

	return gameSave, nil
}

func WriteGameSame(gameSave GameSave) {
	jsonData, err := json.Marshal(gameSave)
	if err != nil {
		panic(fmt.Sprintf("Error writing struct to JSON: %s", err))
	}

	// Create or open the file
	file, err := os.Create(fileName)
	if err != nil {
		panic(fmt.Sprintf("Error creating file: %s", err))
	}
	defer file.Close() // Ensure the file is closed when done

	// Write JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		panic(fmt.Sprintf("Error writing JSON to file: %s", err))
	}
}
