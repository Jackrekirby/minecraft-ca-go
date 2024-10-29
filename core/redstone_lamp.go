package core

import "image/color"

type RedstoneLamp struct {
	IsPowered bool
}

func (b RedstoneLamp) Type() string {
	return "RedstoneLamp"
}

func (b RedstoneLamp) OutputsPowerInDirection(d Direction) bool {
	return b.IsPowered
}

func (b RedstoneLamp) Update(p Vec3, w *World) (Block, bool) {
	var hasUpdated bool = false
	// Up, Down, Left, Right, Front, Back
	for _, d := range [...]Direction{Up, Down, Left, Right, Front, Back} {
		neighbour := w.GetBlock(p.Move(d.GetOppositeDirection()))

		powerEmittingBlock, canOutputPower := neighbour.(PowerEmittingBlock)

		var newIsPowered bool
		if canOutputPower {
			newIsPowered = powerEmittingBlock.OutputsPowerInDirection(d)
		} else {
			newIsPowered = false
		}

		if newIsPowered {
			hasUpdated = newIsPowered != b.IsPowered
			b.IsPowered = newIsPowered
			return b, hasUpdated
		}
	}
	// lamp has not been powered from any direction
	hasUpdated = !b.IsPowered
	return b, hasUpdated
}

func (b RedstoneLamp) ToRune() rune {
	return 'B'
}

func (b RedstoneLamp) ToCuboids() []Cuboid {
	var c color.RGBA
	if b.IsPowered {
		c = color.RGBA{219, 171, 115, 255}
	} else {
		c = color.RGBA{95, 59, 34, 255}
	}

	return []Cuboid{
		MakeAxisAlignedCuboid(
			Point3D{0, 0, 0},
			Point3D{1, 1, 1},
			c,
		),
	}
}
