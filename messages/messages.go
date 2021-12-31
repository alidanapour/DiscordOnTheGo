package messages

import (
	"strings"

	"example.com/external"
	"example.com/games"
	term "example.com/terminal"
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

func MessageEmbedOutHandler(session *discordgo.Session, channelID string, embed *discordgo.MessageEmbed) {
	session.ChannelMessageSendEmbed(channelID, embed)
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
		term.Print(term.COMMAND, message.Author.Username+": "+command)
		switch command {
		case "!ping":
			MessageOutHandler(session, message.ChannelID, "")
		case "!pong":
			MessageOutHandler(session, message.ChannelID, "ping")
		case "!ttt":
			tttResult := games.PlayTTT(message.Author.ID, message.Author.Username, args)
			MessageOutHandler(session, message.ChannelID, tttResult)
		case "!apod":
			apodResult, err := external.ApodRequest()
			if err == nil {
				MessageEmbedOutHandler(session, message.ChannelID, apodResult)
			} else {
				MessageOutHandler(session, message.ChannelID, err.Error())
			}
		default:
			term.Print(term.ERROR, "Unknown Command")
		}

	} else {
		term.Print(term.MESSAGE, message.Author.Username+": "+newMessage)
	}

}
