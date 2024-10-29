package core

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/eiannone/keyboard"
)

// Point3D represents a point in 3D space.
type Point3D struct {
	X, Y, Z float64
}

func Point3DFromScalar(s float64) Point3D {
	return Point3D{s, s, s}
}

func (v Point3D) Add(other Point3D) Point3D {
	return Point3D{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
	}
}

func (v Point3D) Subtract(other Point3D) Point3D {
	return Point3D{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
	}
}

func (v Point3D) Divide(other Point3D) Point3D {
	return Point3D{
		X: v.X / other.X,
		Y: v.Y / other.Y,
		Z: v.Z / other.Z,
	}
}

func Distance(p1, p2 Point3D) float64 {
	p := p2.Subtract(p1)
	// fmt.Println(p1, p2, p, math.Sqrt(p.X*p.X+p.Y*p.Y+p.Z*p.Z))
	return math.Sqrt(p.X*p.X + p.Y*p.Y + p.Z*p.Z)
}

func DegToRad(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// RotateX returns a new Point3D rotated around the X-axis by the given angle (in radians)
func (p Point3D) RotateX(angle float64) Point3D {
	y := p.Y*math.Cos(angle) - p.Z*math.Sin(angle)
	z := p.Y*math.Sin(angle) + p.Z*math.Cos(angle)
	return Point3D{X: p.X, Y: y, Z: z}
}

// RotateY returns a new Point3D rotated around the Y-axis by the given angle (in radians)
func (p Point3D) RotateY(angle float64) Point3D {
	x := p.X*math.Cos(angle) + p.Z*math.Sin(angle)
	z := -p.X*math.Sin(angle) + p.Z*math.Cos(angle)
	return Point3D{X: x, Y: p.Y, Z: z}
}

// RotateZ returns a new Point3D rotated around the Z-axis by the given angle (in radians)
func (p Point3D) RotateZ(angle float64) Point3D {
	x := p.X*math.Cos(angle) - p.Y*math.Sin(angle)
	y := p.X*math.Sin(angle) + p.Y*math.Cos(angle)
	return Point3D{X: x, Y: y, Z: p.Z}
}

// Point2D represents a point in 2D space.
type Point2D struct {
	X, Y float64
}

// Camera represents the camera in 3D space, with position, rotation, and projection parameters
type Camera struct {
	Position    Point3D
	Rotation    Point3D
	FOV         float64
	AspectRatio float64
	Near        float64
	Far         float64
}

// Convert3DTo2D projects a 3D point onto a 2D plane using perspective projection, considering the camera's position and orientation
func (c *Camera) Convert3DTo2D(point Point3D) Point2D {
	// Translate the point by the camera position
	p := point.
		Subtract(c.Position).
		RotateX(c.Rotation.X).
		RotateY(c.Rotation.Y).
		RotateZ(c.Rotation.Z)

	// Check if the point is behind the camera
	if p.Z <= 0 {
		return Point2D{X: 0, Y: 0}
	}

	// Perspective projection
	fovRad := 1.0 / math.Tan(c.FOV*0.5*math.Pi/180)
	q := c.Far / (c.Far - c.Near)

	ndcX := p.X * fovRad * c.AspectRatio
	ndcY := p.Y * fovRad
	ndcZ := p.Z * q

	// Project to screen coordinates
	screenX := ndcX / ndcZ
	screenY := ndcY / ndcZ

	return Point2D{X: screenX, Y: screenY}
}

// Convert3DTo2D projects a 3D point onto a 2D plane using perspective projection.
func Convert3DTo2D(point Point3D, fov, aspectRatio, near, far float64) Point2D {
	fovRad := 1.0 / math.Tan(fov*0.5*math.Pi/180)
	q := far / (far - near)

	ndcX := point.X * fovRad * aspectRatio
	ndcY := point.Y * fovRad
	ndcZ := point.Z * q

	if point.Z <= 0 {
		return Point2D{X: 0, Y: 0}
	}

	screenX := ndcX / ndcZ
	screenY := ndcY / ndcZ

	return Point2D{X: screenX, Y: screenY}
}

func CombineColors(src color.RGBA, dst color.RGBA) color.RGBA {
	// Calculate the new alpha
	alpha := float64(src.A) / 255.0
	newA := uint8(float64(src.A) + float64(dst.A)*(1-alpha))

	// If the new alpha is zero, return transparent black
	if newA == 0 {
		return color.RGBA{0, 0, 0, 0}
	}

	// Blend the colors based on the alpha
	newR := uint8((float64(dst.R)*(1-alpha) + float64(src.R)*alpha) * (float64(newA) / 255.0))
	newG := uint8((float64(dst.G)*(1-alpha) + float64(src.G)*alpha) * (float64(newA) / 255.0))
	newB := uint8((float64(dst.B)*(1-alpha) + float64(src.B)*alpha) * (float64(newA) / 255.0))

	return color.RGBA{R: newR, G: newG, B: newB, A: newA}
}

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

func KeyboardEvents(camera *Camera, quitGame *bool) {
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
		delta := 1.0
		rotation := DegToRad(15)
		// Handle key press
		switch key {
		case 'q':
			fmt.Println("Exiting...")
			*quitGame = true
			return
		case 'w':
			camera.Position = camera.Position.Add(Point3D{0, 0, delta}.RotateY(-camera.Rotation.Y))
		case 'a':
			camera.Position = camera.Position.Add(Point3D{-delta, 0, 0}.RotateY(-camera.Rotation.Y))
		case 's':
			camera.Position = camera.Position.Add(Point3D{0, 0, -delta}.RotateY(-camera.Rotation.Y))
		case 'd':
			camera.Position = camera.Position.Add(Point3D{delta, 0, 0}.RotateY(-camera.Rotation.Y))
		case 'e':
			camera.Position = camera.Position.Add(Point3D{0, 1, 0})
		case 'c':
			camera.Position = camera.Position.Add(Point3D{0, -1, 0})
		case 'z':
			camera.Rotation.Y = camera.Rotation.Y + rotation
		case 'x':
			camera.Rotation.Y = camera.Rotation.Y - rotation
		default:
			fmt.Println("Pressed:", key)
		}
	}
}

func DrawScene(camera *Camera, world *World) {
	imageSize := 512
	img := image.NewRGBA(image.Rect(0, 0, imageSize, imageSize))

	// Calculate aspect ratio based on the image dimensions

	blockSize := 1

	gridSize := 16
	// fHalfGridSize := float64(gridSize) / 2.0

	// Constants for alpha scaling
	const maxAlpha = 1.0
	const minAlpha = 0.0
	const maxDistance = 16.0 // Distance at which alpha should be fully transparent
	const minDistance = 8.0  // Distance at which alpha is fully opaque

	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {
			for z := 0; z < gridSize; z++ {
				block := world.GetBlock(Vec3{x, y, z})
				rb, isRenderable := block.(WireRenderBlock)
				localOffset := Point3D{float64(x), float64(y), float64(z)}
				cameraHorPos := camera.Position
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
						DrawCuboid(mc, *camera, img)
					}
				} else if y == 0 {
					minP := Point3D{
						X: float64(x * blockSize),
						Y: float64(y * blockSize),
						Z: float64(z * blockSize),
					}
					maxP := Point3D{
						X: float64((x + 1) * blockSize),
						Y: float64((y + 1) * blockSize),
						Z: float64((z + 1) * blockSize),
					}
					c := MakeAxisAlignedCuboid(minP, maxP, color.RGBA{255, 255, 255, 100})
					c.Color.A = uint8(float64(c.Color.A) * alphaScaling)
					DrawCuboid(c, *camera, img)
				}
			}
		}
	}

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
