package core

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

func createWorld(world *World) {
	world.SetBlock(
		Vec3{X: 0, Y: 1, Z: 0},
		Lever{Up, true},
	)

	cp := Vec3{X: 8, Y: 2, Z: 4}
	world.SetBlock(cp, RedstoneBlock{})
	for _, d := range [4]Direction{Left, Right, Front, Back} {
		world.SetBlock(
			cp.Move(d),
			RedstoneTorch{Direction: d, IsPowered: true},
		)
	}

	for i := 0; i < 4; i++ {
		world.SetBlock(
			Vec3{X: 3 + i, Y: 8, Z: 0},
			RedstoneTorch{Direction: Up, IsPowered: false},
		)
	}

	world.SetBlock(
		Vec3{X: 3, Y: 7, Z: 0},
		RedstoneLamp{InputPowerType: Weak},
	)
	world.SetBlock(
		Vec3{X: 4, Y: 7, Z: 0},
		RedstoneLamp{InputPowerType: Strong},
	)
	world.SetBlock(
		Vec3{X: 5, Y: 7, Z: 0},
		RedstoneLamp{InputPowerType: Weak},
	)
	world.SetBlock(
		Vec3{X: 6, Y: 7, Z: 0},
		RedstoneLamp{InputPowerType: None},
	)

	for i := 2; i < 7; i++ {
		p := Vec3{X: 3, Y: i, Z: 0}
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
	}
	if iteration%8 == 0 {
		p := Vec3{X: 0, Y: 1, Z: 0}
		b := world.GetBlock(p)
		lever, isLever := b.(Lever)
		if isLever {
			lever.isOn = !lever.isOn
			world.SetBlock(p, lever)
			hasAnyBlockUpdated = true
		}
	}
	if iteration%32 == 4 {
		p := Vec3{X: 1, Y: 2, Z: 0}
		world.SetBlock(p, RedstoneBlock{})
		world.SetBlock(p.Add(Vec3{X: 2, Y: 0, Z: 0}), RedstoneBlock{})
		hasAnyBlockUpdated = true
	}
	if iteration%32 == 20 {
		p := Vec3{X: 1, Y: 2, Z: 0}
		world.SetBlock(p, Air{})
		world.SetBlock(p.Add(Vec3{X: 2, Y: 0, Z: 0}), Air{})
		hasAnyBlockUpdated = true
	}
	return hasAnyBlockUpdated
}

func KeyboardEvents(scene *Scene) {
	// Open the keyboard
	err := keyboard.Open()
	if err != nil {
		fmt.Println("Error opening keyboard:", err)
		return
	}
	defer keyboard.Close()

	fmt.Println("Listening for keyboard inputs. Press 'q' to quit.")

	camera := &scene.Camera

	for {
		// Read key press
		key, _, err := keyboard.GetKey()
		if err != nil {
			fmt.Println("Error reading key:", err)
			break
		}
		delta := 1.0
		rotation := DegToRad(15)
		// Handle key press
		switch key {
		case 'q':
			fmt.Println("Exiting...")
			scene.GameState = Quit
		case 'p':
			if scene.GameState == Paused {
				scene.GameState = Playing
			} else if scene.GameState == Playing {
				scene.GameState = Paused
			}
		case 'o':
			if scene.GameState == Paused {
				scene.GameState = Pausing
			} else if scene.GameState == Playing {
				scene.GameState = Paused
			}
		case 'r':
			scene.World = World{}
			scene.Iteration = 0
			createWorld(&scene.World)
		case 'w':
			camera.Position = camera.Position.Add(Point3D{0, 0, delta}.RotateY(-camera.Rotation.Y))
		case 'a':
			camera.Position = camera.Position.Add(Point3D{-delta, 0, 0}.RotateY(-camera.Rotation.Y))
		case 's':
			camera.Position = camera.Position.Add(Point3D{0, 0, -delta}.RotateY(-camera.Rotation.Y))
		case 'd':
			camera.Position = camera.Position.Add(Point3D{delta, 0, 0}.RotateY(-camera.Rotation.Y))
		case 'e':
			camera.Position = camera.Position.Add(Point3D{0, 1, 0})
		case 'c':
			camera.Position = camera.Position.Add(Point3D{0, -1, 0})
		case 'z':
			camera.Rotation.Y = camera.Rotation.Y + rotation
		case 'x':
			camera.Rotation.Y = camera.Rotation.Y - rotation
		default:
			fmt.Println("Pressed:", key)
		}
	}
}
