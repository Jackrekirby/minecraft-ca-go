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

	world.SetBlock(core.Vec3{X: 4, Y: 1, Z: 2}, core.RedstoneBlock{})
	world.SetBlock(
		core.Vec3{X: 5, Y: 1, Z: 2},
		core.RedstoneTorch{Direction: core.Right, IsPowered: true},
	)
	world.SetBlock(
		core.Vec3{X: 3, Y: 1, Z: 2},
		core.RedstoneTorch{Direction: core.Left, IsPowered: true},
	)
	world.SetBlock(
		core.Vec3{X: 4, Y: 1, Z: 3},
		core.RedstoneTorch{Direction: core.Front, IsPowered: true},
	)
	world.SetBlock(
		core.Vec3{X: 4, Y: 1, Z: 1},
		core.RedstoneTorch{Direction: core.Back, IsPowered: true},
	)

	n := 6
	for i := 1; i < n; i++ {
		p := core.Vec3{X: 1, Y: i, Z: 0}
		torch := core.RedstoneTorch{Direction: core.Up, IsPowered: i%2 == 1}
		world.SetBlock(p, torch)
		if i%2 == 0 {
			world.SetBlock(p.Move(core.Right), core.RedstoneLamp{IsPowered: false})
		}

	}

	// fmt.Println("World [ - ]:")
	core.DrawScene(&world)
	// printWorld(&world, n)
	time.Sleep(500 * time.Millisecond)
	for it := 0; it < 10; it++ {
		hasAnyBlockUpdated := world.UpdateWorld()

		if it == 0 {
			p := core.Vec3{X: 1, Y: 0, Z: 0}
			world.SetBlock(p, core.RedstoneBlock{})
			hasAnyBlockUpdated = true
		}

		if !hasAnyBlockUpdated {
			fmt.Println("No block updates")
			break
		}

		fmt.Println("World [", it, "]:")
		// printWorld(&world, n)
		core.DrawScene(&world)
		time.Sleep(500 * time.Millisecond)
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
