package games

import (
	"errors"
	"fmt"
	"strconv"
)

// By default player1 = X and player2 = O
type TTTGame struct {
	player1     string
	player1ID   string
	player2     string
	player2ID   string
	currentTurn int //P1 = 1, P2 = 2
	board       [3][3]rune
}

var (
	Game TTTGame
)

func init() {
	resetGame()
}

func resetGame() {
	Game.player1 = ""
	Game.player1ID = ""
	Game.player2 = ""
	Game.player2ID = ""
	Game.currentTurn = 1
	Game.board = [3][3]rune{{' ', ' ', ' '}, {' ', ' ', ' '}, {' ', ' ', ' '}}
}

// Returns the board as a string to print in a message. Special characters are added
// to the ends to show message as monospace on discord
func getBoard() string {
	var result string = "`TICTACTOE\n"

	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			result += "[" + string(Game.board[row][col]) + "]"
		}
		result += "\n"
	}
	result += "`"
	return result
}

// Returns if the current player has won.
func checkGameWin() bool {
	var markToCheck rune

	if Game.currentTurn == 1 {
		markToCheck = 'X'
	} else {
		markToCheck = 'O'
	}

	// Check rows
	for row := 0; row < 2; row++ {
		if Game.board[row][0] == markToCheck &&
			Game.board[row][1] == markToCheck &&
			Game.board[row][2] == markToCheck {
			return true
		}
	}

	// Check columns
	for col := 0; col < 2; col++ {
		if Game.board[0][col] == markToCheck &&
			Game.board[1][col] == markToCheck &&
			Game.board[2][col] == markToCheck {
			return true
		}
	}

	// Check diagonals. Hard coded because tictactoe shouldnt change.
	if Game.board[0][0] == markToCheck &&
		Game.board[1][1] == markToCheck &&
		Game.board[2][2] == markToCheck {
		return true
	}
	if Game.board[2][2] == markToCheck &&
		Game.board[1][1] == markToCheck &&
		Game.board[0][0] == markToCheck {
		return true
	}

	return false
}

// Executes a players move and returns and error if unsuccessful.
func playerMove(playerID string, row int, col int) error {
	// Attempt to make a move, return error otherwise
	if Game.board[row-1][col-1] != ' ' {
		return errors.New("TTT Game Error: row " + fmt.Sprint(row) +
			" and col " + fmt.Sprint(col) + " is already taken")
	}

	var playerMark rune
	if Game.currentTurn == 1 {
		playerMark = 'X'
	} else {
		playerMark = 'O'
	}

	Game.board[row-1][col-1] = playerMark
	return nil
}

// Attempts to get two int values from the arguments. If either fail an error is returned.
func getMoveArgs(args []string) (int, int, error) {
	row, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, 0, errors.New("TTT Game error: unknown first argument")
	}

	col, err := strconv.Atoi(args[1])
	if err != nil {
		return 0, 0, errors.New("TTT Game error: unknown second argument")
	}

	// Check that row and col commands are in range
	if row < 1 || row > 3 {
		return 0, 0, errors.New("TTT Game Error: row value out of range [1,3]")
	} else if col < 1 || col > 3 {
		return 0, 0, errors.New("TTT Game Error: col value out of range [1,3]")
	}

	return row, col, nil
}

// Returns true if person who issued the command is playing and in turn; false otherwise.
func validatePlayer(playerID string) error {
	if playerID != Game.player1ID && playerID != Game.player2ID {
		return errors.New("TTT Game Error: non-player issued command")
	}

	if Game.currentTurn == 1 && playerID == Game.player2ID {
		return errors.New("TTT Game Error: Player 2 out of turn")
	} else if Game.currentTurn == 2 && playerID == Game.player1ID {
		return errors.New("TTT Game Error: Player 1 out of turn")
	}

	return nil
}

// Main function that runs tic tac toe. Returns a string representing the game to print.
// Representation could be a message, board, or gameplay error.
func PlayTTT(playerID string, player string, args []string) string {
	var result string = ""

	// Check if args were added
	var command string = ""
	if len(args) > 0 {
		command = args[0]
	}

	switch command {
	// New game or player
	case "":
		if Game.player1ID == "" {
			Game.player1ID = playerID
			Game.player1 = player
			return player +
				" has entered the arena for some tic tac toe, who will join them? " +
				"Type !TTT to join."
		} else if Game.player2ID == "" {
			if playerID == Game.player1ID {
				return "Error: must have different players"
			}
			Game.player2ID = playerID
			Game.player2 = player
			return player + " has accepted the challenge. Player 1's turn..\n" + getBoard()
		} else {
			return "Please enter a command, try '!ttt help'"
		}

	// Help menu
	case "help":
		return "To play tic tac toe, type !TTT. After starting a new game type " +
			"!TTT followed by any of the following commands:\n" +
			"`concede`: To concede the game and let your opponent win.\n" +
			"`# #    `: Two numbers corresponding with the row and column you wish to mark next. " +
			"for example: '1 1' is top left, '2 3' is middle right, etc.\n"

	// Concede
	case "concede":
		return "TODO..."

	// Otherwise assume valid move (i.e two numbers).
	default:
		// Check for valid command/args
		row, col, err := getMoveArgs(args)
		if err != nil {
			return err.Error()
		}

		// Validate player and make next move
		err = validatePlayer(playerID)
		if err != nil {
			return err.Error()
		}

		// execute move (if valid)
		err = playerMove(playerID, row, col)
		if err != nil {
			return err.Error()
		} else {
			result += getBoard()
		}

		if checkGameWin() {
			result += (player + " has won!! Congratulations.\n")
			resetGame()
		} else {
			result += (player + " has marked (" + fmt.Sprint(row) + "," + fmt.Sprint(col) + ")!\n")
			if Game.currentTurn == 1 {
				Game.currentTurn = 2
			} else {
				Game.currentTurn = 1
			}
		}
	}

	// Return message to discord
	return result
}
