package core

import "math"

func FastSin(x float64) float64 {
	// Range reduction to [-π, π]
	for x > math.Pi {
		x -= 2 * math.Pi
	}
	for x < -math.Pi {
		x += 2 * math.Pi
	}

	// Approximation
	if x < 0 {
		return 1.27323954*x + 0.405284735*x*x
	}
	return 1.27323954*x - 0.405284735*x*x
}

func FastCos(x float64) float64 {
	// Range reduction to [-π, π]
	for x > math.Pi {
		x -= 2 * math.Pi
	}
	for x < -math.Pi {
		x += 2 * math.Pi
	}

	x *= x
	return 1 - 0.5*x + 0.0416666664*x*x
}

// Define the size of the lookup table.
const tableSize = 1024

// Precomputed sine values.
var sinTable [tableSize]float64

// Initializes the lookup table with sine values.
func InitSinTable() {
	for i := 0; i < tableSize; i++ {
		angle := float64(i) * (2 * math.Pi / float64(tableSize))
		sinTable[i] = math.Sin(angle)
	}
}

func FastSin2(x float64) float64 {
	// Normalize the input angle to [0, 2π].
	const pi2 = 2 * math.Pi
	x = math.Mod(x, pi2)
	if x < 0 {
		x += 2 * math.Pi
	}

	// Convert the input angle to a table index.
	index := x * float64(tableSize) / pi2
	i := int(index) // Lower index
	// fraction := index - float64(i) // Fractional part for interpolation
	// iNext := (i + 1) % tableSize   // Wrap around for the next index

	// Linearly interpolate between the nearest values.
	return sinTable[i] // (1-fraction)*sinTable[i] + fraction*sinTable[iNext]
}

// Fast cosine approximation using the sine function.
func FastCos2(x float64) float64 {
	// Use the identity cos(x) = sin(x + π/2).
	return FastSin2(x + math.Pi/2)
}

func Cos(x float64) float64 {
	return math.Cos(x)
}

func Sin(x float64) float64 {
	return math.Sin(x)
}
