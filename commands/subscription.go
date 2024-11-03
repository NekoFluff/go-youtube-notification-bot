package commands

import (
	"fmt"

	"log/slog"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/hololive-livestream-notifier-go/data"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Subscription() discord.Command {
	command := "subscription"

	vtuberCommandOption := []*discordgo.ApplicationCommandOption{
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "vtuber",
			Description:  "The vtuber (e.g. Gura)",
			Required:     true,
			Autocomplete: true,
		},
	}

	return discord.Command{
		Command: discordgo.ApplicationCommand{
			Name:        command,
			Description: "Subscription related commands",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "add",
					Description: "Subscribe to get DMs when a livestream goes live! (matches by first or last name)",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options:     vtuberCommandOption,
				},
				{
					Name:        "remove",
					Description: "Unsubscribe to not longer get DMs when a livestream goes live.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options:     vtuberCommandOption,
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
				if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
					handleAutocomplete(add.Options, s, i)
					return
				} else if i.Type == discordgo.InteractionApplicationCommand {
					feedID, _ := primitive.ObjectIDFromHex(add.Options[0].StringValue())

					subscription := data.Subscription{
						User:   user.ID,
						FeedID: feedID,
					}
					data.SaveSubscription(subscription)

					feed, _ := data.GetFeedByID(feedID)

					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("Subscribed to `%v`", feed.FullName()),
						},
					})
					if err != nil {
						slog.Error("Failed to respond to interaction", "error", err)
					}
				}
			} else if remove := optionMap["remove"]; remove != nil {
				if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
					handleAutocomplete(remove.Options, s, i)
					return
				} else if i.Type == discordgo.InteractionApplicationCommand {
					feedID, _ := primitive.ObjectIDFromHex(remove.Options[0].StringValue())

					subscription := data.Subscription{
						User:   user.ID,
						FeedID: feedID,
					}
					data.DeleteSubscription(subscription)

					feed, _ := data.GetFeedByID(feedID)

					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("Removed `%s` subscription", feed.FullName()),
						},
					})
					if err != nil {
						slog.Error("Failed to respond to interaction", "error", err)
					}
				}
			} else if list := optionMap["list"]; list != nil {
				feeds, err := data.GetFeedsForUserBySubscription(user.ID)

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
				if (len(feeds)) == 0 {
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

				// Get all vtubers
				vtubers := make([]string, len(feeds))
				for i, feed := range feeds {
					vtubers[i] = feed.FullName()
				}

				// Join vtubers with a comma
				vtubersString := fmt.Sprintf("Your subscriptions: `%s`", vtubers[0])
				for _, vtuber := range vtubers[1:] {
					vtubersString += fmt.Sprintf(", `%s`", vtuber)
				}

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: vtubersString,
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

func handleAutocomplete(options []*discordgo.ApplicationCommandInteractionDataOption, s *discordgo.Session, i *discordgo.InteractionCreate) {
	vtuber := options[0].StringValue()

	// Get all feeds
	feeds, err := data.GetFeedsByName(vtuber, 20)
	if err != nil {
		slog.Error("Failed to get feeds", "error", err)
		feeds = []data.ChannelFeed{}
	}

	// Create a list of vtubers for the autocomplete
	choices := make([]*discordgo.ApplicationCommandOptionChoice, len(feeds))
	for i, feed := range feeds {
		choices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  feed.FullName(),
			Value: feed.ID,
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})

	if err != nil {
		slog.Error("Failed to respond to interaction", "error", err)
	}
}
