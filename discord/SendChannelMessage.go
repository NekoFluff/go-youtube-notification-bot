package discord

import (
	"github.com/bwmarrin/discordgo"
)

func SendChannelMessage(dg *discordgo.Session, channelName string, message string) {
	for _, guild := range dg.State.Guilds {
		// Get channels for this guild
		channels, _ := dg.GuildChannels(guild.ID)

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
			dg.ChannelMessageSend(
				c.ID,
				message,
			)
		}
	}
}
