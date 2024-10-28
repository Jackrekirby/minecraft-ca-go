package core

import "fmt"

type Vec3 struct {
	X int
	Y int
	Z int
}

func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
	}
}

func (v Vec3) Subtract(other Vec3) Vec3 {
	return Vec3{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
	}
}

func (v Vec3) InRange(min Vec3, max Vec3) bool {
	return v.X < min.X || v.X >= max.X || v.Y < min.Y || v.Y >= max.Y || v.Z < min.Z || v.Z >= max.Z
}

func Vec3FromScalar(scalar int) *Vec3 {
	return &Vec3{scalar, scalar, scalar}
}

func (v Vec3) String() string {
	return fmt.Sprintf("Vec3(%d, %d, %d)", v.X, v.Y, v.Z)
}

func (v Vec3) Move(d Direction) Vec3 {
	return v.Add(d.ToVec3())
}
