package core

import (
	"image/color"
	"math"
)

type Color int

const (
	White Color = iota
	Orange
	Magenta
	LightBlue
	Yellow
	Lime
	Pink
	Gray
	LightGray
	Cyan
	Purple
	Blue
	Brown
	Green
	Red
	Black
)

// Colors stores the RGBA color values for each color.
var Colors = map[Color]color.RGBA{
	White:     {R: 255, G: 255, B: 255, A: 255},
	Orange:    {R: 216, G: 127, B: 51, A: 255},
	Magenta:   {R: 178, G: 76, B: 216, A: 255},
	LightBlue: {R: 102, G: 153, B: 216, A: 255},
	Yellow:    {R: 229, G: 229, B: 51, A: 255},
	Lime:      {R: 127, G: 204, B: 25, A: 255},
	Pink:      {R: 242, G: 127, B: 165, A: 255},
	Gray:      {R: 76, G: 76, B: 76, A: 255},
	LightGray: {R: 153, G: 153, B: 153, A: 255},
	Cyan:      {R: 76, G: 127, B: 153, A: 255},
	Purple:    {R: 127, G: 63, B: 178, A: 255},
	Blue:      {R: 51, G: 76, B: 178, A: 255},
	Brown:     {R: 102, G: 76, B: 51, A: 255},
	Green:     {R: 102, G: 127, B: 51, A: 255},
	Red:       {R: 153, G: 51, B: 51, A: 255},
	Black:     {R: 25, G: 25, B: 25, A: 255},
}

// String returns the string representation of a Color.
func (w Color) String() string {
	switch w {
	case White:
		return "White"
	case Orange:
		return "Orange"
	case Magenta:
		return "Magenta"
	case LightBlue:
		return "Light Blue"
	case Yellow:
		return "Yellow"
	case Lime:
		return "Lime"
	case Pink:
		return "Pink"
	case Gray:
		return "Gray"
	case LightGray:
		return "Light Gray"
	case Cyan:
		return "Cyan"
	case Purple:
		return "Purple"
	case Blue:
		return "Blue"
	case Brown:
		return "Brown"
	case Green:
		return "Green"
	case Red:
		return "Red"
	case Black:
		return "Black"
	default:
		return "Unknown"
	}
}

// GetWoolColorRGBA returns the RGBA color of the given WoolColor.
func (c Color) ToRGBA() color.RGBA {
	return Colors[c]
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

// HSVToRGB converts a hue (0-360), saturation (0-1), and value (0-1) to RGB values.
func HSVToRGB(h, s, v float64) (r, g, b uint8) {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := v - c

	var rf, gf, bf float64
	switch {
	case 0 <= h && h < 60:
		rf, gf, bf = c, x, 0
	case 60 <= h && h < 120:
		rf, gf, bf = x, c, 0
	case 120 <= h && h < 180:
		rf, gf, bf = 0, c, x
	case 180 <= h && h < 240:
		rf, gf, bf = 0, x, c
	case 240 <= h && h < 300:
		rf, gf, bf = x, 0, c
	case 300 <= h && h < 360:
		rf, gf, bf = c, 0, x
	default:
		rf, gf, bf = 0, 0, 0
	}

	// Convert float RGB to uint8 and adjust by m to account for brightness
	r = uint8((rf + m) * 255)
	g = uint8((gf + m) * 255)
	b = uint8((bf + m) * 255)

	return
}
