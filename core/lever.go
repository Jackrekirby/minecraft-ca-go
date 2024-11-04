package core

import "image/color"

type Lever struct {
	Direction Direction
	IsOn      bool
}

func (b Lever) Type() string {
	return "Lever"
}

func (b Lever) OutputsPowerInDirection(d Direction) bool {
	return b.IsOn
}

func (b Lever) OutputsStrongPowerInDirection(d Direction) bool {
	return b.IsOn && b.Direction == d.GetOppositeDirection()
}

func (b Lever) ToCuboids(scene *Scene) []Cuboid {
	s := Point3DFromScalar(16)
	stick := MakeAxisAlignedCuboid(
		Point3D{7, 3, 7}.Divide(s),
		Point3D{9, 11, 9}.Divide(s),
		color.RGBA{160, 127, 81, 255},
		MakeCuboidUVsForSingleTexture("oak_planks", scene),
	)
	var rx float64
	if b.IsOn {
		rx = 45
	} else {
		rx = -45
	}
	for i := 0; i < 8; i++ {
		cuboid := &stick
		translate := Point3D{7, 3, 7}.Divide(s)
		cuboid.vertices[i] = cuboid.vertices[i].Subtract(translate).RotateX(DegToRad(rx)).Add(translate)
	}

	base := MakeAxisAlignedCuboid(
		Point3D{5, 0, 4}.Divide(s),
		Point3D{11, 3, 12}.Divide(s),
		color.RGBA{100, 100, 100, 255},
		MakeCuboidUVsForSingleTexture("stone", scene),
	)

	cuboids := []Cuboid{
		base,
		stick,
	}

	var rot Point3D
	switch b.Direction {
	case Left:
		rot = Point3D{90, 0, 90}
	case Right:
		rot = Point3D{90, 0, -90}
	case Up:
		rot = Point3D{0, 0, 0}
	case Down:
		rot = Point3D{0, 0, 180}
	case Front:
		rot = Point3D{90, 0, 0}
	case Back:
		rot = Point3D{90, 180, 0}
	}

	for j := 0; j < 2; j++ {
		for i := 0; i < 8; i++ {
			cuboid := &cuboids[j]
			translate := Point3D{8, 8, 8}.Divide(s)
			cuboid.vertices[i] = cuboid.vertices[i].
				Subtract(translate).
				RotateZ(DegToRad(rot.Z)).
				RotateX(DegToRad(rot.X)).
				RotateY(DegToRad(rot.Y)).
				Add(translate)
		}
	}

	return cuboids
}
