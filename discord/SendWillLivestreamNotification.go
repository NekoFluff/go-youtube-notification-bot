package discord

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/hololive-livestream-notifier-go/data"
)

func SendWillLivestreamNotification(bot *discord.Bot, livestream data.Livestream) {
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

		// e.g. [Flare Ch. 不知火フレア] Livestream on Mon, 02 Jan 2006 15:04:05 PST
		message := fmt.Sprintf("%s will livestream at <t:%d> - [%s]", livestream.Author, livestream.Date.In(loc).Unix(), livestream.Url)
		slog.Info(message)
		bot.SendChannelMessage("hololive-notifications", message)
	}
}
