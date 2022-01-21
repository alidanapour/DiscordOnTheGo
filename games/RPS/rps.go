package rps

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"example.com/games/RPS/draw"
	st "example.com/games/RPS/state"
)

// Vars and Constants
var botMove st.Move

const help = "Rock Paper Scissors - List of Commands: \n" +
	"`!rps <move>` - Play a game of RPS with the bot. Replace <move> with a valid move.\n" +
	"`!rps set-bg # # #` - Sets the background of the RPS game to the desired RGB values.\n" +
	"`!rps set-circle # # #` - Sets the circles to the desired RGB values."

func init() {
	rand.Seed(time.Now().UnixNano())
	nextMove()
}

func nextMove() {
	botMove = st.Move(rand.Intn(3))
}

func parsePlayerMove(playerMove string) (st.Move, error) {
	move := strings.ToUpper(playerMove)
	switch move {
	case "ROCK":
		return st.ROCK, nil
	case "PAPER":
		return st.PAPER, nil
	case "SCISSORS":
		return st.SCISSORS, nil
	default:
		return 0, errors.New("RPS Game Error: Invalid move specified")
	}
}

// Executes player command.
// May play a game or may change colors
// Returns true if a game was played (and image dir returned)
// Returns false if a message is needed to send to the user
func PlayRPS(playerID string, player string, args []string) (bool, string, error) {

	if len(args) == 0 {
		return false, help, nil
	} else if args[0] == "help" {
		return false, help, nil
	} else if args[0] == "set-bg" || args[0] == "set-circle" || args[0] == "set-default" {
		msg, err := draw.RPS_SetColor(args)
		if err != nil {
			return false, "", err
		}
		return false, msg, nil
	}

	playerMove, err := parsePlayerMove(args[0])
	if err != nil {
		return false, err.Error(), nil
	}

	var gameResult st.Result

	// First check for tie
	if playerMove == botMove {
		gameResult += st.TIE
	}

	// Cases are based on player move
	switch playerMove {
	case st.ROCK:
		if botMove == st.PAPER {
			gameResult += st.LOSE
		} else if botMove == st.SCISSORS {
			gameResult += st.WIN
		}
	case st.PAPER:
		if botMove == st.SCISSORS {
			gameResult += st.LOSE
		} else if botMove == st.ROCK {
			gameResult += st.WIN
		}
	case st.SCISSORS:
		if botMove == st.ROCK {
			gameResult += st.LOSE
		} else if botMove == st.PAPER {
			gameResult += st.WIN
		}
	}

	res, err := draw.RPS_GenerateImage(player, playerMove, botMove, gameResult)
	if err != nil {
		return false, "", err
	}

	nextMove()
	return true, res, nil
}
