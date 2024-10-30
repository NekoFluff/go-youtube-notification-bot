package commands

import (
	"log/slog"

	"github.com/NekoFluff/discord"
	"github.com/bwmarrin/discordgo"
)

func Help() discord.Command {
	command := "help"

	return discord.Command{
		Command: discordgo.ApplicationCommand{
			Name:        command,
			Description: "A guide on how to use the bot.",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			msg := "**Setup:**"
			msg += "\n1. Create a `#hololive-notifications` channel. This is where the bot will send notifications when a new livestream video is published/updated. The bot will also send a 15 minute reminder before any stream starts."
			msg += "\n2. Create a `#hololive-livestream-started` channel. This is where the bot will send notifications when a Hololive member goes live."
			msg += "\n3. (optional) You can use the `/subscription add` command to get DMs when a Hololive member goes live."

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
				},
			})
			if err != nil {
				slog.Error("An error occurred while helping the server", "error", err)
			}
		},
	}
}
