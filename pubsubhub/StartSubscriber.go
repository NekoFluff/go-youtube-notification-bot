package pubsubhub

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/NekoFluff/gobot/data"
	"github.com/NekoFluff/gobot/discord"
	"github.com/NekoFluff/gobot/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/dpup/gohubbub"
)

func StartSubscriber(webpage string, port int, s *discordgo.Session) {
	// Reschedule all notifications from the db
	discord.RecheduleAllLivestreamNotifications(s)

	// Get the youtube channel feeds to subscribe to
	channelFeeds, err := data.GetFeeds()
	if err != nil {
		log.Fatal(err)
	}

	client := gohubbub.NewClient(webpage, "YT Notifier")

	for _, channelFeed := range channelFeeds {
		fmt.Printf("%#v", channelFeed)

		topicURL := channelFeed.TopicURL
		client.DiscoverAndSubscribe(topicURL, func(contentType string, body []byte) {
			// Handle unexpected panics by sending a developer message in discord
			defer func() {
				if r := recover(); r != nil {
					debug.PrintStack()
					str := fmt.Sprintf("Recovered from panic. %v", r)
					log.Println(str)
					discord.SendDeveloperMessage(s, str)
				}
			}()

			// Handle update notification.
			feed, xmlError := ParseXML(string(body))

			if xmlError != nil {
				errorMsg := fmt.Sprintf("XML Parse Error %v", xmlError)
				log.Println(errorMsg)
				discord.SendDeveloperMessage(s, errorMsg)
			} else {
				ProcessFeed(s, feed)
			}
		})
	}
	client.StartAndServe("", port)
}

func ProcessFeed(s *discordgo.Session, feed Feed) {
	discord.SendDeveloperMessage(s, fmt.Sprintf("Processing feed: %#v", feed))
	for _, entry := range feed.Entries {
		log.Printf("%s - %s (%s)", entry.Title, entry.Author.Name, entry.Link)

		livestream, err := ConvertEntryToLivestream(entry)
		if err != nil {
			log.Println(err)
			discord.SendDeveloperMessage(s, fmt.Sprintf("%s is not a livestream. Error: %v", entry.Link.Href, err))

		} else {
			// We need to do this before saving the livestream so we can do some
			// comparison checks with the time the video goes live
			discord.SendWillLivestreamNotification(s, livestream)

			// Save the livestream and set up notifications
			data.SaveLivestream(livestream)
			discord.ScheduleLivestreamNotifications(s, livestream.Url, livestream.Date)
			discord.SendDeveloperMessage(s, fmt.Sprintf("Processed livestream: %s", livestream.Url))
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
