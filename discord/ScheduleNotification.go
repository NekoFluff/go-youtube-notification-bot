package discord

import (
	"fmt"
	"time"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/go-hololive-notification-bot/utils"
	"github.com/robfig/cron"
)

func ScheduleNotification(bot *discord.Bot, t time.Time, channel string, message string, authors []string) *cron.Cron {
	var c = cron.New()
	spec := utils.TimeToCron(t)
	err := c.AddFunc(spec, func() {
		bot.SendChannelMessage(channel, message)
		SendSubscriberMessage(bot, authors, message)
	})
	if err != nil {
		fmt.Printf("Failed to schedule notification: %s", err)
	}
	c.Start()
	return c
}
