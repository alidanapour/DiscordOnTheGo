package messages

import (
	"os"
	"strings"

	"example.com/external"
	rps "example.com/games/RPS"
	ttt "example.com/games/TTT"
	"example.com/games/dad"
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

func MessageFileOutHandler(session *discordgo.Session, channelID string, file string) {
	f, err := os.Open(file)

	if err != nil {
		term.Print(term.ERROR, err.Error())
		return
	}

	_, err = session.ChannelFileSend(channelID, "RPS.png", f)
	if err != nil {
		term.Print(term.ERROR, err.Error())
	}
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

		//TODO: add flag to display or hide all user commands server-side
		term.Print(term.COMMAND, message.Author.Username+": "+newMessage)

		switch command {
		case "!ping":
			MessageOutHandler(session, message.ChannelID, "pong")

		case "!pong":
			MessageOutHandler(session, message.ChannelID, "ping")

		case "!ttt":
			tttResult := ttt.PlayTTT(message.Author.ID, message.Author.Username, args)
			MessageOutHandler(session, message.ChannelID, tttResult)

		case "!rps":
			success, str, err := rps.PlayRPS(message.Author.ID, message.Author.Username, args)
			if err != nil {
				term.Print(term.ERROR, err.Error())
			} else if success {
				MessageFileOutHandler(session, message.ChannelID, str)
			} else {
				MessageOutHandler(session, message.ChannelID, str)
			}

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

		words := strings.Fields(newMessage)
		if len(words) == 2 && (strings.ToUpper(words[0]) == "I'M" || strings.ToUpper(words[0]) == "IM") {
			MessageOutHandler(session, message.ChannelID, dad.HelloDad(words[1]))
		}
	}

}
