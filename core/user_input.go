package core

func createWorld(world *World) {
	cp := Vec3{X: 8, Y: 2, Z: 4}
	world.SetBlock(cp, RedstoneBlock{})
	for _, d := range [4]Direction{Left, Right, Front, Back} {
		world.SetBlock(
			cp.Move(d),
			RedstoneTorch{Direction: d, IsPowered: true},
		)
	}

	for i := 2; i < 7; i++ {
		p := Vec3{X: 3, Y: i, Z: 0}
		if i%2 == 0 {
			world.SetBlock(p, RedstoneLamp{InputPowerType: None})
			world.SetBlock(
				p.Move(Right),
				RedstoneTorch{Direction: Right, IsPowered: true},
			)
		} else {
			world.SetBlock(p.Move(Right), RedstoneLamp{InputPowerType: None})
			world.SetBlock(
				p,
				RedstoneTorch{Direction: Left, IsPowered: false},
			)
		}
	}
	for i := 0; i < 4; i++ {
		world.SetBlock(
			Vec3{X: 3 + i, Y: 7, Z: 0},
			RedstoneLamp{InputPowerType: None},
		)
		world.SetBlock(
			Vec3{X: 3 + i, Y: 8, Z: 0},
			RedstoneTorch{Direction: Up, IsPowered: false},
		)
	}

	for i := 3; i < 8; i++ {
		p := Vec3{X: 1, Y: i, Z: 0}
		torch := RedstoneTorch{Direction: Up, IsPowered: i%2 == 1}
		world.SetBlock(p, torch)
		if i%2 == 0 {
			world.SetBlock(p.Move(Left), RedstoneLamp{InputPowerType: None})
		}

	}
}

func ProcessUserInputs(iteration int, world *World) bool {
	// currently just handles programatic changes to the world to simulate user interaction
	var hasAnyBlockUpdated bool = false
	if iteration == 0 {
		createWorld(world)
		hasAnyBlockUpdated = true
	} else if iteration%32 == 4 {
		p := Vec3{X: 1, Y: 2, Z: 0}
		world.SetBlock(p, RedstoneBlock{})
		world.SetBlock(p.Add(Vec3{X: 2, Y: 0, Z: 0}), RedstoneBlock{})
		hasAnyBlockUpdated = true
	} else if iteration%32 == 20 {
		p := Vec3{X: 1, Y: 2, Z: 0}
		world.SetBlock(p, Air{})
		world.SetBlock(p.Add(Vec3{X: 2, Y: 0, Z: 0}), Air{})
		hasAnyBlockUpdated = true
	}
	return hasAnyBlockUpdated
}
