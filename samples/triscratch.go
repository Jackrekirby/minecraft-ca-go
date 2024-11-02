package samples

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func drawLine(x1, y1, x2, y2 int, img *image.RGBA, col color.Color) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := -1
	sy := -1
	if x1 < x2 {
		sx = 1
	}
	if y1 < y2 {
		sy = 1
	}
	err := dx - dy

	for {
		img.Set(x1, y1, col)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Vertex2 struct {
	X int
	Y int
	U int
	V int
}

func f2i(x float64) int {
	return int(math.Round(x))
}

func renderFlatBottomTriangle(img *image.RGBA, texture *image.RGBA, v [3]Vertex2) {
	// assumes vertices are already ordered such that: v0.Y < v1.Y = v2.Y
	// sort lower vertices so: v1.X < v2.X
	if v[1].X > v[2].X {
		v[1], v[2] = v[2], v[1]
	}

	m01 := float64(v[1].X-v[0].X) / float64(v[1].Y-v[0].Y)
	m02 := float64(v[2].X-v[0].X) / float64(v[2].Y-v[0].Y)

	dy := v[1].Y - v[0].Y
	for y := 0; y <= dy; y++ {
		x01 := f2i(float64(y)*m01) + v[0].X
		x02 := f2i(float64(y)*m02) + v[0].X

		ty := float64(y) / float64(dy)
		u01 := f2i(ty*float64(v[1].U-v[0].U)) + v[0].U
		u02 := f2i(ty*float64(v[2].U-v[0].U)) + v[0].U

		v01 := f2i(ty*float64(v[1].V-v[0].V)) + v[0].V
		v02 := f2i(ty*float64(v[2].V-v[0].V)) + v[0].V
		for x := x01; x <= x02; x++ {
			t := float64(x-x01) / float64(x02-x01)
			u := u01 + f2i(t*float64(u02-u01))
			vv := v01 + f2i(t*float64(v02-v01))
			// v := v01 + int(t*float64(v02-v01))
			// var r, g, b uint8
			// if (u/10)%2 == 0 {
			// 	r = 255
			// } else {
			// 	r = 0
			// }
			// if (vv/10)%2 == 0 {
			// 	g = 255
			// } else {
			// 	g = 0
			// }
			// b = 0
			// img.Set(x, v[0].Y+y, color.RGBA{r, g, b, 255})
			texColor := texture.At(u, vv).(color.RGBA)
			// clr := color.RGBA{uint8(u), uint8(vv), 0, 255}
			img.Set(x, y, texColor)
		}
	}
}

func renderFlatTopTriangle(img *image.RGBA, texture *image.RGBA, v [3]Vertex2) {
	// assumes vertices are already ordered such that: v0.Y = v1.Y < v2.Y
	// sort lower vertices so: v0.X < v1.X
	if v[0].X > v[1].X {
		v[0], v[1] = v[1], v[0]
	}

	m02 := float64(v[2].X-v[0].X) / float64(v[2].Y-v[0].Y)
	m12 := float64(v[2].X-v[1].X) / float64(v[2].Y-v[1].Y)

	for y := v[0].Y; y <= v[2].Y; y++ {
		dy := float64(y - v[2].Y)
		x02 := f2i(dy*m02) + v[2].X
		x12 := f2i(dy*m12) + v[2].X

		ty := float64(y-v[0].Y) / float64(v[2].Y-v[0].Y)

		u02 := f2i(ty*float64(v[2].U-v[0].U)) + v[0].U
		u12 := f2i(ty*float64(v[2].U-v[1].U)) + v[1].U

		v02 := f2i(ty*float64(v[2].V-v[0].V)) + v[0].V
		v12 := f2i(ty*float64(v[2].V-v[1].V)) + v[1].V

		// fmt.Println(y, ty, v02, v12)

		for x := x02; x <= x12; x++ {
			t := float64(x-x02) / float64(x12-x02)
			u := u02 + f2i(t*float64(u12-u02))
			vv := v02 + f2i(t*float64(v12-v02))
			// fmt.Println(" > ", x, t, u, vv)
			// v := v01 + int(t*float64(v02-v01))
			// var r, g, b uint8
			// if (u/10)%2 == 0 {
			// 	r = 255
			// } else {
			// 	r = 0
			// }
			// if (vv/10)%2 == 0 {
			// 	g = 255
			// } else {
			// 	g = 0
			// }
			// b = 0
			texColor := texture.At(u, vv).(color.RGBA)
			// clr := color.RGBA{uint8(u), uint8(vv), 0, 255}
			img.Set(x, y, texColor)
		}
	}

	// m02 := float64(v[2].X-v[0].X) / float64(v[2].Y-v[0].Y)
	// m12 := float64(v[2].X-v[1].X) / float64(v[2].Y-v[1].Y)

	// dy := v[2].Y - v[0].Y
	// for y := 0; y <= dy; y++ {
	// 	x02 := int(float64(y)*m02) + v[0].X
	// 	x12 := int(float64(y)*m12) + v[1].X
	// 	for x := x02; x <= x12; x++ {
	// 		img.Set(x, v[0].Y+y, color.RGBA{0, 255, 0, 255})
	// 	}
	// }
}

func renderTriangle(img *image.RGBA, texture *image.RGBA, v [3]Vertex2) {
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
		renderFlatBottomTriangle(img, texture, v)
	} else if v[0].Y == v[1].Y {
		renderFlatTopTriangle(img, texture, v)
	} else {
		t := float64(v[1].Y-v[0].Y) / float64(v[2].Y-v[0].Y)
		x := v[0].X + int(t*float64(v[2].X-v[0].X))
		a := v[0].U + int(t*float64(v[2].U-v[0].U))
		b := v[0].V + int(t*float64(v[2].V-v[0].V))
		renderFlatBottomTriangle(img, texture, [3]Vertex2{v[0], {x, v[1].Y, a, b}, v[1]})
		renderFlatTopTriangle(img, texture, [3]Vertex2{{x, v[1].Y, a, b}, v[1], v[2]})
	}
}

func drawLine3(img *image.RGBA, x0, y0, x1, y1 int) {
	dx := x1 - x0
	dy := y1 - y0

	var sx, sy int

	if dx > 0 {
		sx = 1
	} else {
		sx = -1
	}

	if dy > 0 {
		sy = 1
	} else {
		sy = -1
	}

	x := x0
	y := y0
	for {
		img.Set(x, y, color.RGBA{255, 0, 0, 255})
		tx := float64(x-x0) / float64(dx)
		ty := float64(y-y0) / float64(dy)
		if tx > ty {
			y += sy
		} else {
			x += sx
		}
		if x == x1 && y == y1 {
			break
		}
	}
}

// func drawLine2(img *image.RGBA, v [2]Vertex2) {
// 	if v[1].X != v[0].X {
// 		if v[0].X > v[1].X {
// 			v[0], v[1] = v[1], v[0]
// 		}
// 		m := (v[1].Y - v[0].Y) / (v[1].X - v[0].X)
// 		fmt.Println("X", v, m)
// 		for x := v[0].X; x <= v[1].X; x++ {
// 			y := m*x + v[0].Y
// 			// fmt.Println(" ", x, y)
// 			img.Set(int(math.Round(x)), int(math.Round(y)), color.RGBA{0, 0, 255, 255})
// 		}
// 	} else if v[1].Y != v[0].Y {
// 		if v[0].Y > v[1].Y {
// 			v[0], v[1] = v[1], v[0]
// 		}
// 		m := (v[1].X - v[0].X) / (v[1].Y - v[0].Y)
// 		fmt.Println("Y", v, m)
// 		for y := v[0].Y; y <= v[1].Y; y++ {
// 			x := m*y + v[0].X
// 			// fmt.Println(" ", x, y)
// 			img.Set(int(math.Round(x)), int(math.Round(y)), color.RGBA{0, 0, 255, 255})
// 		}
// 	} else {
// 		// fmt.Println("P", v)
// 		img.Set(int(math.Round(v[0].X)), int(math.Round(v[0].Y)), color.RGBA{0, 0, 255, 255})
// 	}
// }

func drawWireTriangle(img *image.RGBA, v [3]Vertex2) {
	drawLine(v[0].X, v[0].Y, v[1].X, v[1].Y, img, color.RGBA{0, 0, 255, 255})
	drawLine(v[1].X, v[1].Y, v[2].X, v[2].Y, img, color.RGBA{0, 0, 255, 255})
	drawLine(v[2].X, v[2].Y, v[0].X, v[0].Y, img, color.RGBA{0, 0, 255, 255})
}

func saveImage(img *image.RGBA, filename string) {
	file, err := os.Create(fmt.Sprintf("%s.png", filename))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}

func loadTexture() *image.RGBA {
	// Load the texture
	textureFile, err := os.Open("assets/crafting_table_front.png")
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

	return texture
}

func Main2() {
	fmt.Println("Learning")

	width, height := 128, 128
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{122, 168, 253, 255})
		}
	}

	// randVertex := func() Vertex2 {
	// 	return Vertex2{
	// 		int(rand.Float64() * float64(width-1)),
	// 		int(rand.Float64() * float64(height-1)),
	// 	}
	// }
	// v := [3]Vertex2{
	// 	randVertex(),
	// 	randVertex(),
	// 	randVertex(),
	// }

	// v := [3]Vertex2{
	// 	{50, 0, 0, 0},
	// 	{0, 80, 0, 255},
	// 	{100, 100, 255, 255},
	// }

	// drawLine(24, 8, 16, 32, img, color.RGBA{255, 0, 0, 255})

	texture := loadTexture()

	v := [3]Vertex2{
		{20, 0, 0, 0},
		{100, 50, 16, 0},
		{80, 100, 16, 16},
	}
	renderTriangle(img, texture, v)
	drawWireTriangle(img, v)

	v2 := [3]Vertex2{
		{20, 0, 0, 0},
		{0, 50, 0, 16},
		{80, 100, 16, 16},
	}
	renderTriangle(img, texture, v2)
	drawWireTriangle(img, v2)

	saveImage(img, "output/scratch")
}
