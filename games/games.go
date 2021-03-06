package games

import (
	rps "example.com/games/RPS"
	ttt "example.com/games/TTT"
)

// Returns a string generated by tictactoe. This can be anything from
// gameplay to gameplay logic errors
func PlayTTT(playerID string, player string, args []string) string {
	return ttt.PlayTTT(player, player, args)
}

/*
Returns the following
Bool  : Gameplay was successful and image was generated
String: Message from game and bool = false
Error : Error from game. Incorrect args, failure to generate image, etc
*/
func PlayRPS(playerID string, player string, args []string) (bool, string, error) {
	success, str, err := rps.PlayRPS(playerID, player, args)

	if err != nil {
		return false, "", err
	}
	return success, str, nil

}
