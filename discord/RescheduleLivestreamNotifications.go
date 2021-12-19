package discord

import (
	"fmt"
	"log"

	"github.com/NekoFluff/gobot/data"
	"github.com/bwmarrin/discordgo"
)

func RecheduleAllLivestreamNotifications(dg *discordgo.Session) {
	livestreams, err := data.GetLivestreams()
	if err != nil {
		log.Println(err)
		SendDeveloperMessage(dg, fmt.Sprint(err))
	}

	for _, livestream := range livestreams {
		ScheduleLivestreamNotifications(dg, livestream.Url, livestream.Date)
	}
}
