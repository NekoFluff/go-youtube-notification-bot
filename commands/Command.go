package commands

import (
	"github.com/NekoFluff/gobot/utils"
	"github.com/bwmarrin/discordgo"
)

type DiscordCommand struct {
	Command     string
	Description string
	Execute     func(s *discordgo.Session, m *discordgo.MessageCreate)
}

func NewDiscordCommand(command string, description string, execute func(s *discordgo.Session, m *discordgo.MessageCreate)) *DiscordCommand {
	prefix := utils.GetEnvVar("COMMAND_PREFIX")

	return &DiscordCommand{
		Command:     prefix + command,
		Description: description,
		Execute:     execute,
	}
}
