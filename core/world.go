package core

const (
	WorldWidth  = 16
	WorldHeight = 16
	WorldDepth  = 16
)

type World struct {
	Blocks [WorldWidth * WorldHeight * WorldDepth]Block
}

func (w *World) GetIndex(p Vec3) int {
	return p.Z*WorldWidth*WorldHeight + p.Y*WorldWidth + p.X
}

func (w *World) GetBlock(p Vec3) Block {
	if p.InRange(*Vec3FromScalar(0), *Vec3FromScalar(16)) {
		return Air{}
	}
	block := w.Blocks[w.GetIndex(p)]
	if block == nil {
		return Air{}
	}
	return block
}

func (w *World) SetBlock(p Vec3, block Block) bool {
	if p.InRange(*Vec3FromScalar(0), *Vec3FromScalar(16)) {
		return false
	}
	w.Blocks[w.GetIndex(p)] = block
	return true
}

func (w *World) UpdateBlock(p Vec3) (Block, bool) {
	b := w.GetBlock(p)
	ub, canUpdate := b.(UpdateableBlock)
	if canUpdate {
		return ub.Update(p, w)
	}
	return b, false
}

func (w *World) UpdateWorld() int {
	nextWorld := World{}
	numUpdates := 0
	for x := 0; x < WorldWidth; x++ {
		for y := 0; y < WorldHeight; y++ {
			for z := 0; z < WorldDepth; z++ {
				p := Vec3{X: x, Y: y, Z: z}
				block, hasUpdated := w.UpdateBlock(p)
				if hasUpdated {
					numUpdates += 1
				}
				nextWorld.SetBlock(p, block)
			}
		}
	}
	*w = nextWorld
	return numUpdates
}
