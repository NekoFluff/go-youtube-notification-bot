package discord

import (
	"fmt"
	"log"
	"time"

	"github.com/NekoFluff/gobot/data"
	"github.com/bwmarrin/discordgo"
)

var sentMessages map[string]time.Time = make(map[string]time.Time)

func SendWillLivestreamNotification(s *discordgo.Session, livestream data.Livestream) {
	// Only send a `will livestream on` message if the time the livestream starts has changed
	if sentMessages[livestream.Url] != livestream.Date {
		// load PST time zone
		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			log.Println(err)
		}

		// e.g. [Flare Ch. 不知火フレア] Livestream on Mon, 02 Jan 2006 15:04:05 PST
		message := fmt.Sprintf("[%s] will livestream on [%s] - [%s]", livestream.Author, livestream.Date.In(loc).Format(time.RFC1123), livestream.Url)
		log.Println(message)
		SendChannelMessage(s, "gobot", message)
		sentMessages[livestream.Url] = livestream.Date
	}
}
