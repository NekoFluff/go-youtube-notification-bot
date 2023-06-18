package commands

import (
	"log"

	"github.com/NekoFluff/discord"
	"github.com/bwmarrin/discordgo"
)

func Ping() discord.Command {
	command := "ping"

	return discord.Command{
		Command: discordgo.ApplicationCommand{
			Name:        command,
			Description: "Is the bot online?",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := respondToInteraction(s, i.Interaction, "Pong!")
			if err != nil {
				log.Println("An error occurred while pinging the server")
				log.Println(err)
			}
		},
	}
}
