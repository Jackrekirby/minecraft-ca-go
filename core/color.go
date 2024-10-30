package core

import (
	"image/color"
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
