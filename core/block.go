package core

type Block interface {
	Type() string
}

type UpdateableBlock interface {
	Update(p Vec3, w *World) (Block, bool)
}

type SubUpdateableBlock interface {
	SubUpdate(p Vec3, w *World) (Block, bool)
}

type PowerEmittingBlock interface {
	OutputsPowerInDirection(d Direction) bool
}

type StrongPowerEmittingBlock interface {
	OutputsStrongPowerInDirection(d Direction) bool
}

type WeakPowerEmittingBlock interface {
	OutputsWeakPowerInDirection(d Direction) bool
}

type RenderableBlock interface {
	ToRune() rune
}

type WireRenderBlock interface {
	ToCuboids(scene *Scene) []Cuboid
}

type OpaqueBlock interface {
	IsOpaqueInDirection(d Direction) bool
}

type DirectionalBlock interface {
	Type() string // Block interface
	GetDirection() Direction
	SetDirection(d Direction) DirectionalBlock // would love for this to mutate
}
