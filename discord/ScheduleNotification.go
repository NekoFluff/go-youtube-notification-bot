package discord

import (
	"time"

	"github.com/NekoFluff/gobot/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)

func ScheduleNotification(dg *discordgo.Session, t time.Time, message string) *cron.Cron {
	var c = cron.New()
	spec := utils.TimeToCron(t)
	c.AddFunc(spec, func() {
		SendChannelMessage(dg, "gobot", message)
	})
	c.Start()
	return c
}
