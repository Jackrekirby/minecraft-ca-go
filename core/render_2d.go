package core

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
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

func ImageToRGBA(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	return rgba
}

func loadImage(filePath string) (image.Image, error) {
	file, err := LoadAsset(filePath)
	if err != nil {
		return nil, err
	}
	// defer file.Close()

	reader := bytes.NewReader(file)

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func scaleImage(original *image.RGBA, multiplier float64, interpolator draw.Interpolator) *image.RGBA {
	// Calculate new dimensions
	newWidth := int(float64(original.Bounds().Dx()) * multiplier)
	newHeight := int(float64(original.Bounds().Dy()) * multiplier)
	scaledImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Draw the original image into the new scaled image
	interpolator.Scale(scaledImage, scaledImage.Rect, original, original.Bounds(), draw.Over, nil)

	return scaledImage
}

type TextureMeta struct {
	// normalised coordinates
	U      float64
	V      float64
	Width  float64
	Height float64
}

type Tilemap struct {
	Image image.RGBA
	Metas map[string]TextureMeta
}

// TextureMeta and Tilemap structs defined as before

func GenerateTilemap(dir string, tileSize int) (*Tilemap, error) {
	const doPad = false
	files, err := LoadAssets()
	if err != nil {
		return nil, err
	}

	var images []image.Image
	textureMetas := make(map[string]TextureMeta)
	var imageFiles []fs.DirEntry

	// Load all PNG files in the directory
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".png" {
			imageFiles = append(imageFiles, file)
		}
	}

	for _, file := range imageFiles {
		img, err := loadImage(file.Name())
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("no PNG files found in directory: %s", dir)
	}

	padWidth := 0
	if doPad {
		padWidth = 1
	}

	// Calculate dimensions for the tilemap with border
	tilesPerDim := int(math.Ceil(math.Sqrt(float64(len(images)))))
	tilemapSize := tilesPerDim * (tileSize + 2*padWidth) // Each tile includes a 1-pixel border on each side

	// Create a new tilemap image with adjusted size
	agg_image := image.NewRGBA(image.Rect(0, 0, tilemapSize, tilemapSize))

	// Draw each image into the tilemap and calculate UV coordinates
	for i, img := range images {
		// Calculate destination rectangle with border margin

		x := (i%tilesPerDim)*(tileSize+2*padWidth) + padWidth // Leave 1-pixel margin for the border
		y := (i/tilesPerDim)*(tileSize+2*padWidth) + padWidth
		dstRect := image.Rect(x, y, x+tileSize, y+tileSize)

		// Draw the tile itself
		draw.Draw(agg_image, dstRect, img, image.Point{0, 0}, draw.Over)

		if doPad {
			// Fill the border pixels around each tile
			for bx := 0; bx < tileSize; bx++ {
				// Top border
				agg_image.Set(x+bx, y-1, img.At(bx, 0))
				// Bottom border
				agg_image.Set(x+bx, y+tileSize, img.At(bx, tileSize-1))
			}
			for by := 0; by < tileSize; by++ {
				// Left border
				agg_image.Set(x-1, y+by, img.At(0, by))
				// Right border
				agg_image.Set(x+tileSize, y+by, img.At(tileSize-1, by))
			}

			// Fill corners of the border with neighboring pixel colors
			agg_image.Set(x-1, y-1, img.At(0, 0))                                 // Top-left corner
			agg_image.Set(x+tileSize, y-1, img.At(tileSize-1, 0))                 // Top-right corner
			agg_image.Set(x-1, y+tileSize, img.At(0, tileSize-1))                 // Bottom-left corner
			agg_image.Set(x+tileSize, y+tileSize, img.At(tileSize-1, tileSize-1)) // Bottom-right corner
		}

		// Calculate normalized UV coordinates
		u := float64(x) / float64(tilemapSize)
		v := float64(y) / float64(tilemapSize)
		uWidth := float64(tileSize) / float64(tilemapSize)
		vHeight := float64(tileSize) / float64(tilemapSize)

		// Map filename to UV coordinates
		filename := imageFiles[i].Name()
		nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
		textureMetas[nameWithoutExt] = TextureMeta{u, v, uWidth, vHeight}
	}

	return &Tilemap{*agg_image, textureMetas}, nil
}

func SaveImage(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// func DrawWorld(worldSize int, imageNames []string) {
// 	blockSize := 32 // in pixels
// 	imageSize := worldSize*blockSize + 1

// 	img := image.NewRGBA(image.Rect(0, 0, imageSize, imageSize))

// 	backgroundColor := color.RGBA{18, 91, 167, 255}
// 	lineColor := color.RGBA{255, 255, 255, 255}

// 	for y := 0; y <= imageSize; y++ {
// 		for x := 0; x < imageSize; x++ {
// 			img.Set(x, y, backgroundColor)
// 		}
// 	}

// 	for y := 0; y <= worldSize; y++ {
// 		for x := 0; x < imageSize; x++ {
// 			img.Set(x, (worldSize*blockSize)-(y*blockSize), lineColor)
// 		}
// 	}

// 	for x := 0; x <= worldSize; x++ {
// 		for y := 0; y < imageSize; y++ {
// 			img.Set((worldSize*blockSize)-(x*blockSize), y, lineColor)
// 		}
// 	}

// 	// Load the tile image
// 	tileImage, err := loadImage("assets/test.png") // Change to your tile image path
// 	if err != nil {
// 		panic(err)
// 	}

// 	scaledTileImage := scaleImage(tileImage, 2, draw.NearestNeighbor)

// 	// Draw the tile image on each block in the grid
// 	for y := 0; y < worldSize; y++ {
// 		for x := 0; x < worldSize; x++ {
// 			dstRect := image.Rect(x*blockSize, y*blockSize, (x+1)*blockSize, (y+1)*blockSize)
// 			draw.Draw(img, dstRect, scaledTileImage, image.Point{0, 0}, draw.Over)
// 		}
// 	}

// 	file, err := os.Create("world.png")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()

// 	if err := png.Encode(file, img); err != nil {
// 		panic(err)
// 	}
// }
