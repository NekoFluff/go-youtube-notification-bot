package discord

import (
	"fmt"
	"log"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/go-hololive-notification-bot/data"
)

func RecheduleAllLivestreamNotifications(bot *discord.Bot) {
	livestreams, err := data.GetLivestreams()
	if err != nil {
		log.Println(err)
		bot.SendDeveloperMessage(fmt.Sprint(err))
	}

	for _, livestream := range livestreams {
		ScheduleLivestreamNotifications(bot, livestream, livestream.Date)
	}
}
