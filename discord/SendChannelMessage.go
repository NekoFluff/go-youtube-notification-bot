package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func SendChannelMessage(s *discordgo.Session, channelName string, message string) {
	for _, guild := range s.State.Guilds {
		// Get channels for this guild
		channels, _ := s.GuildChannels(guild.ID)

		for _, c := range channels {
			// Check if channel is a guild text channel and not a voice or DM channel
			if c.Type != discordgo.ChannelTypeGuildText {
				continue
			}

			// Check if channel name matches target name
			if c.Name != channelName {
				continue
			}

			// Send text message
			_, err := s.ChannelMessageSend(
				c.ID,
				message,
			)

			if err != nil {
				fmt.Printf("Failed to send a channel message: %s", err)
			}
		}
	}
}
