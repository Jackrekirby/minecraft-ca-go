package core

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

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
