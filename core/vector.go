package core

import "fmt"

type Vec3 struct {
	X int
	Y int
	Z int
}

func (v Vec3) ToDirection() Direction {
	switch v {
	case Vec3{0, 1, 0}:
		return Up
	case Vec3{0, -1, 0}:
		return Down
	case Vec3{-1, 0, 0}:
		return Left
	case Vec3{1, 0, 0}:
		return Right
	case Vec3{0, 0, 1}:
		return Front
	case Vec3{0, 0, -1}:
		return Back
	default:
		fmt.Println(v)
		panic("Vec3 does not map to a valid Direction")
	}
}

func (v1 Vec3) Equals(v2 Vec3) bool {
	return v1.X == v2.X && v1.Y == v2.Y && v1.Z == v2.Z
}

func (v Vec3) ToPoint3D() Point3D {
	return Point3D{
		X: float64(v.X),
		Y: float64(v.Y),
		Z: float64(v.Z),
	}
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
