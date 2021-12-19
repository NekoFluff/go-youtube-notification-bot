package discord

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)

var liveCronJobs map[string]*cron.Cron = make(map[string]*cron.Cron)
var fifteenMinCronJobs map[string]*cron.Cron = make(map[string]*cron.Cron)

func ScheduleLivestreamNotifications(dg *discordgo.Session, url string, t time.Time) {
	if fifteenMinCronJobs[url] != nil {
		fifteenMinCronJobs[url].Stop()
	}
	fifteenMinCronJobs[url] = ScheduleNotification(dg, t.Add(time.Duration(-15)*time.Minute), "gobot", fmt.Sprintf("Livestream starting in 15 minutes! %s", url))

	if liveCronJobs[url] != nil {
		liveCronJobs[url].Stop()
	}
	liveCronJobs[url] = ScheduleNotification(dg, t, "gobot-live", fmt.Sprintf("Livestream starting! %s", url))
}
