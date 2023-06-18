package discord

import (
	"fmt"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/go-hololive-notification-bot/data"
)

func SendSubscriberMessage(bot *discord.Bot, authors []string, message string) {
	subscriptions, err := data.GetSubscriptions(authors)

	if err != nil {
		bot.SendDeveloperMessage(fmt.Sprintf("Unable to retrieve subscribers: %s", err))
		return
	}

	for _, subscription := range subscriptions {
		ch, err := bot.Session.UserChannelCreate(subscription.User)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = bot.Session.ChannelMessageSend(ch.ID, message)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
