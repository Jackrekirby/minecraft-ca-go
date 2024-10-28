package core

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/image/draw"
)

func DrawChequerBoard(c1 color.RGBA, c2 color.RGBA) {
	// Define the board size
	const (
		boardSize  = 8  // 8x8 squares
		squareSize = 50 // 50 pixels per square
		imageSize  = boardSize * squareSize
	)

	// Create a new image with the specified size
	img := image.NewRGBA(image.Rect(0, 0, imageSize, imageSize))

	// Loop through each square and set the color
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			// Determine the color based on the position
			var c color.Color
			if (x+y)%2 == 0 {
				c = c1
			} else {
				c = c2
			}

			// Fill in the square with the selected color
			for dy := 0; dy < squareSize; dy++ {
				for dx := 0; dx < squareSize; dx++ {
					img.Set(x*squareSize+dx, y*squareSize+dy, c)
				}
			}
		}
	}

	// Create the output file
	file, err := os.Create("chequerboard.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Encode as PNG
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}

func AnimateChequerBoard() {
	n := 10.0
	for i := 0; i < 10; i++ {
		fmt.Println(i)
		c := uint8(float64(i*255.0) / n)
		DrawChequerBoard(color.RGBA{c, 0, 0, 255}, color.RGBA{18, 91, 167, 255})
		time.Sleep(100 * time.Millisecond)
	}
}

func loadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func scaleImage(original image.Image, multiplier float64, interpolator draw.Interpolator) image.Image {
	// Calculate new dimensions
	newWidth := int(float64(original.Bounds().Dx()) * multiplier)
	newHeight := int(float64(original.Bounds().Dy()) * multiplier)
	scaledImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Draw the original image into the new scaled image
	interpolator.Scale(scaledImage, scaledImage.Rect, original, original.Bounds(), draw.Over, nil)

	return scaledImage
}

func GenerateTilemap(dir string, tileSize int) (image.Image, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var images []image.Image

	// Load all PNG files in the directory
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".png" {
			img, err := loadImage(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}
			// Scale image if necessary (optional)
			images = append(images, img)
		}
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("no PNG files found in directory: %s", dir)
	}

	// Calculate dimensions for the tilemap
	tilesPerDim := int(math.Ceil(math.Sqrt(float64(len(images)))))
	tilemapSize := tilesPerDim * tileSize

	fmt.Println(tilesPerDim, tilemapSize)

	// Create a new tilemap image
	tilemap := image.NewRGBA(image.Rect(0, 0, tilemapSize, tilemapSize))

	// Draw each image into the tilemap
	for i, img := range images {
		dstRect := image.Rect(
			(i%tilesPerDim)*tileSize, (i/tilesPerDim)*tileSize,
			(i%tilesPerDim+1)*tileSize, (i/tilesPerDim+1)*tileSize,
		)
		fmt.Println(i, dstRect)
		draw.Draw(tilemap, dstRect, img, image.Point{0, 0}, draw.Over)
	}

	return tilemap, nil
}

func SaveImage(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

func DrawWorld(worldSize int, imageNames []string) {
	blockSize := 32 // in pixels
	imageSize := worldSize*blockSize + 1

	img := image.NewRGBA(image.Rect(0, 0, imageSize, imageSize))

	backgroundColor := color.RGBA{18, 91, 167, 255}
	lineColor := color.RGBA{255, 255, 255, 255}

	for y := 0; y <= imageSize; y++ {
		for x := 0; x < imageSize; x++ {
			img.Set(x, y, backgroundColor)
		}
	}

	for y := 0; y <= worldSize; y++ {
		for x := 0; x < imageSize; x++ {
			img.Set(x, (worldSize*blockSize)-(y*blockSize), lineColor)
		}
	}

	for x := 0; x <= worldSize; x++ {
		for y := 0; y < imageSize; y++ {
			img.Set((worldSize*blockSize)-(x*blockSize), y, lineColor)
		}
	}

	// Load the tile image
	tileImage, err := loadImage("assets/test.png") // Change to your tile image path
	if err != nil {
		panic(err)
	}

	scaledTileImage := scaleImage(tileImage, 2, draw.NearestNeighbor)

	// Draw the tile image on each block in the grid
	for y := 0; y < worldSize; y++ {
		for x := 0; x < worldSize; x++ {
			dstRect := image.Rect(x*blockSize, y*blockSize, (x+1)*blockSize, (y+1)*blockSize)
			draw.Draw(img, dstRect, scaledTileImage, image.Point{0, 0}, draw.Over)
		}
	}

	file, err := os.Create("world.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}
