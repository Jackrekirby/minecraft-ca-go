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
	"path/filepath"

	"github.com/eiannone/keyboard"
)

type KeyboardManager struct {
	scene *Scene
}

// Initialise the keyboard manager
func (km *KeyboardManager) Initialise(scene *Scene) {
	// Open the keyboard
	err := keyboard.Open()
	if err != nil {
		panic(fmt.Sprintf("error opening keyboard: %v", err))
	}

	fmt.Println("Keyboard manager initialised. Press 'q' to quit.")
}

// Update method to handle keyboard events
// func (km *KeyboardManager) Update() {
// 	key, _, err := keyboard.GetKey()
// 	if err != nil {
// 		fmt.Println("Error reading key:", err)
// 		return
// 	}
// 	HandleKeyPress(km.scene, string(key))
// }

// Clean up resources (optional method)
func (km *KeyboardManager) Destroy() {
	keyboard.Close()
	fmt.Println("Keyboard manager closed.")
}

func KeyboardEvents(scene *Scene) {
	if isCPUProfiling {
		return
	}
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
		HandleKeyPress(scene, string(key), 0.5, DegToRad(5))
	}
}

func OutputSceneImage(img *image.RGBA) {
	if isCPUProfiling {
		return
	}
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
	RunEngine2(sceneImage, 1)
}

func FindProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
			return cwd, nil
		}

		parent := filepath.Dir(cwd)
		if parent == cwd { // Reached the root directory
			return "", fmt.Errorf("go.mod not found in any parent directory")
		}

		cwd = parent
	}
}

func CreateProjectRelativePath(relPath string) string {
	projectDir, err := FindProjectRoot()
	if err != nil {
		panic(fmt.Sprintf("Failed to find project dir %v\n", err))
	}
	absPath := filepath.Join(projectDir, relPath)
	return absPath
}

func LoadAsset(filename string) ([]byte, error) {
	relPath := fmt.Sprintf("core/assets/%s", filename)
	absPath := CreateProjectRelativePath(relPath)

	// Check if the file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Printf("Asset file does not exist: %s\n", absPath)
		return nil, err
	}

	// Read the file content
	bytes, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Printf("Failed to read asset '%s': %v\n", absPath, err)
		return nil, err
	}

	return bytes, nil
}

func OpenAsset(filename string) (*os.File, error) {
	relPath := fmt.Sprintf("core/assets/%s", filename)
	absPath := CreateProjectRelativePath(relPath)

	// Check if the file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Printf("Asset file does not exist: %s\n", absPath)
		return nil, err
	}

	// Read the file content
	file, err := os.Open(absPath)
	if err != nil {
		fmt.Printf("Failed to read asset '%s': %v\n", absPath, err)
		return nil, err
	}

	return file, nil
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
