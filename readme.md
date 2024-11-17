# Minecraft Cellular Automata in Go

## About

A recreation of minecraft redstone as a cellular automata in golang.

View the game on the web:
[https://jackrekirby.github.io/minecraft-ca-go/](https://jackrekirby.github.io/minecraft-ca-go/)


## Building

Use `build_wasm.ps1` to build for web assembly.

Use `build_windows.ps1` to build for windows.

Use `go run .` or `& "./main.exe"` to run on windows.

## Third Party Dependencies

Controls on window are done via the terminal using:

`go get -u github.com/eiannone/keyboard`

## Profling

`go tool pprof -http=":8080" cpu.pprof` (requires graphviz)

`go test -bench BenchmarkDrawTriangle3D project_two/core -cpuprofile cpu.prof`

`go test -bench ^BenchmarkDrawTriangle3D$ project_two/core -benchmem -run=^$ -cpuprofile cpu.prof -trace trace.out`

`go test -bench=. -cpuprofile=cpu.out`

`curl -o cpu.prof http://localhost:6060/debug/pprof/profile?seconds=5`

`go tool pprof cpu.prof`
`top15` 
`granularity=lines`
`hide=runtime`

https://stackademic.com/blog/profiling-go-applications-in-the-right-way-with-examples

https://github.com/markfarnan/go-canvas

## Tasks

### General

- [x] Add compilation to WASM
- [x] Add build scripts
- [x] Add build instructions

### Game

- [x] Add generic block
- [x] Add world to store blocks
- [x] Add block update game loop
- [x] Add saving of camera to file [Windows]
- [x] Add block subupdate game loop
- [ ] Add saving of world to file
- [ ] Add variable tick rate
- [x] Add saving of camera to file [WASM]

### Performance
- [x] Investigate how to debug performance in Go
- [x] New game loop minimising go routines and sleep
- [x] Only allocate depth buffer once
- [x] Only allocate image date in JS once
- [x] Use Go's internal nanotime over time.Now()
- [x] Avoid image.Set, goto underlying pixels or use SetRGBA
- [ ] Improve scaling performance
- [ ] Investigate parallelism

### Blocks

- [x] Add redstone torches
- [x] Add redstone block
- [x] Add redstone lamps
- [x] Add solid blocks
- [x] Add levers
- [x] Add weak & strong power
- [ ] Add pistons
- [ ] Add redstone dust
- [ ] Add redstone repeaters 
- [ ] Add multiblock movement
- [ ] Add slimeblocks
- [ ] Add comparators
- [ ] Add observers
- [ ] Add sand

### User Interface

- [x] Add camera movement
- [x] Add camera rotation 
- [x] Add input via terminal [Windows]
- [x] Add statistics (fps, tps ...)
- [x] Add mouse camera movement [WASM]
- [x] Add raycasting for block selection
- [x] Add place, destroy, pick & interact with blocks
- [ ] Add inventory system
- [x] Make block placement directional
- [x] Add ability to toggle control lock
- [x] Show welcome page by default on first page visit

### Rendering

- [x] Add 3D line rendering
- [x] Add wireframe triangle rendering
- [x] Add solid color triangle rendering
- [x] Add text rendering
- [x] Add textured triangle rendering
- [x] Add depth buffer
- [x] Make texture account for perspective
- [ ] Support OBJ format
- [x] Add welcome page
- [x] Support scene scaling
- [x] Fix texture leakage
- [ ] Add model viewer & builder
- [x] Investigate async vs basic gameloop for performance
- [x] Scale text with image size
- [x] Automatically resize on window rescale
- [x] Fix depth / uv coordinates on camera pitching
- [x] Do not render opaque neighbouring faces
- [ ] Cut triangles off outside viewport
- [x] Add crosshair
- [ ] Make crosshair change color for contrast