//go:build !js && !wasm
// +build !js,!wasm

package core

import (
	"image"
	"image/color"
	"testing"
)

func TestDDA3D(t *testing.T) {
	imageSize := 129
	cellSize := 16
	img := image.NewRGBA(image.Rect(0, 0, imageSize, imageSize))

	fCellSize := float64(cellSize)
	for i := 0; i <= imageSize; i += cellSize {
		DrawLine(
			img,
			Int_2D{i, 0},
			Int_2D{i, imageSize - 1},
			color.RGBA{100, 100, 100, 255},
		)
		DrawLine(
			img,
			Int_2D{0, i},
			Int_2D{imageSize - 1, i},
			color.RGBA{100, 100, 100, 255},
		)
	}

	var pointGroups [4][]Point3D = [...][]Point3D{
		DDA3D(4.1, 0.2, 0, 0.8, 4.3, 0),
		DDA3D(0.8, 4.5, 0, 4.4, 7.8, 0),
		DDA3D(4.6, 7.8, 0, 7.8, 4.7, 0),
		DDA3D(7.8, 4.3, 0, 4.5, 0.2, 0),
	}

	for _, points := range pointGroups {
		n := len(points)
		// fmt.Println(i, n)
		// fmt.Println(i, points)
		for i := 1; i < n; i++ {
			p0, p1 := points[i-1], points[i]
			// r, g, b := HSVToRGB(float64(i)/float64(n)*360.0, 1.0, 1.0)
			// clr := color.RGBA{r, g, b, 255}
			DrawLine(
				img,
				Int_2D{int(p0.X * fCellSize), int(p0.Y * fCellSize)},
				Int_2D{int(p1.X * fCellSize), int(p1.Y * fCellSize)},
				color.RGBA{0, 255, 0, 255},
			)
			img.Set(int(p0.X*fCellSize), int(p0.Y*fCellSize), color.RGBA{255, 0, 0, 255})
			img.Set(int(p1.X*fCellSize), int(p1.Y*fCellSize), color.RGBA{255, 0, 0, 255})
		}
	}
	// points := DDA3D(4.1, 0.2, 0, 0.8, 4.3, 0)

	SaveImage(img, CreateProjectRelativePath("output/TestDDA3D.png"))
}
