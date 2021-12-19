package discord

import "github.com/bwmarrin/discordgo"

func StopBot(s *discordgo.Session) {
	// Cleanly close down the Discord session.
	s.Close()
}
