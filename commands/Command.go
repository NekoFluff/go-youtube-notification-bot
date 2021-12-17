package commands

import "github.com/bwmarrin/discordgo"

type DiscordCommand struct {
	Command string
	Execute func(s *discordgo.Session, m *discordgo.MessageCreate)
}
