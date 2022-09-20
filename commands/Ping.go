package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var Ping = &DiscordCommand{Command: "!ping", Execute: func(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
	if err != nil {
		fmt.Println(err)
	}
}}
