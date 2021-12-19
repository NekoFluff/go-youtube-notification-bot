package discord

import (
	"fmt"
	"log"

	"github.com/NekoFluff/gobot/data"
	"github.com/bwmarrin/discordgo"
)

func RecheduleAllLivestreamNotifications(s *discordgo.Session) {
	livestreams, err := data.GetLivestreams()
	if err != nil {
		log.Println(err)
		SendDeveloperMessage(s, fmt.Sprint(err))
	}

	for _, livestream := range livestreams {
		ScheduleLivestreamNotifications(s, livestream.Url, livestream.Date)
	}
}
