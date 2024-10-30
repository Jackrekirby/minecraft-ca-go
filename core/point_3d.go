package core

import "math"

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
