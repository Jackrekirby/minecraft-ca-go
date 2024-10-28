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

	powerEmittingBlock, canOutputPower := neighbour.(PowerEmittingBlock)

	oldIsPowered := b.IsPowered
	if canOutputPower {
		b.IsPowered = !powerEmittingBlock.OutputsPowerInDirection(b.Direction)
	} else {
		b.IsPowered = true
	}

	hasUpdated := oldIsPowered != b.IsPowered
	return b, hasUpdated
}

func (b RedstoneTorch) OutputsPowerInDirection(d Direction) bool {
	return d != b.Direction.GetOppositeDirection() && b.IsPowered
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
	return []Cuboid{
		{
			Point3D{7, 0, 7}.Divide(s),
			Point3D{9, 9, 9}.Divide(s),
			color.RGBA{160, 127, 81, 255},
		},
		{
			Point3D{7, 10, 7}.Divide(s),
			Point3D{9, 12, 9}.Divide(s),
			color.RGBA{c, 0, 0, 255},
		},
	}
}
