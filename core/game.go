package core

type GameState int

const (
	Quit GameState = iota
	Paused
	Playing
	Pausing
)

// String method returns the string representation of the GameState
func (g GameState) String() string {
	switch g {
	case Quit:
		return "Quit"
	case Paused:
		return "Paused"
	case Pausing:
		return "Pausing"
	case Playing:
		return "Playing"
	default:
		panic("GameState not implemented")
	}
}
