package core

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"golang.org/x/image/math/fixed"
)

const skipAdjacentFaces = true // false = slow
const drawTriangleWire = false // true = slow
const DebugUV = false
const TexturePerspective = true

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
		img.Set(p0.X, p0.Y, col)
		// DrawBox2D(p0.X, p0.Y, 1, col, img)
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

func GenerateLine(p0, p1 Int_2D, callback func(int, int)) {
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
		callback(p0.X, p0.Y)
		// DrawBox2D(p0.X, p0.Y, 1, col, img)
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

func calculateDepth(point Point3D, c Camera) float64 {
	// should be returned by project point...
	// Translate the point by the camera position
	p := point. // transform order must match project point code
			Subtract(c.Position).
			RotateY(c.Rotation.Y).
			RotateX(c.Rotation.X).
			RotateZ(c.Rotation.Z)

	// Perspective projection
	q := c.Far / (c.Far - c.Near)
	ndcZ := p.Z * q
	return ndcZ
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

	// fmt.Println(p1, p2, p3)

	// Draw the triangle outline
	// DrawLine(img, p1, p2, White.ToRGBA())
	// DrawLine(img, p2, p3, White.ToRGBA())
	// DrawLine(img, p3, p1, White.ToRGBA())

	// return

	// Fill the triangle with a horizontal scanline
	yStart := max(int(p1.Y), 0)
	yEnd := min(int(p3.Y), imageSize.Y-1)

	for y := yStart; y <= yEnd; y++ {
		var ip0, ip1 ImgPoint
		// Calculate x and z coordinates for each edge intersection at this Y level
		if y < int(p2.Y) {
			ip0 = interpolateImgPoint(p1, p3, y)
			ip1 = interpolateImgPoint(p1, p2, y)

			// rp0 = p3
			// rp1 = p2
		} else {
			ip0 = interpolateImgPoint(p1, p3, y)
			ip1 = interpolateImgPoint(p2, p3, y)
			// rp0 = p1
			// rp1 = p2
		}

		if ip0.X > ip1.X {
			ip0, ip1 = ip1, ip0
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
			// t2 := float64(x-rp0.X) / float64(rp1.X-rp0.X)
			// u := float64(ip0.U) + t*float64(ip1.U-ip0.U)
			// v := float64(ip0.V) + t*float64(ip1.V-ip0.V)

			if z < (*depthBuffer)[dbi] {
				(*depthBuffer)[dbi] = z // Update the depth buffer
				// Sample the texture at the interpolated UV coordinates
				// tx := int(u*float64(texture.Bounds().Dx())) % img.Bounds().Dx()
				// ty := int(v*float64(texture.Bounds().Dy())) % img.Bounds().Dy()
				// fmt.Println(u, v, tx, ty)
				// texColor := texture.At(tx, ty).(color.RGBA)
				// texColor := color.RGBA{uint8(u * 255.0), uint8(v * 255.0), 255, 255}
				img.Set(x, y, col)
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
	U        float64
	V        float64
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

	// plane not in direction of camera
	if DotProduct(normal, Normalize(camera.Position.Subtract(avgV))) < 0 {
		return
	}
	imageSize := Point2D{float64(img.Bounds().Dx()), float64(img.Bounds().Dy())}
	// Project the 3D vertices to 2D screen coordinates
	// p1 := ProjectPoint(v1.Position, camera, imageSize)
	// p2 := ProjectPoint(v2.Position, camera, imageSize)
	// p3 := ProjectPoint(v3.Position, camera, imageSize)

	p1 := camera.Convert3DTo2D(v1.Position)
	p2 := camera.Convert3DTo2D(v2.Position)
	p3 := camera.Convert3DTo2D(v3.Position)

	if p1 != nil && p2 != nil && p3 != nil {
		k := 1.0 // cull edge thickness
		if ((p1.X < -k || p1.X > k) || (p1.Y < -k || p1.Y > k)) &&
			((p2.X < -k || p2.X > k) || (p2.Y < -k || p2.Y > k)) &&
			((p3.X < -k || p3.X > k) || (p3.Y < -k || p3.Y > k)) {
			return
		}

		halfWidth, halfHeight := imageSize.X/2, imageSize.Y/2
		pi1 := Int_2D{
			int(p1.X*halfWidth + halfWidth),
			int(-p1.Y*halfHeight + halfHeight),
		}
		pi2 := Int_2D{
			int(p2.X*halfWidth + halfWidth),
			int(-p2.Y*halfHeight + halfHeight),
		}
		pi3 := Int_2D{
			int(p3.X*halfWidth + halfWidth),
			int(-p3.Y*halfHeight + halfHeight),
		}

		lightDirection := Normalize(Point3D{-0.3, 0.5, 0.8})
		intensity := DotProduct(normal, lightDirection)
		kShade := 0.7 // 0 = black, 1 = no shade
		intensity = kShade + max(0, min(intensity, 1))*(1-kShade)
		shadedColor := ShadeColor(clr, intensity)

		d1 := calculateDepth(v1.Position, camera)
		d2 := calculateDepth(v2.Position, camera)
		d3 := calculateDepth(v3.Position, camera)

		var w1, w2, w3 float64
		if TexturePerspective {
			w1, w2, w3 = d1, d2, d3
		} else {
			w1, w2, w3 = 1.0, 1.0, 1.0
		}

		ip2 := ImgPoint{pi2.X, pi2.Y, d2, v2.U / w2, v2.V / w2}
		ip1 := ImgPoint{pi1.X, pi1.Y, d1, v1.U / w1, v1.V / w1}
		ip3 := ImgPoint{pi3.X, pi3.Y, d3, v3.U / w3, v3.V / w3}
		DrawTriangle2D2(img, ip1, ip2, ip3, shadedColor, depthBuffer, texture, intensity)
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
		// U: int(interpolate(float64(p1.U), float64(p2.U), f1, f2, fy)),
		// V: int(interpolate(float64(p1.V), float64(p2.V), f1, f2, fy)),
	}
	return p
}

type Cuboid struct {
	vertices [8]Point3D
	Color    color.RGBA
	uvs      [][4][2]float64
}

func MakeCuboidUVsForSingleTexture(texture string, scene *Scene) [][4][2]float64 {
	return CreateCuboidUVs(0, 0, 16.0, 16.0, texture, scene)
}

func MakeCuboidUVs(textures [6]string, scene *Scene) [][4][2]float64 {
	uvs := make([][4][2]float64, 6)
	for i := 0; i < 6; i++ {
		// inefficent generation of uvs for same texture and not all faces needed
		uvsForTexture := CreateCuboidUVs(0, 0, 16.0, 16.0, textures[i], scene)
		uvs[i] = uvsForTexture[i]
	}
	return uvs
}

func MakeAxisAlignedCuboid(min, max Point3D, color color.RGBA, uvs [][4][2]float64) Cuboid {
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

	return Cuboid{vertices, color, uvs}
}

func (c Cuboid) Move(offset Point3D) Cuboid {
	for i := 0; i < 8; i++ {
		c.vertices[i] = c.vertices[i].Add(offset)
	}
	return c
}

type DepthBuffer []float64

func denormaliseCuboidUVs(tilemap *Tilemap, uvs [][4][2]float64) [][4][2]float64 {
	w := float64(tilemap.Image.Bounds().Dx())
	h := float64(tilemap.Image.Bounds().Dy())
	for face := 0; face < len(uvs); face++ {
		for vertex := 0; vertex < 4; vertex++ {
			uvs[face][vertex][0] *= w
			uvs[face][vertex][1] *= h
		}
	}
	return uvs
}

func CreateCuboidUVs(u, v, du, dv float64, texture string, scene *Scene) [][4][2]float64 {
	k := 0.01 // prevent floating point error
	u += k
	v += k
	du -= 2 * k
	dv -= 2 * k

	meta := scene.Tilemap.Metas[texture]
	w := float64(scene.Tilemap.Image.Bounds().Dx())
	h := float64(scene.Tilemap.Image.Bounds().Dy())

	// x0, y0, x1, y1 := meta.U*w, meta.V*h, (meta.U+meta.Width)*w, (meta.V+meta.Height)*h
	x0, y0 := (meta.U + u/w), (meta.V + v/h)
	x1, y1 := x0+(du/w), y0+(dv/h)

	uvs := [][4][2]float64{
		{{x1, y1}, {x1, y0}, {x0, y0}, {x0, y1}}, // Front face
		{{x1, y1}, {x0, y1}, {x0, y0}, {x1, y0}}, // Back face
		{{x0, y0}, {x1, y0}, {x1, y1}, {x0, y1}}, // Top face
		{{x0, y0}, {x1, y0}, {x1, y1}, {x0, y1}}, // Bottom face
		{{x0, y0}, {x1, y0}, {x1, y1}, {x0, y1}}, // Left face
		{{x1, y1}, {x1, y0}, {x0, y0}, {x0, y1}}, // Right face
	}
	return uvs
}

func DrawFilledCuboid(
	cuboid Cuboid,
	camera Camera,
	img *image.RGBA,
	depthBuffer *DepthBuffer,
	tilemap *Tilemap,
	facesToRender *[]int,
) {
	// vertices := [...]Point3D{
	// 	{min.X, min.Y, min.Z}, // 0: Left-bottom-front
	// 	{max.X, min.Y, min.Z}, // 1: Right-bottom-front
	// 	{max.X, max.Y, min.Z}, // 2: Right-top-front
	// 	{min.X, max.Y, min.Z}, // 3: Left-top-front
	// 	{min.X, min.Y, max.Z}, // 4: Left-bottom-back
	// 	{max.X, min.Y, max.Z}, // 5: Right-bottom-back
	// 	{max.X, max.Y, max.Z}, // 6: Right-top-back
	// 	{min.X, max.Y, max.Z}, // 7: Left-top-back
	// }

	// Define the faces of the cube with four vertices each
	faces := [][4]int{
		{0, 3, 2, 1}, // Front face
		{4, 5, 6, 7}, // Back face
		{0, 1, 5, 4}, // Top face *bottom
		{2, 3, 7, 6}, // Bottom face *top
		{7, 3, 0, 4}, // Left face
		{1, 2, 6, 5}, // Right face
	}

	// Define UV coordinates for each face's vertices
	// uvs := [][4][2]float64{
	// 	{{1.0, 1.0}, {1.0, 0.0}, {0.0, 0.0}, {0.0, 1.0}}, // Front face
	// 	{{1.0, 1.0}, {0.0, 1.0}, {0.0, 0.0}, {1.0, 0.0}}, // Back face
	// 	{{0.0, 0.0}, {1.0, 0.0}, {1.0, 1.0}, {0.0, 1.0}}, // Top face
	// 	{{0.0, 0.0}, {1.0, 0.0}, {1.0, 1.0}, {0.0, 1.0}}, // Bottom face
	// 	{{0.0, 0.0}, {1.0, 0.0}, {1.0, 1.0}, {0.0, 1.0}}, // Left face
	// 	{{1.0, 1.0}, {1.0, 0.0}, {0.0, 0.0}, {0.0, 1.0}}, // Right face
	// }
	uvs := denormaliseCuboidUVs(tilemap, cuboid.uvs)
	// Draw each face of the cube using two triangles
	for _, i := range *facesToRender {
		// Split each face into two triangles and draw
		// uvs := createCuboidUVs(tilemap, cuboid.textures[i])
		face := &faces[i]
		uv := &uvs[i]
		DrawTriangle3D(
			Vertex{cuboid.vertices[face[0]], uv[0][0], uv[0][1]},
			Vertex{cuboid.vertices[face[1]], uv[1][0], uv[1][1]},
			Vertex{cuboid.vertices[face[2]], uv[2][0], uv[2][1]},
			camera,
			img, cuboid.Color, depthBuffer,
			&tilemap.Image,
		)
		DrawTriangle3D(
			Vertex{cuboid.vertices[face[0]], uv[0][0], uv[0][1]},
			Vertex{cuboid.vertices[face[2]], uv[2][0], uv[2][1]},
			Vertex{cuboid.vertices[face[3]], uv[3][0], uv[3][1]},
			camera,
			img, cuboid.Color, depthBuffer,
			&tilemap.Image,
		)
	}
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

// scratch

type Vertex2 struct {
	X  int
	Y  int
	Z  float64
	U  float64
	V  float64
	TZ float64
}

func f2i(x float64) int {
	return int(math.Round(x))
}

func getUVColor(u, v, depth float64, texture *image.RGBA) color.RGBA {
	if DebugUV {
		r, g, b := HSVToRGB(float64(int(depth*80.0)%360), 1.0, 1.0)
		return color.RGBA{r, g, b, 255}
		return color.RGBA{100 + uint8(u*255)%50, 100 + uint8(v*255)%50, 100, 255}
		return color.RGBA{uint8(u * 255), uint8(v * 255), 0, 255}
	} else {
		// tx := float64(texture.Bounds().Dx())
		// ty := float64(texture.Bounds().Dy())
		return texture.At(int(u), int(v)).(color.RGBA)
	}
}

func renderFlatBottomTriangle(img *image.RGBA, texture *image.RGBA, v [3]Vertex2, depthBuffer *DepthBuffer, shade float64) {
	//texSize := Point2D{float64(texture.Bounds().Dx()), float64(texture.Bounds().Dy())}
	imageSize := Int_2D{img.Bounds().Dx(), img.Bounds().Dy()}
	// assumes vertices are already ordered such that: v0.Y < v1.Y = v2.Y
	// sort lower vertices so: v1.X < v2.X
	if v[1].X > v[2].X {
		v[1], v[2] = v[2], v[1]
	}

	m01 := float64(v[1].X-v[0].X) / float64(v[1].Y-v[0].Y)
	m02 := float64(v[2].X-v[0].X) / float64(v[2].Y-v[0].Y)

	// w := [3]float64{}
	// for i := 0; i < 3; i++ {
	// 	if !TexturePerspective {
	// 		w[i] = 1
	// 	} else {
	// 		// v[i].U = v[i].U / v[i].Z
	// 		// v[i].V = v[i].V / v[i].Z
	// 		w[i] = 1 / v[i].Z
	// 	}
	// }

	dy := v[1].Y - v[0].Y
	for y := v[0].Y; y <= v[1].Y; y++ {
		if y < 0 {
			continue
		} else if y >= imageSize.Y-1 {
			break
		}
		ddy := float64(y - v[0].Y)
		x01 := f2i(ddy*m01) + v[0].X
		x02 := f2i(ddy*m02) + v[0].X

		ty := ddy / float64(dy)
		u01 := ty*(v[1].U-v[0].U) + v[0].U
		u02 := ty*(v[2].U-v[0].U) + v[0].U

		v01 := ty*(v[1].V-v[0].V) + v[0].V
		v02 := ty*(v[2].V-v[0].V) + v[0].V

		d01 := ty*(v[1].Z-v[0].Z) + v[0].Z
		d02 := ty*(v[2].Z-v[0].Z) + v[0].Z

		w01 := ty*(v[1].TZ-v[0].TZ) + v[0].TZ
		w02 := ty*(v[2].TZ-v[0].TZ) + v[0].TZ

		// w01 := ty*(w[1]-w[0]) + w[0]
		// w02 := ty*(w[2]-w[0]) + w[0]
		// fmt.Println(" ", y, d01, d02)
		for x := x01; x <= x02; x++ {
			if x < 0 {
				continue
			} else if x >= imageSize.X-1 {
				break
			}
			t := float64(x-x01) / float64(x02-x01)
			ww := w01 + t*float64(w02-w01)
			u := u01 + t*float64(u02-u01)
			vv := v01 + t*float64(v02-v01)

			// if !TexturePerspective {
			// 	ww = 1
			// }
			// Interpolate U/Z, V/Z, and 1/Z
			// invDepth1 := 1.0 / d01
			// invDepth2 := 1.0 / d02
			// invDepth := invDepth1 + t*(invDepth2-invDepth1)

			// uOverZ := u01/d01 + t*((u02/d02)-(u01/d01))
			// vOverZ := v01/d01 + t*((v02/d02)-(v01/d01))

			// // Calculate perspective-corrected U and V by dividing by interpolated 1/Z
			// u = u / dept
			// vv = vv / invDepth

			depth := d01 + t*(d02-d01)

			u = u / ww
			vv = vv / ww

			dbi := y*int(imageSize.X) + x
			if dbi > len(*depthBuffer) {
				fmt.Println(dbi, len(*depthBuffer), x, y, imageSize.X)
				panic("Index out of range")
			}
			if depth < (*depthBuffer)[dbi] {
				(*depthBuffer)[dbi] = depth
				clr := getUVColor(u, vv, depth, texture)
				clr = ShadeColor(clr, shade)
				// img.Set(x, y, clr)

				// currentColor := img.RGBAAt(x, y)
				// newColor := CombineColors(clr, currentColor)
				// img.SetRGBA(x, y, clr)
				ii := (imageSize.X*y + x) * 4
				s := img.Pix[ii : ii+3] // Small cap improves performance, see https://golang.org/issue/27857
				s[0] = clr.R
				s[1] = clr.G
				s[2] = clr.B
			}
		}
	}

	// DrawLine(img, Int_2D{v[1].X, v[1].Y}, Int_2D{v[2].X, v[2].Y}, White.ToRGBA())
}

func renderFlatTopTriangle(img *image.RGBA, texture *image.RGBA, v [3]Vertex2, depthBuffer *DepthBuffer, shade float64) {
	// texSize := Point2D{float64(texture.Bounds().Dx()), float64(texture.Bounds().Dy())}
	imageSize := Int_2D{img.Bounds().Dx(), img.Bounds().Dy()}
	// assumes vertices are already ordered such that: v0.Y = v1.Y < v2.Y
	// sort lower vertices so: v0.X < v1.X
	if v[0].X > v[1].X {
		v[0], v[1] = v[1], v[0]
	}

	m02 := float64(v[2].X-v[0].X) / float64(v[2].Y-v[0].Y)
	m12 := float64(v[2].X-v[1].X) / float64(v[2].Y-v[1].Y)

	// w := [3]float64{}
	// for i := 0; i < 3; i++ {
	// 	if !TexturePerspective {
	// 		w[i] = 1
	// 	} else {
	// 		// v[i].U = v[i].U / v[i].Z
	// 		// v[i].V = v[i].V / v[i].Z
	// 		w[i] = 1 / v[i].Z
	// 	}
	// }

	for y := v[0].Y; y <= v[2].Y; y++ {
		if y < 0 {
			continue
		} else if y >= imageSize.Y-1 {
			break
		}
		dy := float64(y - v[2].Y)
		x02 := f2i(dy*m02) + v[2].X
		x12 := f2i(dy*m12) + v[2].X

		ty := float64(y-v[0].Y) / float64(v[2].Y-v[0].Y)

		u02 := ty*float64(v[2].U-v[0].U) + v[0].U
		u12 := ty*float64(v[2].U-v[1].U) + v[1].U

		v02 := ty*float64(v[2].V-v[0].V) + v[0].V
		v12 := ty*float64(v[2].V-v[1].V) + v[1].V

		d02 := ty*(v[2].Z-v[0].Z) + v[0].Z
		d12 := ty*(v[2].Z-v[1].Z) + v[1].Z

		w02 := ty*(v[2].TZ-v[0].TZ) + v[0].TZ
		w12 := ty*(v[2].TZ-v[1].TZ) + v[1].TZ

		// w02 := ty*(w[2]-w[0]) + w[0]
		// w12 := ty*(w[2]-w[1]) + w[1]

		// fmt.Println(y, ty, v02, v12)

		for x := x02; x <= x12; x++ {
			if x < 0 {
				continue
			} else if x >= imageSize.X-1 {
				break
			}

			t := float64(x-x02) / float64(x12-x02)
			u := u02 + t*float64(u12-u02)
			vv := v02 + t*float64(v12-v02)
			ww := w02 + t*float64(w12-w02)
			depth := d02 + t*(d12-d02)

			// if !TexturePerspective {
			// 	ww = 1
			// }

			u = u / ww
			vv = vv / ww

			// invDepth1 := 1.0 / d12
			// invDepth2 := 1.0 / d02
			// invDepth := invDepth1 + t*(invDepth2-invDepth1)

			// uOverZ := u02/d02 + t*((u12/d12)-(u02/d02))
			// vOverZ := v02/d02 + t*((v12/d12)-(v02/d02))

			// Calculate perspective-corrected U and V by dividing by interpolated 1/Z
			// u = u / invDepth
			// vv = u / invDepth

			// clr := texture.At(int(u*texSize.X), int(vv*texSize.Y)).(color.RGBA)
			// clr := color.RGBA{100 + uint8(u*255)%50, 100 + uint8(vv*255)%50, 255, 255}

			dbi := y*int(imageSize.X) + x
			if dbi > len(*depthBuffer) {
				fmt.Println(dbi, len(*depthBuffer), x, y, imageSize.X)
				panic("Index out of range")
			}
			if depth < (*depthBuffer)[dbi] {
				(*depthBuffer)[dbi] = depth
				clr := getUVColor(u, vv, depth, texture)
				clr = ShadeColor(clr, shade)
				// currentColor := img.RGBAAt(x, y)
				// newColor := CombineColors(clr, currentColor)
				// img.SetRGBA(x, y, clr)
				ii := (imageSize.X*y + x) * 4
				s := img.Pix[ii : ii+3] // Small cap improves performance, see https://golang.org/issue/27857
				s[0] = clr.R
				s[1] = clr.G
				s[2] = clr.B
			}
		}
	}
}

func DrawTriangle2D2(img *image.RGBA, p1, p2, p3 ImgPoint, col color.RGBA,
	depthBuffer *DepthBuffer, texture *image.RGBA, shade float64,
) {
	v := [3]Vertex2{
		{p1.X, p1.Y, p1.Z, p1.U, p1.V, 1.0 / p1.Z},
		{p2.X, p2.Y, p2.Z, p2.U, p2.V, 1.0 / p2.Z},
		{p3.X, p3.Y, p3.Z, p3.U, p3.V, 1.0 / p3.Z},
	}
	// if TexturePerspective {
	// 	v = [3]Vertex2{
	// 		{p1.X, p1.Y, p1.Z, p1.U / p1.Z, p1.V / p1.Z, 1.0 / p1.Z},
	// 		{p2.X, p2.Y, p2.Z, p2.U / p2.Z, p2.V / p2.Z, 1.0 / p2.Z},
	// 		{p3.X, p3.Y, p3.Z, p3.U / p3.Z, p3.V / p3.Z, 1.0 / p3.Z},
	// 	}
	// } else {
	// 	v = [3]Vertex2{
	// 		{p1.X, p1.Y, p1.Z, p1.U, p1.V, 1.0},
	// 		{p2.X, p2.Y, p2.Z, p2.U, p2.V, 1.0},
	// 		{p3.X, p3.Y, p3.Z, p3.U, p3.V, 1.0},
	// 	}
	// }
	renderTriangle(img, texture, v, depthBuffer, shade)
	if drawTriangleWire {
		DrawLine(img, Int_2D{p1.X, p1.Y}, Int_2D{p2.X, p2.Y}, White.ToRGBA())
		DrawLine(img, Int_2D{p1.X, p1.Y}, Int_2D{p3.X, p3.Y}, White.ToRGBA())
		DrawLine(img, Int_2D{p3.X, p3.Y}, Int_2D{p2.X, p2.Y}, White.ToRGBA())
	}

}

func renderTriangle(img *image.RGBA, texture *image.RGBA, v [3]Vertex2, depthBuffer *DepthBuffer, shade float64) {
	// sort vertices so v0.Y <= v1.Y <= v2.Y
	if v[0].Y > v[1].Y {
		v[0], v[1] = v[1], v[0]
	}
	if v[1].Y > v[2].Y {
		v[1], v[2] = v[2], v[1]
	}
	if v[0].Y > v[1].Y {
		v[0], v[1] = v[1], v[0]
	}

	if v[1].Y == v[2].Y {
		renderFlatBottomTriangle(img, texture, v, depthBuffer, shade)
	} else if v[0].Y == v[1].Y {
		renderFlatTopTriangle(img, texture, v, depthBuffer, shade)
	} else {
		t := float64(v[1].Y-v[0].Y) / float64(v[2].Y-v[0].Y)
		x := v[0].X + int(t*float64(v[2].X-v[0].X))
		z := v[0].Z + t*(v[2].Z-v[0].Z)
		a := v[0].U + t*(v[2].U-v[0].U)
		b := v[0].V + t*(v[2].V-v[0].V)
		tz := v[0].TZ + t*(v[2].TZ-v[0].TZ)
		// fmt.Println(z, v[0].Z, v[2].Z)
		renderFlatBottomTriangle(img, texture, [3]Vertex2{v[0], {x, v[1].Y, z, a, b, tz}, v[1]}, depthBuffer, shade)
		renderFlatTopTriangle(img, texture, [3]Vertex2{{x, v[1].Y, z, a, b, tz}, v[1], v[2]}, depthBuffer, shade)
	}
}

// assumes world of 16x16x16
func Convert1DTo3D(i int) Point3D {
	// Extract x, y, z using bit-shifting
	x := i & 0xF        // Mask the lowest 4 bits
	y := (i >> 4) & 0xF // Shift right by 4 bits and mask the next 4 bits
	z := (i >> 8) & 0xF // Shift right by 8 bits and mask the next 4 bits

	return Point3D{float64(x), float64(y), float64(z)}
}

func DrawObjects(scene *Scene, img *image.RGBA, depthBuffer *DepthBuffer) {
	// {0, 3, 2, 1}, // Front face
	// {4, 5, 6, 7}, // Back face
	// {0, 1, 5, 4}, // Top face
	// {2, 3, 7, 6}, // Bottom face
	// {7, 3, 0, 4}, // Left face
	// {1, 2, 6, 5}, // Right face
	var Directions = [6]Direction{Back, Front, Down, Up, Left, Right} // inconsistent direction order
	for i, block := range scene.World.Blocks {
		rb, isRenderable := block.(WireRenderBlock)

		if isRenderable {
			// if a block and is neighbour are opaque on their shared face then dont render
			var faces []int
			position := Convert1DTo3D(i)
			opaqueBlock, isOpaqueBlock := block.(OpaqueBlock)
			if skipAdjacentFaces && isOpaqueBlock {
				for face, direction := range Directions {
					if !opaqueBlock.IsOpaqueInDirection(direction) {
						continue
					}
					np := position.ToVec3().Move(direction)
					neighbour := scene.World.GetBlock(np)
					nOpaqueBlock, nIsOpaqueBlock := neighbour.(OpaqueBlock)
					if !nIsOpaqueBlock || !nOpaqueBlock.IsOpaqueInDirection(direction.GetOppositeDirection()) {
						faces = append(faces, face)
					}
				}
			} else {
				faces = []int{0, 1, 2, 3, 4, 5}
			}

			for _, c := range rb.ToCuboids(scene) {
				// generating a new cuboid is bad. mutate or move vertices dynamically inside function
				movedCuboid := c.Move(position)
				DrawFilledCuboid(movedCuboid, scene.Camera, img, depthBuffer, &scene.Tilemap, &faces)

			}
		}
	}

	// Calculate aspect ratio based on the image dimensions

	// blockSize := 1

	// gridSize := 16
	// fHalfGridSize := float64(gridSize) / 2.0

	// Constants for alpha scaling
	// const maxAlpha = 1.0
	// const minAlpha = 0.0
	// const maxDistance = 16.0 // Distance at which alpha should be fully transparent
	// const minDistance = 8.0  // Distance at which alpha is fully opaque
	// fmt.Println("Drawing", scene.Iteration)
	// for _, bd := range GetSortedBlockPositions(scene.Camera.Position) {
	// 	p := Vec3{int(bd.Position.X), int(bd.Position.Y), int(bd.Position.Z)}
	// 	block := scene.World.GetBlock(p)
	// 	rb, isRenderable := block.(WireRenderBlock)
	// 	localOffset := bd.Position
	// 	// cameraHorPos := scene.Camera.Position
	// 	// cameraHorPos.Y = localOffset.Y
	// 	// distance := Distance(localOffset, cameraHorPos)
	// 	var alphaScaling float64 = 1.0

	// 	// if distance >= maxDistance {
	// 	// 	alphaScaling = minAlpha
	// 	// } else if distance <= minDistance {
	// 	// 	alphaScaling = maxAlpha
	// 	// } else {
	// 	// 	alphaScaling = maxAlpha * (1 - (distance / maxDistance))
	// 	// }

	// 	if isRenderable {
	// 		for _, c := range rb.ToCuboids(scene) {

	// 			mc := c.Move(localOffset)

	// 			mc.Color.A = uint8(float64(mc.Color.A) * alphaScaling)
	// 			// fmt.Println(mc, globalOffset, localOffset)
	// 			DrawFilledCuboid(mc, scene.Camera, img, depthBuffer, &scene.Tilemap)
	// 			// DrawCuboid(mc, scene.Camera, img)
	// 			// fmt.Println(mc)
	// 		}
	// 	} else if p.Y == 0 {
	// 		// minP := Point3D{
	// 		// 	X: float64(p.X * blockSize),
	// 		// 	Y: float64(p.Y * blockSize),
	// 		// 	Z: float64(p.Z * blockSize),
	// 		// }
	// 		// maxP := Point3D{
	// 		// 	X: float64((p.X + 1) * blockSize),
	// 		// 	Y: float64((p.Y + 1) * blockSize),
	// 		// 	Z: float64((p.Z + 1) * blockSize),
	// 		// }
	// 		// c := MakeAxisAlignedCuboid(minP, maxP, color.RGBA{255, 255, 255, 100}, [6]string{"test", "test", "test", "test", "test", "test"})
	// 		// // c.Color.A = uint8(float64(c.Color.A) * alphaScaling)
	// 		// DrawFilledCuboid(c, scene.Camera, img, depthBuffer, &scene.Tilemap)
	// 	}
	// }
}

func DrawDebugInformation(scene *Scene, img *image.RGBA) {
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

	DrawText(img, 4, fontSize*3, fmt.Sprintf(
		"XYZ: %.1f %.1f %.1f, Rot: %.1f %.1f %.1f",
		scene.Camera.Position.X,
		scene.Camera.Position.Y,
		scene.Camera.Position.Z,
		RadToDeg(scene.Camera.Rotation.X),
		RadToDeg(scene.Camera.Rotation.Y),
		RadToDeg(scene.Camera.Rotation.Z),
	), Cyan.ToRGBA(), scene.FontFace)
}

func DrawTestTriangles(scene *Scene, img *image.RGBA, depthBuffer *DepthBuffer) {
	k := 1.0
	//front
	DrawTriangle3D(
		Vertex{Point3D{0, 0, 0}, 0, 0},
		Vertex{Point3D{k, 0, 0}, 16, 0},
		Vertex{Point3D{k, k, 0}, 16, 16},
		scene.Camera, img, Red.ToRGBA(), depthBuffer, &scene.Tilemap.Image,
	)

	// DrawTriangle3D(
	// 	Vertex{Point3D{0, 0, 0}, 0, 0},
	// 	Vertex{Point3D{k, k, 0}, 16, 16},
	// 	Vertex{Point3D{0, k, 0}, 0, 16},
	// 	scene.Camera, img, Orange.ToRGBA(), &depthBuffer, &scene.Tilemap.Image,
	// )

	// // left
	// DrawTriangle3D(Point3D{0, 0, 0}, Point3D{0, 0, k}, Point3D{0, k, k}, scene.Camera, img, Yellow.ToRGBA(), &depthBuffer)
	// DrawTriangle3D(Point3D{0, 0, 0}, Point3D{0, k, k}, Point3D{0, k, 0}, scene.Camera, img, Green.ToRGBA(), &depthBuffer)

	// // back
	// DrawTriangle3D(Point3D{0, 0, k}, Point3D{k, k, k}, Point3D{k, 0, k}, scene.Camera, img, Lime.ToRGBA(), &depthBuffer)

	// DrawTriangle3D(Point3D{0, 0, k}, Point3D{0, k, k}, Point3D{k, k, k}, scene.Camera, img, LightBlue.ToRGBA(), &depthBuffer)

	// // right
	// DrawTriangle3D(Point3D{k, 0, 0}, Point3D{k, 0, k}, Point3D{k, k, k}, scene.Camera, img, Cyan.ToRGBA(), &depthBuffer)
	// DrawTriangle3D(Point3D{k, 0, 0}, Point3D{k, k, k}, Point3D{k, k, 0}, scene.Camera, img, Blue.ToRGBA(), &depthBuffer)
}

func GetSign(x float64) int {
	if x < 0 {
		return -1
	}
	return 1
}

func Distance2(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return dx*dx + dy*dy
}

func SqDist3(x1, y1, z1, x2, y2, z2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	dz := x2 - x1
	return dx*dx + dy*dy + dz*dz
}

// functional 2D line-grid intersection algorithm (z const)
func DDA3D2(x1, y1, z1, x2, y2, z2 float64) []Point3D {
	// Store the line points
	var points []Point3D
	points = append(points, Point3D{x1, y1, z1})

	dx, dy := x2-x1, y2-y1
	sx, sy := GetSign(dx), GetSign(dy)
	dydx, dxdy := dy/dx, dx/dy
	var x, y float64
	if sx == -1 {
		x = math.Floor(x1)
	} else {
		x = math.Ceil(x1)
	}
	if sy == -1 {
		y = math.Floor(y1)
	} else {
		y = math.Ceil(y1)
	}

	for ((sx == 1 && x < x2) || (sx == -1 && x > x2)) ||
		((sy == 1 && y < y2) || (sy == -1 && y > y2)) {
		yn, xn := dydx*(x-x1)+y1, dxdy*(y-y1)+x1
		ddx, ddy := Distance2(x1, y1, xn, y), Distance2(x1, y1, x, yn)
		var p Point3D
		if ddx > ddy {
			p = Point3D{x, yn, z1}
			x += float64(sx)
		} else {
			p = Point3D{xn, y, z1}
			y += float64(sy)
		}
		points = append(points, p)
	}
	points = append(points, Point3D{x2, y2, z2})

	return points
}

func boolRound(x float64, doFloor bool) float64 {
	if doFloor {
		return math.Floor(x)
	} else {
		return math.Ceil(x)
	}
}

func replaceNaNWithInf(value float64) float64 {
	if math.IsNaN(value) {
		return math.Inf(1) // Replace NaN with positive infinity
	}
	return value
}

// functional 3D line-grid intersection algorithm
func DDA3D(x1, y1, z1, x2, y2, z2 float64) []Point3D {
	// Store the line points
	var points []Point3D
	points = append(points, Point3D{x1, y1, z1})

	// Calculate differences
	dx, dy, dz := x2-x1, y2-y1, z2-z1
	sx, sy, sz := GetSign(dx), GetSign(dy), GetSign(dz)

	x, y, z := boolRound(x1, sx == -1), boolRound(y1, sy == -1), boolRound(z1, sz == -1)

	for ((sx == 1 && x < x2) || (sx == -1 && x > x2)) ||
		((sy == 1 && y < y2) || (sy == -1 && y > y2)) ||
		((sz == 1 && z < z2) || (sz == -1 && z > z2)) {

		tx, ty, tz := (x-x1)/dx, (y-y1)/dy, (z-z1)/dz

		ytx, ztx := (y2-y1)*tx+y1, (z2-z1)*tx+z1
		xty, zty := (x2-x1)*ty+x1, (z2-z1)*ty+z1
		xtz, ytz := (x2-x1)*tz+x1, (y2-y1)*tz+y1

		// fmt.Println("p1", x1, y1, z1)

		// fmt.Println("d", y, dy, ty)
		// // fmt.Println("x", x, ytx, ztx)
		// fmt.Println("y", xty, y, zty)
		// // fmt.Println("z", xtz, ytz, z)
		// fmt.Println("p2", x2, y2, z2)

		ddx := replaceNaNWithInf(SqDist3(x1, y1, z1, x, ytx, ztx))
		ddy := replaceNaNWithInf(SqDist3(x1, y1, z1, xty, y, zty))
		ddz := replaceNaNWithInf(SqDist3(x1, y1, z1, xtz, ytz, z))

		// fmt.Println("dd", ddx, ddy, ddz)

		var p Point3D
		if ddx < ddy && ddx < ddz {
			p = Point3D{x, ytx, ztx}
			x += float64(sx)
		} else if ddy < ddz {
			p = Point3D{xty, y, zty}
			y += float64(sy)
		} else {
			p = Point3D{xtz, ytz, z}
			z += float64(sz)
		}
		points = append(points, p)
	}

	points = append(points, Point3D{x2, y2, z2})

	return points
}

func clearSceneImage(img *image.RGBA) {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	color := color.RGBA{122, 168, 253, 255}
	pxCount := width * height
	// Directly modify the pixel array
	for i := 0; i < pxCount; i++ {
		// Calculate the pixel's starting index in the Pix slice (4 bytes per pixel)
		offset := i * 4
		// Set the pixel color (RGBA)
		img.Pix[offset] = color.R   // Red
		img.Pix[offset+1] = color.G // Green
		img.Pix[offset+2] = color.B // Blue
		img.Pix[offset+3] = color.A // Alpha
	}
}

func clearDepthBuffer(depthBuffer *DepthBuffer) {
	for i := range *depthBuffer {
		(*depthBuffer)[i] = 1e9 // A large value representing 'infinity'
	}
}

func DebugDrawRayCast(scene *Scene, img *image.RGBA, startPos Point3D, rotation Point3D) {
	// endPos := CalculateEndPosition(scene.Camera.Position, scene.Camera.Rotation, 5)
	// startPos := Point3D{0.1, 9.7, -2.7}
	endPos := CalculateEndPosition(startPos, rotation, 8)

	x1, y1, z1 := startPos.ToComponents()
	x2, y2, z2 := endPos.ToComponents()
	points := DDA3D(x1, y1, z1, x2, y2, z2)
	// fmt.Println(points)
	k := 0.1
	clr := color.RGBA{0, 255, 0, 255}
	clr2 := color.RGBA{255, 0, 0, 255}
	clr3 := color.RGBA{0, 0, 255, 255}
	uvs := MakeCuboidUVsForSingleTexture("redstone_lamp_on", scene)
	for _, p := range points {
		p1 := Point3D{math.Floor(p.X), math.Floor(p.Y), math.Floor(p.Z)}
		p2 := p1.Add(Point3D{1, 1, 1})

		p3 := Point3D{p.X - k/2, p.Y - k/2, p.Z - k/2}
		p4 := p3.Add(Point3D{k, k, k})
		// fmt.Println(i, p1, p2)
		cuboid1 := MakeAxisAlignedCuboid(p1, p2, clr, uvs)
		DrawCuboid(cuboid1, scene.Camera, img)

		cuboid2 := MakeAxisAlignedCuboid(p3, p4, clr2, uvs)
		DrawCuboid(cuboid2, scene.Camera, img)
	}

	// attempt to work out face of ray cast
	for i := 2; i < len(points); i++ {
		p0 := points[i-1]
		p1 := points[i]

		fp1 := Point3D{math.Floor(p1.X), math.Floor(p1.Y), math.Floor(p1.Z)}

		dp0 := p0.Subtract(fp1)

		var dir Direction
		if dp0.X == 0 || dp0.X == -1 {
			dir = Left
		} else if dp0.X == 1 {
			dir = Right
		} else if dp0.Y == 0 || dp0.Y == -1 {
			dir = Down
		} else if dp0.Y == 1 {
			dir = Up
		} else if dp0.Z == 0 || dp0.Z == -1 {
			dir = Front
		} else if dp0.Z == 1 {
			dir = Back
		} else {
			fmt.Println(i, p0, dp0)
			panic("position not along any face")
		}
		dp := dir.ToVec3().ToPoint3D().Scale(0.1)
		p2 := p0.Add(dp)
		// fmt.Println(i, p0, dp0, dir)
		cuboid2 := MakeAxisAlignedCuboid(p0, p2, clr3, uvs)
		DrawCuboid(cuboid2, scene.Camera, img)
	}
}

func GetRayCastPositions(scene *Scene) (previous *Vec3, selected *Vec3) {
	// selected will return nil if all air blocks
	// previous will return nil is selected nil, or selected is current position
	startPos := scene.Camera.Position
	r := scene.Camera.Rotation
	endPos := CalculateEndPosition(startPos, r.Scale(-1), 8)

	x1, y1, z1 := startPos.ToComponents()
	x2, y2, z2 := endPos.ToComponents()
	points := DDA3D(x1, y1, z1, x2, y2, z2)

	midPoints := make([]Point3D, len(points)-1)

	for i := 0; i < len(points)-1; i++ {
		p1, p2 := points[i], points[i+1]
		mp := p1.Add(p2).Scale(0.5)
		midPoints[i] = mp
	}

	var fpoints []Vec3

	// get the unique block positions
	for _, p := range midPoints {
		fp := p.Floor().ToVec3()
		n := len(fpoints)
		// fmt.Println(fpoints, len(fpoints))
		if n == 0 {
			fpoints = append(fpoints, fp)
		} else {
			// fmt.Println(fpoints, len(fpoints), i, i-1)
			if fpoints[n-1].Equals(fp) {
				continue
			}
			fpoints = append(fpoints, fp)
		}
	}

	for i, p := range fpoints {
		b := scene.World.GetBlock(p)
		_, isAir := b.(Air)
		if !isAir {
			if i == 0 {
				return nil, &p
			}
			return &fpoints[i-1], &p
		}
	}

	return nil, nil
}

func DrawGetRayCastPositions(scene *Scene, img *image.RGBA, startPos Point3D, rotation Point3D) {
	// selected will return nil if all air blocks
	// previous will return nil is selected nil, or selected is current position
	endPos := CalculateEndPosition(startPos, rotation, 8)

	x1, y1, z1 := startPos.ToComponents()
	x2, y2, z2 := endPos.ToComponents()
	points := DDA3D(x1, y1, z1, x2, y2, z2)

	midPoints := make([]Point3D, len(points)-1)

	for i := 0; i < len(points)-1; i++ {
		p1, p2 := points[i], points[i+1]
		mp := p1.Add(p2).Scale(0.5)
		midPoints[i] = mp
	}

	// var fpoints []Vec3

	// // get the unique block positions
	// for _, p := range midPoints {
	// 	fp := p.ToVec3()
	// 	n := len(fpoints)
	// 	// fmt.Println(fpoints, len(fpoints))
	// 	if n == 0 {
	// 		fpoints = append(fpoints, fp)
	// 	} else {
	// 		// fmt.Println(fpoints, len(fpoints), i, i-1)
	// 		if fpoints[n-1].Equals(fp) {
	// 			continue
	// 		}
	// 		fpoints = append(fpoints, fp)
	// 	}
	// }
	uvs := MakeCuboidUVsForSingleTexture("redstone_lamp_on", scene)
	clr := color.RGBA{0, 0, 255, 255}
	clr2 := color.RGBA{100, 100, 255, 255}
	for _, p := range midPoints {
		dd := Point3D{0.1, 0.1, 0.1}
		p1 := p.Subtract(dd.Scale(0.5))
		p2 := p1.Add(dd)
		cuboid1 := MakeAxisAlignedCuboid(p1, p2, clr, uvs)
		DrawCuboid(cuboid1, scene.Camera, img)

		p1 = p.Floor().Add(Point3D{0.5, 0.5, 0.5}).Subtract(dd.Scale(0.5))
		p2 = p1.Add(dd)
		cuboid1 = MakeAxisAlignedCuboid(p1, p2, clr2, uvs)
		DrawCuboid(cuboid1, scene.Camera, img)
	}
}

func DrawWireCube(scene *Scene, img *image.RGBA, p Point3D) {
	clr := color.RGBA{0, 255, 0, 255}
	uvs := MakeCuboidUVsForSingleTexture("redstone_lamp_on", scene)

	p2 := p.Add(Point3D{1, 1, 1})
	cuboid := MakeAxisAlignedCuboid(p, p2, clr, uvs)
	DrawCuboid(cuboid, scene.Camera, img)
}

func DrawRayCast(scene *Scene, img *image.RGBA) {
	startPos := scene.Camera.Position
	r := scene.Camera.Rotation
	endPos := CalculateEndPosition(startPos, r.Scale(-1), 8)
	// startPos := Point3D{5.2, 5.3, -5.4}
	// endPos := CalculateEndPosition(startPos, Point3D{0, DegToRad(22.5), 0}, 5)

	x1, y1, z1 := startPos.ToComponents()
	x2, y2, z2 := endPos.ToComponents()
	points := DDA3D(x1, y1, z1, x2, y2, z2)

	clr := color.RGBA{255, 255, 255, 255}
	uvs := MakeCuboidUVsForSingleTexture("redstone_lamp_on", scene)

	k := 0.1
	m := 0.5 - k/2
	for _, p := range points {
		p1 := Point3D{math.Floor(p.X) + m, math.Floor(p.Y) + m, math.Floor(p.Z) + m}
		// b := scene.World.GetBlock(p1.ToVec3())
		// _, isAir := b.(Air)
		// if !isAir {
		p2 := p1.Add(Point3D{k, k, k})
		cuboid := MakeAxisAlignedCuboid(p1, p2, clr, uvs)
		DrawCuboid(cuboid, scene.Camera, img)
		// }
	}

}

func drawCrossHair(img *image.RGBA) {
	x, y := img.Bounds().Dx()/2, img.Bounds().Dy()/2
	k := 5
	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}

	callback := func(x, y int) {
		inColor := img.At(x, y)
		r, g, b, _ := inColor.RGBA()

		r, g, b = r>>8, g>>8, b>>8

		avgChannel := (r + g + b) / 3
		var clr color.RGBA
		if avgChannel > 125 {
			clr = black
		} else {
			clr = white
		}
		img.SetRGBA(x, y, clr)
	}

	GenerateLine(Int_2D{x + k, y}, Int_2D{x + 1, y}, callback)
	GenerateLine(Int_2D{x - k, y}, Int_2D{x - 1, y}, callback)
	GenerateLine(Int_2D{x, y + k}, Int_2D{x, y - k}, callback)
	// DrawLine(img, Int_2D{x + k, y}, Int_2D{x - k, y}, clr)
	// DrawLine(img, Int_2D{x, y + k}, Int_2D{x, y - k}, clr)
}

func DrawScene(scene *Scene, img *image.RGBA, depthBuffer *DepthBuffer) {
	clearDepthBuffer(depthBuffer)
	clearSceneImage(img)

	// DrawTestTriangles(scene, img, depthBuffer)

	DrawObjects(scene, img, depthBuffer)

	// startPos := Point3D{0.1, 9.7, -2.7}
	// rotation := Point3D{DegToRad(332.5), DegToRad(349.5), 0}.Scale(-1)

	// startPos := scene.Player.Position
	// rotation := scene.Player.Rotation.Scale(-1)

	// DrawRayCast(scene, img)
	// DrawGetRayCastPositions(scene, img, startPos, rotation)
	// DebugDrawRayCast(scene, img, startPos, rotation)

	_, selectedPos := GetRayCastPositions(scene)
	if selectedPos != nil {
		DrawWireCube(scene, img, selectedPos.ToPoint3D())
	}

	drawCrossHair(img)
	DrawDebugInformation(scene, img)
}
