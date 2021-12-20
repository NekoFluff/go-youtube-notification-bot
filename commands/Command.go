package commands

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

type DiscordCommand struct {
	Command     string
	Description string
	Execute     func(s *discordgo.Session, m *discordgo.MessageCreate)
}

func NewDiscordCommand(command string, description string, execute func(s *discordgo.Session, m *discordgo.MessageCreate)) *DiscordCommand {
	prefix := os.Getenv("COMMAND_PREFIX")
	if prefix == "" {
		prefix = "!"
	}

	return &DiscordCommand{
		Command:     prefix + command,
		Description: description,
		Execute:     execute,
	}
}
