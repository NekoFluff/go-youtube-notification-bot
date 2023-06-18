package commands

import (
	"fmt"
	"log"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/go-hololive-notification-bot/data"
	"github.com/bwmarrin/discordgo"
)

func Subscription() discord.Command {
	command := "subscription"

	creatorCommandOption := []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "creator",
			Description: "The creator (e.g. gura)",
			Required:    true,
		},
	}

	return discord.Command{
		Command: discordgo.ApplicationCommand{
			Name:        command,
			Description: "Subscription related commands",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "add",
					Description: "Subscribe to get DMs when a livestream goes live!",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options:     creatorCommandOption,
				},
				{
					Name:        "remove",
					Description: "Unsubscribe to not longer get DMs when a livestream goes live.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options:     creatorCommandOption,
				},
				{
					Name:        "list",
					Description: "List all our subscriptions.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if add := optionMap["add"]; add != nil {
				creator := add.Options[0].StringValue()

				subscription := data.Subscription{
					User:         i.Interaction.Member.User.ID,
					Subscription: creator,
				}
				data.SaveSubscription(subscription)

				_, err := s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Subscribed to %v", creator))
				if err != nil {
					log.Println(err)
				}
			} else if remove := optionMap["remove"]; remove != nil {
				creator := remove.Options[0].StringValue()

				subscription := data.Subscription{
					User:         i.Interaction.Member.User.ID,
					Subscription: creator,
				}
				data.DeleteSubscription(subscription)

				_, err := s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Removed subscriptions from %s", creator))
				if err != nil {
					log.Println(err)
				}
			} else if list := optionMap["list"]; list != nil {
				_, err := s.ChannelMessageSend(i.ChannelID, "Not implemented yet")
				if err != nil {
					log.Println(err)
				}
			} else {
				_, err := s.ChannelMessageSend(i.ChannelID, "Unknown subcommand provided")
				if err != nil {
					log.Println(err)
				}
			}
		},
	}
}
