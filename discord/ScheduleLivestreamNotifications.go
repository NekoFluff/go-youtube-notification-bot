package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/NekoFluff/gobot/data"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)

var liveCronJobs map[string]*cron.Cron = make(map[string]*cron.Cron)
var fifteenMinCronJobs map[string]*cron.Cron = make(map[string]*cron.Cron)

func ScheduleLivestreamNotifications(s *discordgo.Session, livestream data.Livestream, t time.Time) {
	url := livestream.Url

	authors := strings.Split(strings.ToLower(livestream.Author), " ")
	if fifteenMinCronJobs[url] != nil {
		fifteenMinCronJobs[url].Stop()
	}
	fifteenMinCronJobs[url] = ScheduleNotification(s, t.Add(time.Duration(-15)*time.Minute), "hololive-notifications", fmt.Sprintf("Livestream starting in 15 minutes! %s", url), authors)

	if liveCronJobs[url] != nil {
		liveCronJobs[url].Stop()
	}
	liveCronJobs[url] = ScheduleNotification(s, t, "hololive-stream-started", fmt.Sprintf("Livestream starting! %s", url), authors)
}
