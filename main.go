// Some of the code derived from ping pong example
// -> https://github.com/bwmarrin/discordgo/blob/master/examples/pingpong/main.go

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"example.com/messages"
	term "example.com/terminal"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line arguments
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	term.Print(term.SUCCESS, "Starting DiscordOnTheGo ...")
	session, err := discordgo.New("Bot " + Token)

	if err != nil {
		term.Print(term.ERROR, "Could not start a Discord Session")
		log.Fatal(err)
		return
	}

	if Token == "" {
		term.Print(term.ERROR, "No token provided")
		return
	}

	//Register messageCreate func as a callback for MessageCreate events
	session.AddHandler(messages.MessageInHandler)

	// We only care about receiving messages
	session.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket to discord and listen
	err = session.Open()
	if err != nil {
		term.Print(term.ERROR, "Could not open WebSocket")
		return
	}

	// Main loop
	term.Print(term.SUCCESS, "Bot is connected and running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up and exit
	session.Close()
}
