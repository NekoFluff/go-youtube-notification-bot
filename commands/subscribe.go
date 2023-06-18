package commands

import (
	"fmt"
	"log"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/go-hololive-notification-bot/data"
	"github.com/bwmarrin/discordgo"
)

func Subscribe() discord.Command {
	command := "subscribe"

	return discord.Command{
		Command: discordgo.ApplicationCommand{
			Name:        command,
			Description: "Subscribe to get notifications from certain vtubers!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "creator",
					Description: "The creator (e.g. gura)",
					Required:    true,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			creator := optionMap["creator"].StringValue()

			subscription := data.Subscription{
				User:         i.Message.Author.ID,
				Subscription: creator,
			}
			data.SaveSubscription(subscription)

			_, err := s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Subscribed to %v", creator))
			if err != nil {
				log.Println(err)
			}
		},
	}
}
