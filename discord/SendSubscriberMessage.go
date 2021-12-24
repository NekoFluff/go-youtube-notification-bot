package discord

import (
	"fmt"

	"github.com/NekoFluff/gobot/data"
	"github.com/bwmarrin/discordgo"
)

func SendSubscriberMessage(s *discordgo.Session, authors []string, message string) {
	subscriptions, err := data.GetSubscriptions(authors)

	if err != nil {
		SendDeveloperMessage(s, fmt.Sprintf("Unable to retrieve subscribers: %s", err))
		return
	}

	for _, subscription := range subscriptions {
		ch, err := s.UserChannelCreate(subscription.User)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = s.ChannelMessageSend(ch.ID, message)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
