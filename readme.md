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

`go test -bench=. -cpuprofile=cpu.out`

`curl -o cpu.prof http://localhost:6060/debug/pprof/profile?seconds=5`

`go tool pprof cpu.prof`
`top15` 
`granularity=lines`
`hide=runtime`

https://stackademic.com/blog/profiling-go-applications-in-the-right-way-with-examples

## Tasks

### General

- [x] Add compilation to WASM
- [ ] Add build scripts
- [ ] Add build instructions
- [ ] Add task list

### Game

- [x] Add generic block
- [x] Add world to store blocks
- [x] Add block update game loop
- [x] Add saving of camera to file [Windows]
- [x] Add block subupdate game loop
- [ ] Add saving of world to file
- [ ] Add variable tick rate

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
- [ ] Add mouse camera movement [WASM]
- [ ] Add place & destroy blocks
- [ ] Add inventory system

### Rendering

- [x] Add 3D line rendering
- [x] Add wireframe triangle rendering
- [x] Add solid color triangle rendering
- [x] Add text rendering
- [x] Add textured triangle rendering
- [x] Add depth buffer
- [x] Make texture account for perspective
- [ ] Investigate rendering performance
- [ ] Improve rendering performance
- [ ] Support OBJ format
- [ ] Add welcome page
- [ ] Support scene scaling
- [ ] Fix texture leakage
