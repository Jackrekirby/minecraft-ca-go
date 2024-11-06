package core

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"golang.org/x/image/colornames"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func DrawTextExample() {
	// Create a blank image with a white background
	width, height := 600, 300
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.White}, image.Point{}, draw.Src)

	// Load Cascadia Mono font with a larger font size
	fontFace, err := LoadTrueTypeFont("CascadiaMono.ttf", 48)
	if err != nil {
		log.Fatalf("failed to load font: %v", err)
	}

	// Draw text on the image
	DrawText(img, 50, 150, "Hello, Cascadia Mono!", colornames.Black, fontFace)

	// Save the image to a PNG file
	outFile, err := os.Create("output.png")
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer outFile.Close()
	png.Encode(outFile, img)
}

// loadTrueTypeFont loads a TTF font from a file and sets the font size.
func LoadTrueTypeFont(fontPath string, fontSize float64) (font.Face, error) {
	fontBytes, err := LoadAsset(fontPath)
	if err != nil {
		fmt.Println("Failed to read ttf file")
		return nil, err
	}
	ttf, err := opentype.Parse(fontBytes)
	if err != nil {
		fmt.Println("Failed to parse ttf file")
		return nil, err
	}
	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		fmt.Println("Failed to create face from ttf")
	}
	return face, err
}

// addLabel draws text on an image at the given position (x, y) with the specified font face and color
func DrawText(img *image.RGBA, x, y int, label string, col color.Color, face font.Face) {
	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot: fixed.Point26_6{
			X: fixed.I(x),
			Y: fixed.I(y),
		},
	}
	drawer.DrawString(label)
}
