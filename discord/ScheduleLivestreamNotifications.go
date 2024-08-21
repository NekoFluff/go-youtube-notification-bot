package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/hololive-livestream-notifier-go/data"
	"github.com/robfig/cron"
)

var liveCronJobs map[string]*cron.Cron = make(map[string]*cron.Cron)
var fifteenMinCronJobs map[string]*cron.Cron = make(map[string]*cron.Cron)

func ScheduleLivestreamNotifications(bot *discord.Bot, livestream data.Livestream, t time.Time) {
	url := livestream.Url

	authors := strings.Split(strings.ToLower(livestream.Author), " ")
	if fifteenMinCronJobs[url] != nil {
		fifteenMinCronJobs[url].Stop()
	}
	fifteenMinCronJobs[url] = ScheduleNotification(bot, t.Add(time.Duration(-15)*time.Minute), "hololive-notifications", fmt.Sprintf("[%s] Livestream starting in 15 minutes! %s", livestream.Author, url), authors)

	if liveCronJobs[url] != nil {
		liveCronJobs[url].Stop()
	}
	liveCronJobs[url] = ScheduleNotification(bot, t, "hololive-stream-started", fmt.Sprintf("[%s] Livestream starting! %s", livestream.Author, url), authors)
}
