package discord

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/hololive-livestream-notifier-go/data"
)

func SendWillLivestreamNotification(bot *discord.Bot, livestream data.Livestream, force bool) {
	if livestream.Date.Before(time.Now()) {
		return
	}

	storedLivestream, _ := data.GetLivestream(livestream.Url)

	dateChanged := false
	if storedLivestream != nil {
		dateChanged = storedLivestream.Date.Sub(livestream.Date) != time.Duration(0)
	}

	// Only send a `will livestream on` message if the time the livestream starts has changed
	if storedLivestream == nil || dateChanged || force {
		// load PST time zone
		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			slog.Error("Fail to load America/Los_Angeles location", "error", err)

			bot.SendDeveloperMessage(fmt.Sprint(err))
		}

		ts := livestream.Date.In(loc).Unix()
		message := fmt.Sprintf("%s livestream @<t:%d:F> <t:%d:R>\n%s", livestream.Author, ts, ts, livestream.Url)
		slog.Info(message)
		bot.SendChannelMessage("hololive-notifications", message)

		authors := strings.Split(strings.ToLower(livestream.Author), " ")
		SendSubscriberMessage(bot, authors, message)
	}
}
