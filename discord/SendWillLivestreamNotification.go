package discord

import (
	"fmt"
	"log"
	"time"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/go-hololive-notification-bot/data"
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
			log.Println(err)
			bot.SendDeveloperMessage(fmt.Sprint(err))
		}

		// e.g. [Flare Ch. 不知火フレア] Livestream on Mon, 02 Jan 2006 15:04:05 PST
		message := fmt.Sprintf("%s will livestream on [%s] - [%s]", livestream.Author, livestream.Date.In(loc).Format(time.RFC1123), livestream.Url)
		log.Println(message)
		bot.SendChannelMessage("hololive-notifications", message)
	}
}
