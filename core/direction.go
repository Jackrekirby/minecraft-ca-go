package core

type Direction int

const (
	Up Direction = iota
	Down
	Left  // -x
	Right // +x
	Front // +y
	Back  // -y
)

func (d Direction) String() string {
	return [...]string{"up", "down", "left", "right", "front", "back"}[d]
}

func (d Direction) GetOppositeDirection() Direction {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	case Front:
		return Back
	case Back:
		return Front
	default:
		panic("Direction not implemented")
	}
}

func (d Direction) ToVec3() Vec3 {
	switch d {
	case Up:
		return Vec3{0, 1, 0} // Up is positive y-axis
	case Down:
		return Vec3{0, -1, 0} // Down is negative y-axis
	case Left:
		return Vec3{-1, 0, 0} // Left is negative x-axis
	case Right:
		return Vec3{1, 0, 0} // Right is positive x-axis
	case Front:
		return Vec3{0, 0, 1} // Front is positive z-axis
	case Back:
		return Vec3{0, 0, -1} // Back is negative z-axis
	default:
		panic("Direction not implemented")
	}
}
