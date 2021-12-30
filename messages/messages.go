package messages

import (
	"fmt"
	"strings"

	"example.com/games"
	"github.com/bwmarrin/discordgo"
)

func MessageInHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore bot messages
	if message.Author.ID == session.State.User.ID {
		return
	} else {
		newMessage(session, message)
	}
}

func MessageOutHandler(session *discordgo.Session, channelID string, text string) {
	session.ChannelMessageSend(channelID, text)
}

// This function handles all incoming messages
func newMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	newMessage := message.Content

	if strings.HasPrefix(newMessage, "!") {
		// Get the command and arguments
		args := strings.Fields(message.Content)
		command := args[0]
		args = append(args[1:]) // remove command from args

		//TODO: add flag to display or hide all user commands
		fmt.Println("Command " + command + " sent by " + message.Author.Username)

		switch command {
		case "!ping":
			MessageOutHandler(session, message.ChannelID, "pong")
		case "!pong":
			MessageOutHandler(session, message.ChannelID, "ping")
		case "!TTT":
			result := games.PlayTTT(message.Author.ID, message.Author.Username, args)
			MessageOutHandler(session, message.ChannelID, result)
		default:
			fmt.Println("Unknown command")
		}

	} else {
		fmt.Println("Message: " + newMessage)
	}

}
