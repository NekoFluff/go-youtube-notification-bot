package pubsubhub

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"strconv"
	"time"

	mydiscord "github.com/NekoFluff/discord"
	"github.com/NekoFluff/go-hololive-notification-bot/data"
	"github.com/NekoFluff/go-hololive-notification-bot/discord"
	"github.com/NekoFluff/go-hololive-notification-bot/utils"
	"github.com/dpup/gohubbub"
)

func StartSubscriber(webpage string, port int, bot *mydiscord.Bot) {
	// Reschedule all notifications from the db
	discord.RecheduleAllLivestreamNotifications(bot)

	// Get the youtube channel feeds to subscribe to
	channelFeeds, err := data.GetFeeds()
	if err != nil {
		slog.Error("Failed to get feeds to subscribe to", "error", err)

		panic(err)
	}

	client := gohubbub.NewClient(webpage, "YT Notifier")

	for _, channelFeed := range channelFeeds {
		topicURL := channelFeed.TopicURL
		err = client.DiscoverAndSubscribe(topicURL, func(contentType string, body []byte) {
			// Handle unexpected panics by sending a developer message in discord
			defer func() {
				if r := recover(); r != nil {
					debug.PrintStack()
					str := fmt.Sprintf("Recovered from panic. %v", r)
					slog.Info(str)
					bot.SendDeveloperMessage(str)
				}
			}()

			// Handle update notification.
			feed, xmlError := ParseXML(string(body))

			if xmlError != nil {
				errorMsg := fmt.Sprintf("XML Parse Error %v", xmlError)
				slog.Error(errorMsg)
				bot.SendDeveloperMessage(errorMsg)
			} else {
				ProcessFeed(bot, feed)
			}
		})

		if err != nil {
			slog.Error("Failed to subscribe to feed", "error", err)
		}
	}
	client.StartAndServe("", port)
}

func ProcessFeed(bot *mydiscord.Bot, feed Feed) {
	j, err := json.MarshalIndent(feed, "", "  ")
	if err != nil {
		slog.Error("Failed to marshal indent feed", "error", err)
	} else {
		bot.SendDeveloperMessage(fmt.Sprintf("Processing feed:\n```json\n%s```", string(j)))
	}
	for _, entry := range feed.Entries {
		slog.Info("%s - %s (%s)\n", entry.Title, entry.Author.Name, entry.Link)

		livestream, err := ConvertEntryToLivestream(entry)
		if err != nil {
			slog.Error("Failed to convert the feed data into a Livestream object", "error", err)
			bot.SendDeveloperMessage(fmt.Sprintf("%s is not a livestream. Error: %v", entry.Link.Href, err))

		} else {
			// We need to do this before saving the livestream so we can do some
			// comparison checks with the time the video goes live
			discord.SendWillLivestreamNotification(bot, livestream)

			// Save the livestream and set up notifications
			data.SaveLivestream(livestream)
			discord.ScheduleLivestreamNotifications(bot, livestream, livestream.Date)
			bot.SendDeveloperMessage(fmt.Sprintf("Processed livestream: %s", livestream.Url))
		}
	}
}

func ConvertEntryToLivestream(entry Entry) (livestream data.Livestream, err error) {
	maximumAttempts := 3

	var livestreamUnixTime time.Time
	for i := 0; i < maximumAttempts; i++ {
		livestreamUnixTime, err = GetLivestreamUnixTime(entry.Link.Href)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(5) * time.Minute)
		}
	}

	// Failed all attempts to get the livestream unix time
	if err != nil {
		return
	}

	livestream = data.Livestream{
		Author:  entry.Author.Name,
		Url:     entry.Link.Href,
		Date:    livestreamUnixTime,
		Title:   entry.Title,
		Updated: entry.Updated,
	}

	return
}

func GetLivestreamUnixTime(url string) (t time.Time, err error) {
	html, err := utils.GetHTMLContent(url)
	if err != nil {
		return
	}
	params := utils.GetParams(`(?:"scheduledStartTime":")(?P<timestamp>\d+)`, string(html))
	// fmt.Printf("%s\n", params)

	// Translate port string into int
	if timestampStr, ok := params["timestamp"]; ok {
		var timestampInt int64
		timestampInt, err = strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			return
		}
		t = time.Unix(timestampInt, 0)
	} else {
		err = errors.New("no timestamp found")
		return
	}

	return
}
