package core2

// func BenchmarkRunIterations(b *testing.B) {
// 	state := State{}
// 	state.Initialise(3, 8)
// 	blocks := state.GetCurrentIterationBlocks()
// 	state.SetBlock(1, 5, RedstoneTorch_On, blocks)
// 	state.SetBlock(1, 4, RedstoneTorch_Off, blocks)
// 	state.SetBlock(1, 3, RedstoneTorch_On, blocks)
// 	state.SetBlock(1, 2, RedstoneTorch_Off, blocks)
// 	// state.SetBlock(1, 1, RedstoneBlock, blocks)
// 	output := func(s *State) {}
// 	for i := 0; i < b.N; i++ {
// 		state.RunIterations(output)
// 	}
// }
