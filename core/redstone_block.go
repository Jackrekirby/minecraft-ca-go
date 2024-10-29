package core

import "image/color"

type RedstoneBlock struct {
}

func (b RedstoneBlock) Type() string {
	return "RedstoneBlock"
}

func (b RedstoneBlock) OutputsPowerInDirection(d Direction) bool {
	return true
}

func (b RedstoneBlock) OutputsStrongPowerInDirection(d Direction) bool {
	return true
}

func (b RedstoneBlock) ToRune() rune {
	return 'B'
}

func (b RedstoneBlock) ToCuboids() []Cuboid {
	return []Cuboid{
		MakeAxisAlignedCuboid(
			Point3D{0, 0, 0},
			Point3D{1, 1, 1},
			color.RGBA{255, 0, 0, 255},
		),
	}
}
