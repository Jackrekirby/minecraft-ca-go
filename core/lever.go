package core

import "image/color"

type Lever struct {
	direction Direction
	isOn      bool
}

func (b Lever) Type() string {
	return "Lever"
}

func (b Lever) OutputsPowerInDirection(d Direction) bool {
	return b.isOn
}

func (b Lever) OutputsStrongPowerInDirection(d Direction) bool {
	return b.isOn && b.direction == d
}

func (b Lever) ToCuboids() []Cuboid {
	s := Point3DFromScalar(16)
	stick := MakeAxisAlignedCuboid(
		Point3D{7, 3, 7}.Divide(s),
		Point3D{9, 11, 9}.Divide(s),
		color.RGBA{160, 127, 81, 255},
	)
	var rx float64
	if b.isOn {
		rx = 45
	} else {
		rx = -45
	}
	for i := 0; i < 8; i++ {
		cuboid := &stick
		translate := Point3D{7, 3, 7}.Divide(s)
		cuboid.vertices[i] = cuboid.vertices[i].Subtract(translate).RotateX(DegToRad(rx)).Add(translate)
	}
	return []Cuboid{
		MakeAxisAlignedCuboid(
			Point3D{5, 0, 4}.Divide(s),
			Point3D{11, 3, 12}.Divide(s),
			color.RGBA{100, 100, 100, 255},
		),
		stick,
	}
}
