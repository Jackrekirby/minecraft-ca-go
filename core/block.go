package core

type Block interface {
	Type() string
}

type UpdateableBlock interface {
	Update(p Vec3, w *World) (Block, bool)
}

type PowerEmittingBlock interface {
	OutputsPowerInDirection(d Direction) bool
}

type RenderableBlock interface {
	ToRune() rune
}

type WireRenderBlock interface {
	ToCuboids() []Cuboid
}
