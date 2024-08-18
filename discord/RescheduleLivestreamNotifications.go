package discord

import (
	"fmt"
	"log/slog"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/hololive-livestream-notifier-go/data"
)

func RecheduleAllLivestreamNotifications(bot *discord.Bot) {
	livestreams, err := data.GetLivestreams()
	if err != nil {
		slog.Error("Failed to reschedule all livestream notifications", "error", err)

		bot.SendDeveloperMessage(fmt.Sprint(err))
	}

	for _, livestream := range livestreams {
		ScheduleLivestreamNotifications(bot, livestream, livestream.Date)
	}
}
