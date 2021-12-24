package discord

import (
	"time"

	"github.com/NekoFluff/gobot/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)

func ScheduleNotification(s *discordgo.Session, t time.Time, channel string, message string, authors []string) *cron.Cron {
	var c = cron.New()
	spec := utils.TimeToCron(t)
	c.AddFunc(spec, func() {
		SendChannelMessage(s, channel, message)
		SendSubscriberMessage(s, authors, message)
	})
	c.Start()
	return c
}
