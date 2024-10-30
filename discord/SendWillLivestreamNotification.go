package discord

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/hololive-livestream-notifier-go/data"
)

func SendWillLivestreamNotification(bot *discord.Bot, livestream data.Livestream) {
	if livestream.Date.Before(time.Now()) {
		return
	}

	storedLivestream, err := data.GetLivestream(livestream.Url)

	dateChanged := false
	if storedLivestream != nil {
		dateChanged = storedLivestream.Date.Sub(livestream.Date) != time.Duration(0)
	}

	// Only send a `will livestream on` message if the time the livestream starts has changed
	if err != nil || dateChanged {
		// load PST time zone
		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			slog.Error("Fail to load America/Los_Angeles location", "error", err)

			bot.SendDeveloperMessage(fmt.Sprint(err))
		}

		message := fmt.Sprintf("%s livestream @<t:%d:F> <t:%d:R>\n%s", livestream.Author, livestream.Date.In(loc).Unix(), livestream.Url)
		slog.Info(message)
		bot.SendChannelMessage("hololive-notifications", message)
	}
}
