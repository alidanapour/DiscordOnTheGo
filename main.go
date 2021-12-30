// Some of the code derived from ping pong example
// -> https://github.com/bwmarrin/discordgo/blob/master/examples/pingpong/main.go

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"example.com/messages"
	"example.com/terminal"
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
	fmt.Println(terminal.Green + "Starting up Discord bot..." + terminal.Reset)
	session, err := discordgo.New("Bot " + Token)

	if err != nil {
		fmt.Println(terminal.Red + "Error starting a Discord session" + terminal.Reset)
		log.Fatal(err)
		return
	}

	if Token == "" {
		fmt.Println(terminal.Red + "Error: no token provided.. exiting" + terminal.Reset)
		return
	}

	//Register messageCreate func as a callback for MessageCreate events
	//session.AddHandler(messageCreate)
	session.AddHandler(messages.MessageInHandler)

	// We only care about receiving messages
	session.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket to discord and listen
	err = session.Open()
	if err != nil {
		fmt.Println(terminal.Red+"Error opening websocket: "+terminal.Reset, err)
		return
	}

	// Main loop
	fmt.Println(terminal.Green + "Bot is now running..." + terminal.Reset)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up and exit
	session.Close()
}
