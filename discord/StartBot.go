package discord

import (
	"fmt"

	"github.com/NekoFluff/gobot/commands"
	"github.com/bwmarrin/discordgo"
)

func StartBot(Token string) (s *discordgo.Session, err error) {
	// Create a new Discord session using the provided bot token.
	s, err = discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	s.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	s.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Help command has to be separated from all commands to prevent cyclic behavior
	if m.Content == commands.Help.Command {
		commands.Help.Execute(s, m)
	}

	for _, c := range commands.AllCommands {
		if m.Content == c.Command {
			c.Execute(s, m)
		}
	}
}
