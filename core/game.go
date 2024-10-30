package core

type GameState int

const (
	Quit GameState = iota
	Pause
	Play
)
