package pubsubhub

import (
	"encoding/xml"
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

func StartSubscriber(webpage string, port int, dg *discordgo.Session) {
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
			// Handle update notification.
			var feed Feed
			xmlError := xml.Unmarshal(body, &feed)

			if xmlError != nil {
				errorMsg := fmt.Sprintf("XML Parse Error %v", xmlError)
				log.Println(errorMsg)
				discord.SendDeveloperMessage(dg, errorMsg)

			} else {
				discord.SendDeveloperMessage(dg, fmt.Sprintf("Processing feed: %#v", feed))
				discord.SendDeveloperMessage(dg, fmt.Sprintf("Processing feed: %#v", feed))
				for _, entry := range feed.Entries {
					log.Printf("%s - %s (%s)", entry.Title, entry.Author.Name, entry.Link)

					livestream, err := convertEntryToLivestream(entry)
					if err != nil {
						log.Println(err)
						discord.SendDeveloperMessage(dg, fmt.Sprintf("%v", err))

					} else {
						data.SaveLivestream(livestream)
						// discord.SendChannelMessage(dg, [Author] will livestream on [Date] - [Link])
						discord.ScheduleLivestreamNotifications(dg, livestream.Url, livestream.Date)
						discord.SendDeveloperMessage(dg, fmt.Sprintf("Processed livestream %s", livestream.Url))

					}
				}
			}
		})
	}
	client.StartAndServe("", port)
}

func convertEntryToLivestream(entry Entry) (livestream data.Livestream, err error) {
	livestreamUnixTime, err := getLivestreamUnixTime(entry.Link)

	if err != nil {
		return
	}

	livestream = data.Livestream{
		Author: entry.Author.Name,
		Url:    entry.Link,
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
