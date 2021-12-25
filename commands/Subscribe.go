package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/NekoFluff/gobot/data"
	"github.com/bwmarrin/discordgo"
)

var Subscribe = NewDiscordCommand(
	"subscribe",
	"Subscribe to get notifications from certain vtubers!",
	func(s *discordgo.Session, m *discordgo.MessageCreate) {
		args := strings.Split(strings.ToLower(m.Content), " ")[1:]
		if len(args) == 0 {
			_, err := s.ChannelMessageSend(m.ChannelID, "No arguments provided")
			if err != nil {
				log.Println(err)
			}
			return
		}

		for _, arg := range args {
			subscription := data.Subscription{
				User:         m.Author.ID,
				Subscription: arg,
			}
			data.SaveSubscription(subscription)

			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Subscribed to %v", arg))
			if err != nil {
				log.Println(err)
			}
		}
	},
)
