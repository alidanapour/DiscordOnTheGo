package state

// Enums for rps
type Move int

const (
	ROCK Move = iota
	PAPER
	SCISSORS
)

type Result int

const (
	TIE Result = iota
	LOSE
	WIN
)
