package samples

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

// Point represents a 2D point or vector
type Point struct {
	X, Y int
}

// TexCoord represents 2D UV coordinates in texture space
type TexCoord struct {
	U, V float64
}

// Vertex defines a vertex with position and texture coordinates
type Vertex struct {
	Pos Point
	UV  TexCoord
}

// Triangle defines a triangle with three vertices
type Triangle struct {
	V0, V1, V2 Vertex
}

// drawTexturedTriangle rasterizes a triangle with texture mapping
func drawTexturedTriangle(img *image.RGBA, tri Triangle, texture *image.RGBA) {
	// Sort vertices by Y-coordinate (v0.Y <= v1.Y <= v2.Y)
	v := []Vertex{tri.V0, tri.V1, tri.V2}
	if v[0].Pos.Y > v[1].Pos.Y {
		v[0], v[1] = v[1], v[0]
	}
	if v[1].Pos.Y > v[2].Pos.Y {
		v[1], v[2] = v[2], v[1]
	}
	if v[0].Pos.Y > v[1].Pos.Y {
		v[0], v[1] = v[1], v[0]
	}

	// Check for flat-bottom and flat-top triangles
	if v[1].Pos.Y == v[2].Pos.Y {
		fillTexturedFlatBottomTriangle(img, v[0], v[1], v[2], texture)
	} else if v[0].Pos.Y == v[1].Pos.Y {
		fillTexturedFlatTopTriangle(img, v[0], v[1], v[2], texture)
	} else {
		// Split the triangle into a flat-bottom and a flat-top
		v4 := Vertex{
			Pos: Point{
				X: v[0].Pos.X + int(float64(v[2].Pos.X-v[0].Pos.X)*float64(v[1].Pos.Y-v[0].Pos.Y)/float64(v[2].Pos.Y-v[0].Pos.Y)),
				Y: v[1].Pos.Y,
			},
			UV: TexCoord{
				U: v[0].UV.U + (v[2].UV.U-v[0].UV.U)*float64(v[1].Pos.Y-v[0].Pos.Y)/float64(v[2].Pos.Y-v[0].Pos.Y),
				V: v[0].UV.V + (v[2].UV.V-v[0].UV.V)*float64(v[1].Pos.Y-v[0].Pos.Y)/float64(v[2].Pos.Y-v[0].Pos.Y),
			},
		}
		fillTexturedFlatBottomTriangle(img, v[0], v[1], v4, texture)
		fillTexturedFlatTopTriangle(img, v[1], v4, v[2], texture)
	}
}

// fillTexturedFlatBottomTriangle fills a flat-bottom triangle with texture mapping
func fillTexturedFlatBottomTriangle(img *image.RGBA, v0, v1, v2 Vertex, texture *image.RGBA) {
	slope1 := float64(v1.Pos.X-v0.Pos.X) / float64(v1.Pos.Y-v0.Pos.Y)
	slope2 := float64(v2.Pos.X-v0.Pos.X) / float64(v2.Pos.Y-v0.Pos.Y)

	startX := float64(v0.Pos.X)
	endX := float64(v0.Pos.X)
	startUV := v0.UV
	endUV := v0.UV

	fmt.Println(startX, endX, v0.Pos.Y, v1.Pos.Y)

	for y := v0.Pos.Y; y <= v1.Pos.Y; y++ {
		lerpTexCoords(img, int(startX), int(endX), y, startUV, endUV, texture)
		startX += slope1
		endX += slope2

		// Linear interpolation of UV coordinates
		t := float64(y-v0.Pos.Y) / float64(v1.Pos.Y-v0.Pos.Y)
		startUV = lerpTexCoord(v0.UV, v1.UV, t)
		endUV = lerpTexCoord(v0.UV, v2.UV, t)
	}
}

// fillTexturedFlatTopTriangle fills a flat-top triangle with texture mapping
func fillTexturedFlatTopTriangle(img *image.RGBA, v0, v1, v2 Vertex, texture *image.RGBA) {
	slope1 := float64(v2.Pos.X-v0.Pos.X) / float64(v2.Pos.Y-v0.Pos.Y)
	slope2 := float64(v2.Pos.X-v1.Pos.X) / float64(v2.Pos.Y-v1.Pos.Y)

	startX := float64(v2.Pos.X)
	endX := float64(v2.Pos.X)
	startUV := v2.UV
	endUV := v2.UV

	for y := v2.Pos.Y; y > v0.Pos.Y; y-- {
		lerpTexCoords(img, int(startX), int(endX), y, startUV, endUV, texture)
		startX -= slope1
		endX -= slope2

		// Linear interpolation of UV coordinates
		t := float64(v2.Pos.Y-y) / float64(v2.Pos.Y-v0.Pos.Y)
		startUV = lerpTexCoord(v2.UV, v0.UV, t)
		endUV = lerpTexCoord(v2.UV, v1.UV, t)
	}
}

// lerpTexCoords fills a scanline between two points with interpolated UV coordinates
func lerpTexCoords(img *image.RGBA, xStart, xEnd, y int, uvStart, uvEnd TexCoord, texture *image.RGBA) {
	for x := xStart; x <= xEnd; x++ {
		t := float64(x-xStart) / float64(xEnd-xStart)
		uv := lerpTexCoord(uvStart, uvEnd, t)

		// Sample the color from the texture
		tx := int(uv.U * float64(texture.Bounds().Dx()))
		ty := int(uv.V * float64(texture.Bounds().Dy()))
		color := texture.At(tx, ty)
		fmt.Println(x, y, color)
		img.Set(x, y, color)
	}
}

// lerpTexCoord performs linear interpolation between two texture coordinates
func lerpTexCoord(tc0, tc1 TexCoord, t float64) TexCoord {
	return TexCoord{
		U: tc0.U + t*(tc1.U-tc0.U),
		V: tc0.V + t*(tc1.V-tc0.V),
	}
}

func ImageToRGBA(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	return rgba
}

func Main() {
	// Load the texture
	textureFile, err := os.Open("assets/test.png")
	if err != nil {
		panic(err)
	}
	defer textureFile.Close()

	textureImg, err := png.Decode(textureFile)
	if err != nil {
		panic(err)
	}

	// Convert texture to RGBA
	texture := ImageToRGBA(textureImg)

	// Create a blank RGBA image
	width, height := 300, 300
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Define a triangle with texture coordinates
	triangle := Triangle{
		V0: Vertex{Pos: Point{0, 0}, UV: TexCoord{0.0, 0.0}},
		V2: Vertex{Pos: Point{0, 300}, UV: TexCoord{0.0, 1.0}},
		V1: Vertex{Pos: Point{300, 150}, UV: TexCoord{1.0, 0.5}},
	}

	// Rasterize the triangle with the texture
	drawTexturedTriangle(img, triangle, texture)

	// Save the image to a file
	file, err := os.Create("textured_triangle.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}
