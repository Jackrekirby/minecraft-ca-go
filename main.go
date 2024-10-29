package main

import (
	"fmt"
	"project_two/core"
	"time"
)

func printWorld(w *core.World, n int) {
	for y := n - 1; y >= 0; y-- {
		p := core.Vec3{X: 1, Y: y, Z: 0}
		b := w.GetBlock(p)
		rb, canRender := b.(core.RenderableBlock)
		if canRender {
			fmt.Println(y, string(rb.ToRune()))
		} else {
			fmt.Println(y, " ")
		}
	}
}

func runWorld() {
	world := core.World{}

	quitGame := false

	camera := core.Camera{
		Position:    core.Point3D{X: 3.5, Y: 5.5, Z: -4},
		Rotation:    core.Point3D{X: core.DegToRad(0), Y: core.DegToRad(0), Z: core.DegToRad(0)},
		FOV:         90.0,
		AspectRatio: 1.0,
		Near:        0.1,
		Far:         100.0,
	}

	go core.KeyboardEvents(&camera, &quitGame)

	cp := core.Vec3{X: 8, Y: 2, Z: 4}
	world.SetBlock(cp, core.RedstoneBlock{})
	for _, d := range [4]core.Direction{core.Left, core.Right, core.Front, core.Back} {
		world.SetBlock(
			cp.Move(d),
			core.RedstoneTorch{Direction: d, IsPowered: true},
		)
	}

	for i := 2; i < 7; i++ {
		p := core.Vec3{X: 3, Y: i, Z: 0}
		if i%2 == 0 {
			world.SetBlock(p, core.RedstoneLamp{InputPowerType: core.None})
			world.SetBlock(
				p.Move(core.Right),
				core.RedstoneTorch{Direction: core.Right, IsPowered: true},
			)
		} else {
			world.SetBlock(p.Move(core.Right), core.RedstoneLamp{InputPowerType: core.Strong})
			world.SetBlock(
				p,
				core.RedstoneTorch{Direction: core.Left, IsPowered: false},
			)
		}
	}

	for i := 3; i < 8; i++ {
		p := core.Vec3{X: 1, Y: i, Z: 0}
		torch := core.RedstoneTorch{Direction: core.Up, IsPowered: i%2 == 1}
		world.SetBlock(p, torch)
		if i%2 == 0 {
			world.SetBlock(p.Move(core.Left), core.RedstoneLamp{InputPowerType: core.None})
		}

	}

	// fmt.Println("World [ - ]:")
	core.DrawScene(&camera, &world)
	// printWorld(&world, n)
	time.Sleep(500 * time.Millisecond)
	iteration := 0
	for !quitGame {
		hasAnyBlockUpdated := world.UpdateWorld()

		if iteration%32 == 4 {
			p := core.Vec3{X: 1, Y: 2, Z: 0}
			world.SetBlock(p, core.RedstoneBlock{})
			world.SetBlock(p.Add(core.Vec3{X: 2, Y: 0, Z: 0}), core.RedstoneBlock{})
			hasAnyBlockUpdated = true
		} else if iteration%32 == 20 {
			p := core.Vec3{X: 1, Y: 2, Z: 0}
			world.SetBlock(p, core.Air{})
			world.SetBlock(p.Add(core.Vec3{X: 2, Y: 0, Z: 0}), core.Air{})
			hasAnyBlockUpdated = true
		}

		// if !hasAnyBlockUpdated {
		// 	fmt.Println("No block updates")
		// 	// break
		// }

		fmt.Println(iteration, hasAnyBlockUpdated)
		// printWorld(&world, n)
		core.DrawScene(&camera, &world)
		time.Sleep(500 * time.Millisecond)
		iteration = iteration + 1
	}
}

func main() {
	runWorld()
	// core.AnimateChequerBoard()
	// tilemap, err := core.GenerateTilemap("./assets", 32)
	// if err != nil {
	// 	panic(err)
	// }
	// core.SaveImage(tilemap, "tilemap.png")
}
