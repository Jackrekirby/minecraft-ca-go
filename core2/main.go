package core2

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

type Block int8
type Direction int8
type Power int8

const (
	Air Block = iota
	RedstoneTorch_On_Left
	RedstoneTorch_On_Right
	RedstoneTorch_On_Up
	RedstoneTorch_On_Down
	RedstoneTorch_Off_Left
	RedstoneTorch_Off_Right
	RedstoneTorch_Off_Up
	RedstoneTorch_Off_Down
	RedstoneBlock
)

const (
	Left Direction = iota
	Right
	Up
	Down
)

type Vec2i struct {
	X, Y int
}

func (d Direction) Opposite() Direction {
	switch d {
	case Left:
		return Right
	case Right:
		return Left
	case Up:
		return Down
	case Down:
		return Up
	default:
		panic(fmt.Sprintf("unknown direction: %v", d))
	}
}

func (d Direction) ToVector() Vec2i {
	switch d {
	case Left:
		return Vec2i{-1, 0}
	case Right:
		return Vec2i{1, 0}
	case Up:
		return Vec2i{0, 1}
	case Down:
		return Vec2i{0, -1}
	default:
		panic(fmt.Sprintf("unknown direction: %v", d))
	}
}

const (
	On Power = iota
	Off
)

// Define a map between Block and color
var blockColorMap = map[Block]color.RGBA{
	Air:                     {R: 255, G: 255, B: 255, A: 100},
	RedstoneTorch_On_Left:   {R: 255, G: 0, B: 0, A: 255},
	RedstoneTorch_On_Right:  {R: 255, G: 0, B: 0, A: 255},
	RedstoneTorch_On_Up:     {R: 255, G: 0, B: 0, A: 255},
	RedstoneTorch_On_Down:   {R: 255, G: 0, B: 0, A: 255},
	RedstoneTorch_Off_Left:  {R: 100, G: 0, B: 0, A: 255},
	RedstoneTorch_Off_Right: {R: 100, G: 0, B: 0, A: 255},
	RedstoneTorch_Off_Up:    {R: 100, G: 0, B: 0, A: 255},
	RedstoneTorch_Off_Down:  {R: 100, G: 0, B: 0, A: 255},
	RedstoneBlock:           {R: 0, G: 255, B: 0, A: 255},
}

type State struct {
	Width     int
	Height    int
	Blocks0   []Block
	Blocks1   []Block
	Iteration int
}

// Create an image from the blocks
func (s *State) WriteBlocksToImage(filename string, img *image.RGBA) error {
	blocks := s.GetCurrentIterationBlocks()

	// Map the blocks to colors
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			block := s.GetBlock(x, s.Height-1-y, blocks) // Flip vertically for correct orientation
			clr, exists := blockColorMap[block]
			if !exists {
				clr = color.RGBA{0, 0, 0, 255} // Default to black if not mapped
			}
			img.SetRGBA(x, y, clr)
		}
	}

	// Create the output file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the image as PNG
	return png.Encode(file, img)
}

func (s *State) Initialise(width, height int) {
	n := width * height
	s.Width = width
	s.Height = height
	s.Blocks0 = make([]Block, n)
	s.Blocks1 = make([]Block, n)
	s.Iteration = 0
}

func (s *State) GetCurrentIterationBlocks() []Block {
	if s.Iteration%2 == 0 {
		return s.Blocks0
	}
	return s.Blocks1
}

func (s *State) GetNextIterationBlocks() []Block {
	if s.Iteration%2 == 0 {
		return s.Blocks1
	}
	return s.Blocks0
}

func (s *State) PrintBlocks() {
	fmt.Println(s.Iteration)
	b := s.GetCurrentIterationBlocks()
	for i := s.Height - 1; i >= 0; i-- {
		start := i * s.Width
		end := start + s.Width
		fmt.Println(b[start:end])
	}
}

func (s *State) GetIndex(x, y int) int {
	return x + y*s.Width
}

func (s *State) SetBlock(x, y int, b Block, blocks []Block) {
	blocks[s.GetIndex(x, y)] = b
}

func (s *State) GetBlock(x, y int, blocks []Block) Block {
	if x < 0 || x >= s.Width || y < 0 || y >= s.Height {
		return Air
	}
	return blocks[s.GetIndex(x, y)]
}

func (b Block) isRedstoneTorch() bool {
	return b == RedstoneTorch_On_Left ||
		b == RedstoneTorch_On_Right ||
		b == RedstoneTorch_On_Up ||
		b == RedstoneTorch_On_Down ||
		b == RedstoneTorch_Off_Left ||
		b == RedstoneTorch_Off_Right ||
		b == RedstoneTorch_Off_Up ||
		b == RedstoneTorch_Off_Down
}

func (b Block) getPower() Power {
	switch b {
	case RedstoneTorch_On_Left, RedstoneTorch_On_Right, RedstoneTorch_On_Up, RedstoneTorch_On_Down:
		return On
	case RedstoneTorch_Off_Left, RedstoneTorch_Off_Right, RedstoneTorch_Off_Up, RedstoneTorch_Off_Down:
		return Off
	default:
		panic(fmt.Sprintf("unknown block: %v", b))
	}
}

func (b Block) getDirection() Direction {
	switch b {
	case RedstoneTorch_On_Left, RedstoneTorch_Off_Left:
		return Left
	case RedstoneTorch_On_Right, RedstoneTorch_Off_Right:
		return Right
	case RedstoneTorch_On_Up, RedstoneTorch_Off_Up:
		return Up
	case RedstoneTorch_On_Down, RedstoneTorch_Off_Down:
		return Down
	default:
		panic(fmt.Sprintf("unknown block: %v", b))
	}
}

func (b Block) getPowerInDirection(direction Direction) Power {
	if b == RedstoneBlock {
		return On
	} else if b.isRedstoneTorch() && b.getPower() == On && b.getDirection() != direction.Opposite() {
		return On
	} else {
		return Off
	}
}

func BuildRedstoneTorch(power Power, direction Direction) Block {
	switch power {
	case On:
		switch direction {
		case Left:
			return RedstoneTorch_On_Left
		case Right:
			return RedstoneTorch_On_Right
		case Up:
			return RedstoneTorch_On_Up
		case Down:
			return RedstoneTorch_On_Down
		default:
			panic(fmt.Sprintf("unknown direction: %v", direction))
		}
	case Off:
		switch direction {
		case Left:
			return RedstoneTorch_Off_Left
		case Right:
			return RedstoneTorch_Off_Right
		case Up:
			return RedstoneTorch_Off_Up
		case Down:
			return RedstoneTorch_Off_Down
		default:
			panic(fmt.Sprintf("unknown direction: %v", direction))
		}
	default:
		panic(fmt.Sprintf("unknown power state: %v", power))
	}
}

func (s *State) UpdateRedstoneTorch(self Block, x, y int, blocks []Block) (Block, bool) {
	d0 := self.getDirection()
	v0 := d0.ToVector()
	neighbour := s.GetBlock(x-v0.X, y-v0.Y, blocks)

	p0 := self.getPower()
	p1 := neighbour.getPowerInDirection(d0)
	if p0 == Off && p1 == Off {
		return BuildRedstoneTorch(On, d0), true
	} else if p0 == On && p1 == On {
		return BuildRedstoneTorch(Off, d0), true
	}
	return self, false
}

func (s *State) UpdateBlock(x, y int, blocks []Block) (Block, bool) {
	self := s.GetBlock(x, y, blocks)

	if self.isRedstoneTorch() {
		return s.UpdateRedstoneTorch(self, x, y, blocks)
	}

	return self, false
}

func (s *State) UpdateBlocks() (updateCount int) {
	blocks := s.GetCurrentIterationBlocks()
	updatedBlocks := s.GetNextIterationBlocks()
	updateCount = 0
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			newBlock, updated := s.UpdateBlock(x, y, blocks)
			s.SetBlock(x, y, newBlock, updatedBlocks)
			if updated {
				updateCount++
			}
		}
	}
	s.Iteration++
	return updateCount
}

func (s *State) RunIterations(callback func(s *State)) {
	callback(s)
	for s.Iteration < 200 {
		updateCount := s.UpdateBlocks()
		callback(s)
		if updateCount == 0 {
			break
		}
	}
}

func RunProgram() {
	state := State{}
	state.Initialise(8, 8)

	// define blocks
	b := state.GetCurrentIterationBlocks()

	state.SetBlock(2, 2, RedstoneTorch_On_Left, b)
	state.SetBlock(3, 2, RedstoneTorch_Off_Left, b)
	state.SetBlock(4, 2, RedstoneTorch_On_Left, b)

	state.SetBlock(5, 4, RedstoneTorch_Off_Down, b)
	state.SetBlock(5, 3, RedstoneTorch_On_Down, b)
	state.SetBlock(5, 2, RedstoneTorch_Off_Down, b)

	state.SetBlock(2, 5, RedstoneTorch_Off_Right, b)
	state.SetBlock(3, 5, RedstoneTorch_On_Right, b)
	state.SetBlock(4, 5, RedstoneTorch_Off_Right, b)
	state.SetBlock(5, 5, RedstoneTorch_On_Right, b)

	state.SetBlock(1, 5, RedstoneTorch_On_Up, b)
	state.SetBlock(1, 4, RedstoneTorch_Off_Up, b)
	state.SetBlock(1, 3, RedstoneTorch_On_Up, b)
	state.SetBlock(1, 2, RedstoneTorch_On_Left, b)

	imgPath := "output/gpu.png"
	img := image.NewRGBA(image.Rect(0, 0, state.Width, state.Height))

	var callback func(s *State)
	if PRINT_BLOCKS == 0 {
		callback = func(s *State) { s.PrintBlocks() }
	} else if PRINT_BLOCKS == 1 {
		callback = func(s *State) {
			s.WriteBlocksToImage(imgPath, img)
			time.Sleep(100 * time.Millisecond)

			if s.Iteration == 4 {
				state.SetBlock(1, 1, RedstoneBlock, b)
			}
		}
	} else {
		callback = func(s *State) {}
	}

	state.RunIterations(callback)
}

func ProfileCPU() func() {
	cpuProfile, err := os.Create("cpu.pprof")
	if err != nil {
		panic(fmt.Sprint("Could not create CPU profile:", err))
	}

	err = pprof.StartCPUProfile(cpuProfile)
	if err != nil {
		panic(fmt.Sprint("Could not start CPU profile:", err))
	}

	return func() {
		defer cpuProfile.Close()
		defer pprof.StopCPUProfile()
	}
}

func ProfileMemory() func() {
	memProfile, err := os.Create("mem.pprof")
	if err != nil {
		panic(fmt.Sprint("Could not create memory profile:", err))
	}

	// Force garbage collection to capture up-to-date memory allocation information
	runtime.GC()

	err = pprof.WriteHeapProfile(memProfile)
	if err != nil {
		panic(fmt.Sprint("Could not write memory profile:", err))
	}
	return func() {
		defer memProfile.Close()
	}
}

const PROFILE = false
const PRINT_BLOCKS = 1

func Main() {
	fmt.Println("core2.Main()")

	if PROFILE {
		defer ProfileCPU()()
	}

	RunProgram()

	if PROFILE {
		defer ProfileMemory()
	}

}
