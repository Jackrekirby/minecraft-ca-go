package core

import "image/color"

type RedstoneTorch struct {
	Direction Direction
	IsPowered bool
}

func (b RedstoneTorch) Type() string {
	return "RedstoneTorch"
}

func (b RedstoneTorch) Update(p Vec3, w *World) (Block, bool) {
	neighbour := w.GetBlock(p.Move(b.Direction.GetOppositeDirection()))

	powerEmittingBlock, canOutputPower := neighbour.(StrongPowerEmittingBlock)

	oldIsPowered := b.IsPowered
	if canOutputPower {
		b.IsPowered = !powerEmittingBlock.OutputsStrongPowerInDirection(b.Direction)
	} else {
		b.IsPowered = true
	}

	hasUpdated := oldIsPowered != b.IsPowered
	return b, hasUpdated
}

func (b RedstoneTorch) OutputsPowerInDirection(d Direction) bool {
	return d != b.Direction.GetOppositeDirection() && b.IsPowered
}

func (b RedstoneTorch) OutputsStrongPowerInDirection(d Direction) bool {
	return d == Up && b.IsPowered
}

func (b RedstoneTorch) OutputsWeakPowerInDirection(d Direction) bool {
	return b.OutputsPowerInDirection(d) && !b.OutputsStrongPowerInDirection(d)
}

func (b RedstoneTorch) ToRune() rune {
	if b.IsPowered {
		return 'T'
	} else {
		return 't'
	}
}

func (b RedstoneTorch) ToCuboids() []Cuboid {
	s := Point3DFromScalar(16)
	var c uint8
	if b.IsPowered {
		c = 255
	} else {
		c = 150
	}

	torch_base := MakeAxisAlignedCuboid(
		Point3D{7, 0, 7}.Divide(s),
		Point3D{9, 9, 9}.Divide(s),
		color.RGBA{160, 127, 81, 255},
	)

	torch_head := MakeAxisAlignedCuboid(
		Point3D{7, 10, 7}.Divide(s),
		Point3D{9, 12, 9}.Divide(s),
		color.RGBA{c, 0, 0, 255},
	)

	cuboids := []Cuboid{
		torch_base,
		torch_head,
	}

	var ry, rz float64 = 0, 0
	var offset Point3D = Point3D{0, 0, 0}
	switch b.Direction {
	case Left:
		ry = 0
		rz = 45
		offset = Point3D{0.25, 0, 0}
	case Right:
		ry = 180
		rz = 45
		offset = Point3D{-0.25, 0, 0}
	case Front:
		ry = 90
		rz = 45
		offset = Point3D{0, 0, -0.25}
	case Back:
		ry = 270
		rz = 45
		offset = Point3D{0, 0, 0.25}
	}

	for j := range cuboids {
		for i := 0; i < 8; i++ {
			cuboid := &cuboids[j]
			translate := Point3D{8, 4.5, 8}.Divide(s)
			cuboid.vertices[i] = cuboid.vertices[i].Subtract(translate).RotateZ(DegToRad(rz)).RotateY(DegToRad(ry)).Add(translate).Add(offset)
		}
	}

	return cuboids

}
