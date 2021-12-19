package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var Help = &DiscordCommand{
	Command:     "!help",
	Description: "Display all available commands!",
	Execute: func(s *discordgo.Session, m *discordgo.MessageCreate) {
		embedFields := []*discordgo.MessageEmbedField{}

		// Setup MessageEmbedField
		embedFields = append(embedFields, &discordgo.MessageEmbedField{
			Name:   "Setup",
			Value:  "Just make text channels named `#gobot` and `#gobot-live`. All push notifications from Youtube will be sent here! (e.g. When a video is uploaded/updated)",
			Inline: false,
		})

		// Build all the commands into MessageEmbedFields
		for _, c := range AllCommands {
			embedField := &discordgo.MessageEmbedField{
				Name:   c.Command,
				Value:  c.Description,
				Inline: false,
			}
			embedFields = append(embedFields, embedField)
		}

		// Build the embed
		embed := &discordgo.MessageEmbed{
			Type:   discordgo.EmbedTypeRich,
			Title:  "Help Page",
			Fields: embedFields,
		}

		_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
		log.Println(err)
	}}
