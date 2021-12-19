package pubsubhub

import (
	"fmt"
	"log"
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
				processFeed(s, feed)
			}
		})
	}
	client.StartAndServe("", port)
}

func processFeed(s *discordgo.Session, feed Feed) {
	discord.SendDeveloperMessage(s, fmt.Sprintf("Processing feed: %#v", feed))
	for _, entry := range feed.Entries {
		log.Printf("%s - %s (%s)", entry.Title, entry.Author.Name, entry.Link)

		livestream, err := convertEntryToLivestream(entry)
		if err != nil {
			log.Println(err)
			discord.SendDeveloperMessage(s, fmt.Sprintf("%s it not a livestream. Error: %v", entry.Link.Href, err))

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

func convertEntryToLivestream(entry Entry) (livestream data.Livestream, err error) {
	livestreamUnixTime, err := getLivestreamUnixTime(entry.Link.Href)

	if err != nil {
		return
	}

	livestream = data.Livestream{
		Author: entry.Author.Name,
		Url:    entry.Link.Href,
		Date:   livestreamUnixTime,
	}

	return
}

func getLivestreamUnixTime(url string) (time.Time, error) {
	html, err := utils.GetHTMLContent(url)
	if err != nil {
		return time.Now(), err
	}
	params := utils.GetParams(`(?:"scheduledStartTime":")(?P<timestamp>\d+)`, string(html))
	// fmt.Printf("%s\n", params)

	// Translate port string into int
	timestamp, err := strconv.ParseInt(params["timestamp"], 10, 64)
	if err != nil {
		return time.Now(), err
	}

	return time.Unix(timestamp, 0), nil
}
