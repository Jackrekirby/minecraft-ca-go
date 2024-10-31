package core

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sort"

	"golang.org/x/image/math/fixed"
)

func DrawBox2D(x, y, size int, col color.RGBA, img *image.RGBA) {
	for i := -size / 2; i <= size/2; i++ {
		for j := -size / 2; j <= size/2; j++ {
			// Ensure we stay within image bounds
			px, py := x+i, y+j
			if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
				// img.Set(px, py, col)
				currentColor := img.RGBAAt(px, py)           // Get current color at (px, py)
				newColor := CombineColors(col, currentColor) // Combine new color with current color
				img.Set(px, py, newColor)                    // Set the combined color
			}
		}
	}
}

// DrawLine draws a line between two 2D points on the image using Bresenham's algorithm.
func DrawLine(img *image.RGBA, p1, p2 Point2D, col color.RGBA) {
	x0, y0 := int(p1.X), int(p1.Y)
	x1, y1 := int(p2.X), int(p2.Y)

	dx := int(math.Abs(float64(x1 - x0)))
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	dy := -int(math.Abs(float64(y1 - y0)))
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy

	for {
		// img.Set(x0, y0, col)
		DrawBox2D(x0, y0, 1, col, img)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func DrawLine3D(
	p0, p1 Point3D,
	camera Camera,
	img *image.RGBA,
	color color.RGBA,
) {
	// Get image dimensions
	imageWidth := float64(img.Bounds().Dx())
	imageHeight := float64(img.Bounds().Dy())

	// Project the points to 2D
	p2Start := camera.Convert3DTo2D(p0)
	p2End := camera.Convert3DTo2D(p1)

	if (p2Start.X == 0 && p2Start.Y == 0) || (p2End.X == 0 && p2End.Y == 0) {
		return
	}

	// Adjust coordinates to fit within the image bounds
	p2Start.X = p2Start.X*imageWidth/2 + imageWidth/2
	p2Start.Y = -p2Start.Y*imageHeight/2 + imageHeight/2
	p2End.X = p2End.X*imageWidth/2 + imageWidth/2
	p2End.Y = -p2End.Y*imageHeight/2 + imageHeight/2

	// Draw the line
	DrawLine(img, p2Start, p2End, color)
}

type Cuboid struct {
	vertices [8]Point3D
	Color    color.RGBA
}

func MakeAxisAlignedCuboid(min, max Point3D, color color.RGBA) Cuboid {
	vertices := [...]Point3D{
		{min.X, min.Y, min.Z}, // 0: Left-bottom-front
		{max.X, min.Y, min.Z}, // 1: Right-bottom-front
		{max.X, max.Y, min.Z}, // 2: Right-top-front
		{min.X, max.Y, min.Z}, // 3: Left-top-front
		{min.X, min.Y, max.Z}, // 4: Left-bottom-back
		{max.X, min.Y, max.Z}, // 5: Right-bottom-back
		{max.X, max.Y, max.Z}, // 6: Right-top-back
		{min.X, max.Y, max.Z}, // 7: Left-top-back
	}
	return Cuboid{vertices, color}
}

func (c Cuboid) Move(offset Point3D) Cuboid {
	for i := 0; i < 8; i++ {
		c.vertices[i] = c.vertices[i].Add(offset)
	}
	return c
}

func DrawCuboid(
	cuboid Cuboid,
	camera Camera,
	img *image.RGBA,
) {
	// Define the edges of the cube by connecting vertex indices
	edges := [][2]int{
		{0, 1}, {1, 2}, {2, 3}, {3, 0}, // Front face
		{4, 5}, {5, 6}, {6, 7}, {7, 4}, // Back face
		{0, 4}, {1, 5}, {2, 6}, {3, 7}, // Connecting edges
	}

	// Draw each edge of the cube
	for _, edge := range edges {
		DrawLine3D(
			cuboid.vertices[edge[0]],
			cuboid.vertices[edge[1]],
			camera,
			img, cuboid.Color,
		)
	}
}

// BlockDistance is a struct to hold block information along with distance from camera
type BlockDistance struct {
	Position Point3D
	Distance float64
}

// ByDistance implements sort.Interface for []BlockDistance based on the Distance field
type ByDistance []BlockDistance

func (a ByDistance) Len() int           { return len(a) }
func (a ByDistance) Less(i, j int) bool { return a[i].Distance < a[j].Distance }
func (a ByDistance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func GetSortedBlockPositions(origin Point3D) []BlockDistance {
	var blocks []BlockDistance
	gridSize := 16
	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {
			for z := 0; z < gridSize; z++ {
				p := Point3D{float64(x), float64(y), float64(z)}
				distance := Distance(p, origin)
				blocks = append(blocks, BlockDistance{Position: p, Distance: distance})
			}
		}
	}
	// Sort the blocks by distance from the camera
	sort.Sort(ByDistance(blocks))

	// Print sorted blocks
	// for _, bd := range blocks {
	// 	fmt.Printf("Block Position: (%v, %v, %v), Distance: %v\n", bd.Position.X, bd.Position.Y, bd.Position.Z, bd.Distance)
	// }
	return blocks
}

func Int26_6ToInt(value fixed.Int26_6) int {
	return int(value >> 6)
}

func DrawScene(scene *Scene) {
	imageSize := 512
	img := image.NewRGBA(image.Rect(0, 0, imageSize, imageSize))

	// Calculate aspect ratio based on the image dimensions

	blockSize := 1

	// gridSize := 16
	// fHalfGridSize := float64(gridSize) / 2.0

	// Constants for alpha scaling
	const maxAlpha = 1.0
	const minAlpha = 0.0
	const maxDistance = 16.0 // Distance at which alpha should be fully transparent
	const minDistance = 8.0  // Distance at which alpha is fully opaque

	for _, bd := range GetSortedBlockPositions(scene.Camera.Position) {
		p := Vec3{int(bd.Position.X), int(bd.Position.Y), int(bd.Position.Z)}
		block := scene.World.GetBlock(p)
		rb, isRenderable := block.(WireRenderBlock)
		localOffset := bd.Position
		cameraHorPos := scene.Camera.Position
		cameraHorPos.Y = localOffset.Y
		distance := Distance(localOffset, cameraHorPos)
		var alphaScaling float64

		if distance >= maxDistance {
			alphaScaling = minAlpha
		} else if distance <= minDistance {
			alphaScaling = maxAlpha
		} else {
			alphaScaling = maxAlpha * (1 - (distance / maxDistance))
		}

		if isRenderable {
			for _, c := range rb.ToCuboids() {

				mc := c.Move(localOffset)

				mc.Color.A = uint8(float64(mc.Color.A) * alphaScaling)
				// fmt.Println(mc, globalOffset, localOffset)
				DrawCuboid(mc, scene.Camera, img)
			}
		} else if p.Y == 0 {
			minP := Point3D{
				X: float64(p.X * blockSize),
				Y: float64(p.Y * blockSize),
				Z: float64(p.Z * blockSize),
			}
			maxP := Point3D{
				X: float64((p.X + 1) * blockSize),
				Y: float64((p.Y + 1) * blockSize),
				Z: float64((p.Z + 1) * blockSize),
			}
			c := MakeAxisAlignedCuboid(minP, maxP, color.RGBA{255, 255, 255, 100})
			c.Color.A = uint8(float64(c.Color.A) * alphaScaling)
			// DrawCuboid(c, scene.Camera, img)
		}
	}

	fontSize := Int26_6ToInt(scene.FontFace.Metrics().Height)
	DrawText(img, 4, fontSize,
		fmt.Sprintf("I: %d, U/I %d, sU/I %d, sI/I %d, S: %s",
			scene.Iteration,
			scene.NumBlockUpdatesInStep,
			scene.NumBlockSubUpdatesInStep,
			scene.NumBlockSubUpdateIterationsInStep,
			scene.GameState.String(),
		), Cyan.ToRGBA(), scene.FontFace)

	DrawText(img, 4, fontSize*2, fmt.Sprintf("%.1f, %.1f, %.1f", scene.Camera.Position.X, scene.Camera.Position.Y, scene.Camera.Position.Z), Cyan.ToRGBA(), scene.FontFace)

	// Create the output file
	file, err := os.Create("scene.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}
