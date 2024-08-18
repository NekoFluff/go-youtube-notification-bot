package commands

import (
	"fmt"

	"log/slog"

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

			user := i.Interaction.User
			if i.Interaction.Member != nil {
				user = i.Interaction.Member.User
			}

			if add := optionMap["add"]; add != nil {
				creator := add.Options[0].StringValue()

				subscription := data.Subscription{
					User:         user.ID,
					Subscription: creator,
				}
				data.SaveSubscription(subscription)

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Subscribed to `%v`", creator),
					},
				})
				if err != nil {
					slog.Error("Failed to respond to interaction", "error", err)
				}
			} else if remove := optionMap["remove"]; remove != nil {
				creator := remove.Options[0].StringValue()

				subscription := data.Subscription{
					User:         user.ID,
					Subscription: creator,
				}
				data.DeleteSubscription(subscription)

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Removed `%s` subscriptions", creator),
					},
				})
				if err != nil {
					slog.Error("Failed to respond to interaction", "error", err)
				}
			} else if list := optionMap["list"]; list != nil {
				subscriptions, err := data.GetSubscriptionsForUser(user.ID)

				if err != nil {
					slog.Error("Failed to get subscription for the user", "error", err, "user", user.ID)

					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Failed to get your subscriptions",
						},
					})
					if err != nil {
						slog.Error("Failed to respond to interaction", "error", err)
					}
					return
				}

				// Check if we have any subscriptions
				if (len(subscriptions)) == 0 {
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "You have no subscriptions",
						},
					})
					if err != nil {
						slog.Error("Failed to respond to interaction", "error", err)
					}
					return
				}

				// Get all creators
				creators := make([]string, len(subscriptions))
				for i, sub := range subscriptions {
					creators[i] = sub.Subscription
				}

				// Join creators with a comma
				creatorsString := fmt.Sprintf("Your subscriptions: `%s`", creators[0])
				for _, creator := range creators[1:] {
					creatorsString += fmt.Sprintf(", `%s`", creator)
				}

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: creatorsString,
					},
				})
				if err != nil {
					slog.Error("Failed to respond to interaction", "error", err)
				}
			} else {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Unknown subcommand provided",
					},
				})
				if err != nil {
					slog.Error("Failed to respond to interaction", "error", err)
				}
			}
		},
	}
}
