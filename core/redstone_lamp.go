package core

import "image/color"

type PowerType int

const (
	Strong PowerType = iota
	Weak
	None
)

type RedstoneLamp struct {
	InputPowerType PowerType
}

func (b RedstoneLamp) Type() string {
	return "RedstoneLamp"
}

func (b RedstoneLamp) OutputsPowerInDirection(d Direction) bool {
	return b.InputPowerType == Strong
}

func (b RedstoneLamp) OutputsStrongPowerInDirection(d Direction) bool {
	return false
}

func (b RedstoneLamp) OutputsWeakPowerInDirection(d Direction) bool {
	return b.InputPowerType == Strong
}

func UpdateInputPowerType(p Vec3, w *World) PowerType {
	var inputPowerType PowerType = None
	// Up, Down, Left, Right, Front, Back
	for _, d := range [...]Direction{Up, Down, Left, Right, Front, Back} {
		neighbour := w.GetBlock(p.Move(d.GetOppositeDirection()))

		strongPowerEmittingBlock, canOutputStrongPower := neighbour.(StrongPowerEmittingBlock)

		if canOutputStrongPower && inputPowerType != Strong && strongPowerEmittingBlock.OutputsStrongPowerInDirection(d) {
			inputPowerType = Strong
			break
		} else {
			weakPowerEmittingBlock, canOutputWeakPower := neighbour.(WeakPowerEmittingBlock)
			if canOutputWeakPower && inputPowerType == None && weakPowerEmittingBlock.OutputsWeakPowerInDirection(d) {
				inputPowerType = Weak
			}
		}
	}
	return inputPowerType
}

func (b RedstoneLamp) SubUpdate(p Vec3, w *World) (Block, bool) {
	var newInputPowerType PowerType = UpdateInputPowerType(p, w)
	hasUpdated := newInputPowerType != b.InputPowerType
	b.InputPowerType = newInputPowerType
	return b, hasUpdated
}

func (b RedstoneLamp) isPowered() bool {
	return b.InputPowerType != None
}

func (b RedstoneLamp) ToRune() rune {
	return 'B'
}

func (b RedstoneLamp) ToCuboids() []Cuboid {
	var c color.RGBA
	if b.InputPowerType == Strong {
		c = color.RGBA{219, 171, 115, 255}
	} else if b.InputPowerType == Weak {
		c = color.RGBA{0, 255, 0, 255}
		// c = color.RGBA{219, 171, 115, 255}
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
