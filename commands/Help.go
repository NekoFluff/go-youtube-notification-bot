package commands

import (
	"github.com/bwmarrin/discordgo"
)

var Help = &DiscordCommand{Command: "!help", Execute: func(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Help is on it's way! Hooray!")
}}

// func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	if m.Content == "!help" {
// 		mux.Greet()
// 		s.ChannelMessageSend(m.ChannelID, "Help is on it's way!")
// 	}
// }
