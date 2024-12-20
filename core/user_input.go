package core

import "fmt"

func createSimpleWorld(world *World) {
	// levers arond a lamp
	world.SetBlock(Vec3{X: 2, Y: 2, Z: 2}, RedstoneLamp{InputPowerType: None})
	world.SetBlock(Vec3{X: 2, Y: 3, Z: 2}, RedstoneLamp{InputPowerType: None})
	world.SetBlock(Vec3{X: 3, Y: 2, Z: 2}, RedstoneLamp{InputPowerType: None})
	world.SetBlock(Vec3{X: 1, Y: 2, Z: 2}, RedstoneLamp{InputPowerType: None})
	// world.SetBlock(p.Move(Up), WoolBlock{Cyan, None})
	// world.SetBlock(p.Move(Up), RedstoneTorch{Direction: Up, IsPowered: true})

	// for _, d := range [6]Direction{Left, Right, Front, Back, Up, Down} {
	// 	world.SetBlock(
	// 		p.Move(d),
	// 		Lever{Direction: d, IsOn: false},
	// 	)
	// }
}

func createWorld(world *World) {
	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			p := Vec3{X: x, Y: 0, Z: z}
			world.SetBlock(
				p,
				WoolBlock{Color(x % 16), None},
			)
		}
	}

	// levers arond a lamp
	p := Vec3{X: 2, Y: 1, Z: 13}
	world.SetBlock(p, RedstoneLamp{InputPowerType: None})
	for _, d := range [6]Direction{Left, Right, Front, Back, Up, Down} {
		world.SetBlock(
			p.Move(d),
			Lever{Direction: d, IsOn: false},
		)
	}

	// torches around a redstone block
	cp := Vec3{X: 5, Y: 1, Z: 13}
	world.SetBlock(cp, RedstoneBlock{})
	for _, d := range [4]Direction{Left, Right, Front, Back} {
		world.SetBlock(
			cp.Move(d),
			RedstoneTorch{Direction: d, IsPowered: true},
		)
	}

	// torch horizontal
	for i := 0; i < 4; i++ {
		world.SetBlock(
			Vec3{X: 3 + i, Y: 8, Z: 2},
			RedstoneTorch{Direction: Up, IsPowered: false},
		)
	}

	// torch tower lamps on top
	p = Vec3{X: 2, Y: 7, Z: 2}
	world.SetBlock(
		p.Add(Vec3{1, 0, 0}),
		RedstoneLamp{InputPowerType: Weak},
	)
	world.SetBlock(
		p.Add(Vec3{2, 0, 0}),
		RedstoneLamp{InputPowerType: Strong},
	)
	world.SetBlock(
		p.Add(Vec3{3, 0, 0}),
		RedstoneLamp{InputPowerType: Weak},
	)
	world.SetBlock(
		p.Add(Vec3{4, 0, 0}),
		RedstoneLamp{InputPowerType: None},
	)

	world.SetBlock(Vec3{X: 2, Y: 2, Z: 2}, Lever{Left, false})
	for i := 2; i < 7; i++ {
		p := Vec3{X: 3, Y: i, Z: 2}
		if i%2 == 0 {
			world.SetBlock(p, RedstoneLamp{InputPowerType: None})
			world.SetBlock(
				p.Move(Right),
				RedstoneTorch{Direction: Right, IsPowered: true},
			)
		} else {
			world.SetBlock(p.Move(Right), RedstoneLamp{InputPowerType: Strong})
			world.SetBlock(
				p,
				RedstoneTorch{Direction: Left, IsPowered: false},
			)
		}
	}

	// vertical invalid torch tower
	world.SetBlock(Vec3{X: 12, Y: 2, Z: 2}, Lever{Left, false})
	world.SetBlock(Vec3{X: 13, Y: 2, Z: 2}, WoolBlock{Cyan, None})
	for i := 3; i < 8; i++ {
		p := Vec3{X: 13, Y: i, Z: 2}
		torch := RedstoneTorch{Direction: Up, IsPowered: i%2 == 1}
		world.SetBlock(p, torch)
		if i%2 == 0 {
			world.SetBlock(p.Move(Left), RedstoneLamp{InputPowerType: None})
		}
	}
}

func toggleLever(p Vec3, world *World) bool {
	b := world.GetBlock(
		p,
	)
	lever, isLever := b.(Lever)
	if isLever {
		lever.IsOn = !lever.IsOn
		world.SetBlock(p, lever)
		return true
	}
	return false
}

func ProcessUserInputs(iteration int, world *World) bool {
	// currently just handles programatic changes to the world to simulate user interaction
	var hasAnyBlockUpdated bool = false
	if iteration == 0 {
		createWorld(world)
		// createSimpleWorld(world)
		hasAnyBlockUpdated = true
	}
	// if iteration%32 == 4 || iteration%32 == 20 {
	// 	if toggleLever(Vec3{X: 0, Y: 2, Z: 2}, world) {
	// 		hasAnyBlockUpdated = true
	// 	}
	// 	if toggleLever(Vec3{X: 2, Y: 2, Z: 2}, world) {
	// 		hasAnyBlockUpdated = true
	// 	}
	// }
	// {
	// 	d := [6]Direction{Left, Right, Front, Back, Up, Down}[iteration%12/2]
	// 	p := Vec3{X: 2, Y: 1, Z: 13}.Move(d)
	// 	if toggleLever(p, world) {
	// 		hasAnyBlockUpdated = true
	// 	}
	// }

	return hasAnyBlockUpdated
}

func HandleKeyPress(scene *Scene, key string, moveDelta float64, rotDelta float64) {
	camera := &scene.Camera
	// delta := 0.5
	// rotation := DegToRad(15)
	// Handle key press
	switch key {
	case "q":
		fmt.Println("Exiting...")
		scene.GameState = Quit
	case "p":
		if scene.GameState == Paused {
			scene.GameState = Playing
		} else if scene.GameState == Playing {
			scene.GameState = Paused
		}
	case "o":
		if scene.GameState == Paused {
			scene.GameState = Pausing
		} else if scene.GameState == Playing {
			scene.GameState = Paused
		}
	case "r":
		scene.World = World{}
		scene.Iteration = 0
		createWorld(&scene.World)
	case "w":
		camera.Position = camera.Position.Add(Point3D{0, 0, moveDelta}.RotateY(-camera.Rotation.Y))
	case "a":
		camera.Position = camera.Position.Add(Point3D{-moveDelta, 0, 0}.RotateY(-camera.Rotation.Y))
	case "s":
		camera.Position = camera.Position.Add(Point3D{0, 0, -moveDelta}.RotateY(-camera.Rotation.Y))
	case "d":
		camera.Position = camera.Position.Add(Point3D{moveDelta, 0, 0}.RotateY(-camera.Rotation.Y))
	case "e":
		camera.Position = camera.Position.Add(Point3D{0, moveDelta, 0})
	case "c":
		camera.Position = camera.Position.Add(Point3D{0, -moveDelta, 0})
	case "z":
		camera.Rotation.Y = camera.Rotation.Y + rotDelta
	case "x":
		camera.Rotation.Y = camera.Rotation.Y - rotDelta
	default:
		fmt.Println("Pressed:", key)
	}
}
