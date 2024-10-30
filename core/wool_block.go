package core

type WoolBlock struct {
	Color          Color
	InputPowerType PowerType
}

func (b WoolBlock) Type() string {
	return "WoolBlock"
}

func (b WoolBlock) OutputsPowerInDirection(d Direction) bool {
	return b.InputPowerType != None
}

func (b WoolBlock) OutputsStrongPowerInDirection(d Direction) bool {
	return b.InputPowerType == Strong
}

func (b WoolBlock) OutputsWeakPowerInDirection(d Direction) bool {
	return b.InputPowerType == Weak
}

func (b WoolBlock) Update(p Vec3, w *World) (Block, bool) {
	var newInputPowerType PowerType = UpdateInputPowerType(p, w)
	hasUpdated := newInputPowerType != b.InputPowerType
	b.InputPowerType = newInputPowerType
	return b, hasUpdated
}

func (b WoolBlock) ToCuboids() []Cuboid {
	return []Cuboid{
		MakeAxisAlignedCuboid(
			Point3D{0, 0, 0},
			Point3D{1, 1, 1},
			b.Color.ToRGBA(),
		),
	}
}