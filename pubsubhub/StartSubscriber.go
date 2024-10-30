package pubsubhub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"runtime/debug"
	"time"

	mydiscord "github.com/NekoFluff/discord"
	"github.com/NekoFluff/hololive-livestream-notifier-go/data"
	"github.com/NekoFluff/hololive-livestream-notifier-go/discord"
	"github.com/NekoFluff/hololive-livestream-notifier-go/utils"
	"github.com/dpup/gohubbub"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
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
		slog.Info(fmt.Sprintf("%s - %s (%s)\n", entry.Title, entry.Author.Name, entry.Link))

		livestream, err := ConvertEntryToLivestream(entry)
		if err != nil {
			slog.Error("Failed to convert the feed data into a Livestream object", "error", err)
			bot.SendDeveloperMessage(fmt.Sprintf("%s is not a livestream. Error: %v", entry.Link.Href, err))

		} else {
			if livestream.Date.Before(time.Now()) {
				bot.SendDeveloperMessage(fmt.Sprintf("Livestream %s has already started or has ended", livestream.Url))
				continue
			}

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
	livestreamUnixTime, err := GetLivestreamUnixTime(entry.Link.Href)

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
	svc, err := youtube.NewService(context.Background(), option.WithAPIKey(utils.GetEnvVar("YOUTUBE_API_KEY")))
	if err != nil {
		slog.Error("Failed to create youtube service", "error", err)
		return
	}

	// Usage:
	videoID, err := GetVideoID(url)
	if err != nil {
		slog.Error("Failed to get video ID", "error", err)
		return
	}

	call := svc.Videos.List([]string{"liveStreamingDetails"}).Id(videoID)
	response, err := call.Do()

	if err != nil {
		slog.Error("Failed to get video details", "error", err)
		return
	}

	if len(response.Items) == 0 {
		slog.Error("No video details found")
		return
	}

	if response.Items[0].LiveStreamingDetails == nil {
		err = errors.New("no live streaming details found for the video")
		return
	}

	if response.Items[0].LiveStreamingDetails.ScheduledStartTime == "" {
		err = errors.New("video is not scheduled")
		return
	}

	t, err = time.Parse(time.RFC3339, response.Items[0].LiveStreamingDetails.ScheduledStartTime)
	if err != nil {
		slog.Error("Failed to parse scheduled start time", "error", err)
		return
	}

	return t, nil
}

func GetVideoID(url string) (string, error) {
	videoIDRegex := `\?v=(.*)`
	re := regexp.MustCompile(videoIDRegex)
	match := re.FindStringSubmatch(url)
	if len(match) < 2 {
		return "", errors.New("failed to parse video ID from URL")
	}
	return match[1], nil
}
