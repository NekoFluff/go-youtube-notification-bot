package discord

import (
	"fmt"

	"github.com/NekoFluff/gobot/commands"
	"github.com/bwmarrin/discordgo"
)

var AllCommands = []*commands.DiscordCommand{
	commands.Help,
	commands.GoRoutine,
}

func StartBot(Token string) (dg *discordgo.Session, err error) {
	// Create a new Discord session using the provided bot token.
	dg, err = discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
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

	for _, c := range AllCommands {
		if m.Content == c.Command {
			c.Execute(s, m)
		}
	}
}
