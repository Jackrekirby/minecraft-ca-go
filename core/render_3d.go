package core

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

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
func DrawLine(img *image.RGBA, p0, p1 Int_2D, col color.RGBA) {
	dx := int(math.Abs(float64(p1.X - p0.X)))
	sx := -1
	if p0.X < p1.X {
		sx = 1
	}
	dy := -int(math.Abs(float64(p1.Y - p0.Y)))
	sy := -1
	if p0.Y < p1.Y {
		sy = 1
	}
	err := dx + dy

	for {
		// img.Set(p0.X, p0.Y, col)
		DrawBox2D(p0.X, p0.Y, 1, col, img)
		if p0.X == p1.X && p0.Y == p1.Y {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			p0.X += sx
		}
		if e2 <= dx {
			err += dx
			p0.Y += sy
		}
	}
}

func ProjectPoint(p Point3D, camera Camera, imageSize Point2D) *Int_2D {
	p2 := camera.Convert3DTo2D(p)
	if p2 == nil {
		return nil
	}
	int_2d := Int_2D{
		int(p2.X*imageSize.X/2 + imageSize.X/2),
		int(-p2.Y*imageSize.Y/2 + imageSize.Y/2),
	}
	return &int_2d
}

func DrawLine3D(
	p0, p1 Point3D,
	camera Camera,
	img *image.RGBA,
	color color.RGBA,
) {
	// Get image dimensions
	imageSize := Point2D{float64(img.Bounds().Dx()), float64(img.Bounds().Dy())}

	// Project the points to 2D
	p2Start := ProjectPoint(p0, camera, imageSize)
	p2End := ProjectPoint(p1, camera, imageSize)

	if p2Start != nil && p2End != nil {
		DrawLine(img, *p2Start, *p2End, color)
	}
}

func DrawTriangle2D(img *image.RGBA, p1, p2, p3 ImgPoint, col color.RGBA,
	depthBuffer *DepthBuffer, texture *image.RGBA,
) {
	imageSize := Int_2D{img.Bounds().Dx(), img.Bounds().Dy()}
	// Sort vertices by Y-coordinate (p1 is the topmost, p3 is the bottommost)
	if p2.Y < p1.Y {
		p1, p2 = p2, p1
	}
	if p3.Y < p1.Y {
		p1, p3 = p3, p1
	}
	if p3.Y < p2.Y {
		p2, p3 = p3, p2
	}

	fmt.Println(p1, p2, p3)

	// col = color.RGBA{uint8(p1.X * 10 % 255), uint8(p2.X * 10 % 255), uint8(p3.X * 10 % 255), 255}
	// Draw the triangle outline
	// DrawLine(img, p1, p2, White.ToRGBA())
	// DrawLine(img, p2, p3, White.ToRGBA())
	// DrawLine(img, p3, p1, White.ToRGBA())

	// return

	// Fill the triangle with a horizontal scanline
	yStart := max(int(p1.Y), 0)
	yEnd := min(int(p3.Y), imageSize.Y-1)

	for y := yStart; y <= yEnd; y++ {
		// var xStart, xEnd int
		// var zStart, zEnd float64
		var ip0, ip1 ImgPoint
		// Calculate x and z coordinates for each edge intersection at this Y level
		if y < int(p2.Y) {
			ip0 = interpolateImgPoint(p1, p3, y)
			ip1 = interpolateImgPoint(p1, p2, y)
			// xStart = interpolateX(p1, p3, y)
			// xEnd = interpolateX(p1, p2, y)
			// zStart = interpolateZ(p1, p3, y)
			// zEnd = interpolateZ(p1, p2, y)
		} else {
			ip0 = interpolateImgPoint(p1, p3, y)
			ip1 = interpolateImgPoint(p2, p3, y)
			// xStart = interpolateX(p1, p3, y)
			// xEnd = interpolateX(p2, p3, y)
			// zStart = interpolateZ(p1, p3, y)
			// zEnd = interpolateZ(p2, p3, y)
		}

		if ip0.X > ip1.X {
			ip0, ip1 = ip1, ip0
			// xStart, xEnd = xEnd, xStart
			// zStart, zEnd = zEnd, zStart
		}
		ip0.X = max(ip0.X, 0)
		ip1.X = min(ip1.X, imageSize.X-1)

		// Draw the horizontal line, interpolating depth along the line
		for x := ip0.X; x <= ip1.X; x++ {
			// Only draw if this pixel is closer than the current depth buffer value
			dbi := y*int(imageSize.X) + x
			if dbi > len(*depthBuffer) {
				fmt.Println(dbi, len(*depthBuffer), x, y, imageSize.X)
				panic("Index out of range")
			}

			t := float64(x-ip0.X) / float64(ip1.X-ip0.X)
			z := ip0.Z + t*(ip1.Z-ip0.Z)
			u := float64(ip0.U) + t*float64(ip1.U-ip0.U)
			v := float64(ip0.V) + t*float64(ip1.V-ip0.V)

			if z < (*depthBuffer)[dbi] {
				(*depthBuffer)[dbi] = z // Update the depth buffer
				// Sample the texture at the interpolated UV coordinates
				tx := int(u*float64(texture.Bounds().Dx())) % img.Bounds().Dx()
				ty := int(v*float64(texture.Bounds().Dy())) % img.Bounds().Dy()
				// fmt.Println(u, v, tx, ty)
				texColor := texture.At(tx, ty).(color.RGBA)
				// texColor := color.RGBA{uint8(u * 255.0), uint8(v * 255.0), 255, 255}
				img.Set(x, y, texColor)
			}
		}
	}
}

// Calculate the normal vector of the triangle given three vertices in 3D space.
func CalculateNormal(v1, v2, v3 Point3D) Point3D {
	// Calculate two edges of the triangle
	edge1 := Point3D{v2.X - v1.X, v2.Y - v1.Y, v2.Z - v1.Z}
	edge2 := Point3D{v3.X - v1.X, v3.Y - v1.Y, v3.Z - v1.Z}

	// Compute the cross product of the two edges
	normal := Point3D{
		edge1.Y*edge2.Z - edge1.Z*edge2.Y,
		edge1.Z*edge2.X - edge1.X*edge2.Z,
		edge1.X*edge2.Y - edge1.Y*edge2.X,
	}

	return Normalize(normal)
}

func ShadeColor(baseColor color.RGBA, intensity float64) color.RGBA {
	// Ensure intensity is clamped between 0 and 1
	if intensity < 0 {
		intensity = 0
	} else if intensity > 1 {
		intensity = 1
	}

	// Scale the color components by the intensity
	r := uint8(float64(baseColor.R) * intensity)
	g := uint8(float64(baseColor.G) * intensity)
	b := uint8(float64(baseColor.B) * intensity)
	return color.RGBA{r, g, b, baseColor.A}
}

type Vertex struct {
	Position Point3D
	U        int
	V        int
}

func DrawTriangle3D(
	v1, v2, v3 Vertex, // Three vertices of the triangle in 3D space
	camera Camera,
	img *image.RGBA,
	clr color.RGBA,
	depthBuffer *DepthBuffer,
	texture *image.RGBA,
) {
	// Calculate the normal of the triangle
	normal := CalculateNormal(v1.Position, v2.Position, v3.Position)
	avgV := v1.Position.Add(v2.Position).Add(v3.Position).Divide(Point3D{3, 3, 3})

	// RotateVector(Point3D{0, 0, 1}, camera.Rotation)
	// plane not in direction of camera
	if DotProduct(normal, Normalize(camera.Position.Subtract(avgV))) < 0 {
		return
	}

	imageSize := Point2D{float64(img.Bounds().Dx()), float64(img.Bounds().Dy())}
	// Project the 3D vertices to 2D screen coordinates
	p1 := ProjectPoint(v1.Position, camera, imageSize)
	p2 := ProjectPoint(v2.Position, camera, imageSize)
	p3 := ProjectPoint(v3.Position, camera, imageSize)

	// depthPerVertex := Point3D{Distance(v1, camera.Position), Distance(v2, camera.Position), Distance(v3, camera.Position)}

	// fmt.Println(v1, v2, v3, camera.Position, depthPerVertex)

	lightDirection := Normalize(Point3D{-0.3, 0.5, 0.8})

	// Calculate the shading intensity based on alignment with the light direction
	// Dot product between normal and light direction
	intensity := DotProduct(normal, lightDirection)
	// fmt.Println(normal, lightDirection, intensity)
	intensity = 0.7 + max(0, min(intensity, 1))*0.3

	// Shade the color based on the lighting intensity
	shadedColor := ShadeColor(clr, intensity)

	// xx := Point3D{255, 255, 255}.Multiply(normal)
	// shadedColor = color.RGBA{uint8(xx.X), uint8(xx.Y), uint8(xx.Z), uint8(255)}

	// Draw the triangle on the 2D screen using the projected points
	if p1 != nil && p2 != nil && p3 != nil {
		ip1 := ImgPoint{p1.X, p1.Y, Distance(v1.Position, camera.Position), v1.U, v1.V}
		ip2 := ImgPoint{p2.X, p2.Y, Distance(v2.Position, camera.Position), v2.U, v2.V}
		ip3 := ImgPoint{p3.X, p3.Y, Distance(v3.Position, camera.Position), v3.U, v3.V}
		DrawTriangle2D(img, ip1, ip2, ip3, shadedColor, depthBuffer, texture)
	}
}

// interpolateX calculates the X coordinate for a given Y using linear interpolation.
func interpolateX(p1, p2 ImgPoint, y int) int {
	if p1.Y == p2.Y || p1.X == p2.X { // Avoid division by zero if points are horizontal
		return p1.X
	}
	return int(p1.X + (y-p1.Y)*(p2.X-p1.X)/(p2.Y-p1.Y))
}

func interpolateZ(p1, p2 ImgPoint, y int) float64 {
	t := float64(y-p1.Y) / float64(p2.Y-p1.Y)
	return p1.Z + t*(p2.Z-p1.Z)
}

func interpolate(x0, x1, y0, y1, y float64) float64 {
	if y1 == y0 {
		return x0
	}
	t := (y - y0) / (y1 - y0)
	return x0 + t*(x1-x0)
}

func interpolateImgPoint(p1, p2 ImgPoint, y int) ImgPoint {
	f1 := float64(p1.Y)
	f2 := float64(p2.Y)
	fy := float64(y)
	p := ImgPoint{
		X: int(interpolate(float64(p1.X), float64(p2.X), f1, f2, fy)),
		Z: interpolate(float64(p1.Z), float64(p2.Z), f1, f2, fy),
		U: int(interpolate(float64(p1.U), float64(p2.U), f1, f2, fy)),
		V: int(interpolate(float64(p1.V), float64(p2.V), f1, f2, fy)),
	}
	return p
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

type DepthBuffer []float64

func DrawFilledCuboid(
	cuboid Cuboid,
	camera Camera,
	img *image.RGBA,
	depthBuffer *DepthBuffer,
	texture *image.RGBA,
) {
	// Define the faces of the cube with four vertices each
	// faces := [][4]int{
	// 	{0, 3, 2, 1}, // Front face
	// 	{4, 5, 6, 7}, // Back face
	// 	{0, 1, 5, 4}, // Top face
	// 	{2, 3, 7, 6}, // Bottom face
	// 	{7, 3, 0, 4}, // Left face
	// 	{1, 2, 6, 5}, // Right face
	// }

	// // Draw each face of the cube using two triangles
	// for _, face := range faces {
	// 	// Split each face into two triangles and draw
	// 	DrawTriangle3D(
	// 		cuboid.vertices[face[0]],
	// 		cuboid.vertices[face[1]],
	// 		cuboid.vertices[face[2]],
	// 		camera,
	// 		img, cuboid.Color, depthBuffer,
	// 		texture,
	// 	)
	// 	DrawTriangle3D(
	// 		cuboid.vertices[face[0]],
	// 		cuboid.vertices[face[2]],
	// 		cuboid.vertices[face[3]],
	// 		camera,
	// 		img, cuboid.Color, depthBuffer,
	// 		texture,
	// 	)
	// }
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
func (a ByDistance) Less(i, j int) bool { return a[i].Distance > a[j].Distance }
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
	// sort.Sort(ByDistance(blocks))

	// Print sorted blocks
	// for _, bd := range blocks {
	// 	fmt.Printf("Block Position: (%v, %v, %v), Distance: %v\n", bd.Position.X, bd.Position.Y, bd.Position.Z, bd.Distance)
	// }
	return blocks
}

func Int26_6ToInt(value fixed.Int26_6) int {
	return int(value >> 6)
}

func DrawObjects(scene *Scene, img *image.RGBA, depthBuffer *DepthBuffer) {
	// Calculate aspect ratio based on the image dimensions

	blockSize := 1

	// gridSize := 16
	// fHalfGridSize := float64(gridSize) / 2.0

	// Constants for alpha scaling
	// const maxAlpha = 1.0
	// const minAlpha = 0.0
	// const maxDistance = 16.0 // Distance at which alpha should be fully transparent
	// const minDistance = 8.0  // Distance at which alpha is fully opaque
	// fmt.Println("Drawing", scene.Iteration)
	for _, bd := range GetSortedBlockPositions(scene.Camera.Position) {
		p := Vec3{int(bd.Position.X), int(bd.Position.Y), int(bd.Position.Z)}
		block := scene.World.GetBlock(p)
		rb, isRenderable := block.(WireRenderBlock)
		localOffset := bd.Position
		// cameraHorPos := scene.Camera.Position
		// cameraHorPos.Y = localOffset.Y
		// distance := Distance(localOffset, cameraHorPos)
		var alphaScaling float64 = 1.0

		// if distance >= maxDistance {
		// 	alphaScaling = minAlpha
		// } else if distance <= minDistance {
		// 	alphaScaling = maxAlpha
		// } else {
		// 	alphaScaling = maxAlpha * (1 - (distance / maxDistance))
		// }

		if isRenderable {
			for _, c := range rb.ToCuboids() {

				mc := c.Move(localOffset)

				mc.Color.A = uint8(float64(mc.Color.A) * alphaScaling)
				// fmt.Println(mc, globalOffset, localOffset)
				DrawFilledCuboid(mc, scene.Camera, img, depthBuffer, &scene.Texture)
				// fmt.Println(mc)
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
}

func DrawScene(scene *Scene) {
	imageSize := 512
	img := image.NewRGBA(image.Rect(0, 0, imageSize, imageSize))

	depthBuffer := make(DepthBuffer, imageSize*imageSize)
	for i := range depthBuffer {
		depthBuffer[i] = 1e9 // A large value representing 'infinity'
	}

	for x := 0; x < imageSize; x++ {
		for y := 0; y < imageSize; y++ {
			img.Set(x, y, color.RGBA{122, 168, 253, 255})
		}
	}

	k := 10.0
	//front
	DrawTriangle3D(
		Vertex{Point3D{0, 0, 0}, 0, 0},
		Vertex{Point3D{k, 0, 0}, 1, 0},
		Vertex{Point3D{k, k, 0}, 1, 1},
		scene.Camera, img, Red.ToRGBA(), &depthBuffer, &scene.Texture,
	)

	// DrawTriangle3D(Point3D{0, 0, 0}, Point3D{k, k, 0}, Point3D{0, k, 0}, scene.Camera, img, Orange.ToRGBA(), &depthBuffer)

	// // left
	// DrawTriangle3D(Point3D{0, 0, 0}, Point3D{0, 0, k}, Point3D{0, k, k}, scene.Camera, img, Yellow.ToRGBA(), &depthBuffer)
	// DrawTriangle3D(Point3D{0, 0, 0}, Point3D{0, k, k}, Point3D{0, k, 0}, scene.Camera, img, Green.ToRGBA(), &depthBuffer)

	// // back
	// DrawTriangle3D(Point3D{0, 0, k}, Point3D{k, k, k}, Point3D{k, 0, k}, scene.Camera, img, Lime.ToRGBA(), &depthBuffer)

	// DrawTriangle3D(Point3D{0, 0, k}, Point3D{0, k, k}, Point3D{k, k, k}, scene.Camera, img, LightBlue.ToRGBA(), &depthBuffer)

	// // right
	// DrawTriangle3D(Point3D{k, 0, 0}, Point3D{k, 0, k}, Point3D{k, k, k}, scene.Camera, img, Cyan.ToRGBA(), &depthBuffer)
	// DrawTriangle3D(Point3D{k, 0, 0}, Point3D{k, k, k}, Point3D{k, k, 0}, scene.Camera, img, Blue.ToRGBA(), &depthBuffer)

	//DrawObjects(scene, img, &depthBuffer)

	fontSize := Int26_6ToInt(scene.FontFace.Metrics().Height)
	DrawText(img, 4, fontSize,
		fmt.Sprintf("I: %d, U/I %d, sU/I %d, sI/I %d",
			scene.Iteration,
			scene.NumBlockUpdatesInStep,
			scene.NumBlockSubUpdatesInStep,
			scene.NumBlockSubUpdateIterationsInStep,
		), Cyan.ToRGBA(), scene.FontFace)

	DrawText(img, 4, fontSize*2,
		fmt.Sprintf("F/S: %d, I/S %d, S: %s",
			scene.RecordedFramesPerSecond,
			scene.RecordedStepsPerSecond,
			scene.GameState.String(),
		), Cyan.ToRGBA(), scene.FontFace)

	DrawText(img, 4, fontSize*3, fmt.Sprintf("XYZ: %.1f, %.1f, %.1f", scene.Camera.Position.X, scene.Camera.Position.Y, scene.Camera.Position.Z), Cyan.ToRGBA(), scene.FontFace)

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
