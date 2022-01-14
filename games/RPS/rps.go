package rps

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"example.com/games/RPS/draw"
	st "example.com/games/RPS/state"
)

// Vars
var botMove st.Move

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

// Returns a string denoting if a player has won RPS against the bot
// TODO: track wins/losses with playerID
func PlayRPS(playerID string, player string, args []string) (bool, string, error) {

	if len(args) == 0 {
		return false, "To play Rock Paper Scissors type !rps <move>. Where <move>" +
			"is replaced with Rock, Paper, or Scissors.", nil
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
