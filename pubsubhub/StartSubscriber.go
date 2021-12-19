package pubsubhub

import (
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/NekoFluff/gobot/data"
	"github.com/NekoFluff/gobot/utils"
	"github.com/dpup/gohubbub"
)

func StartSubscriber(webpage string, port int) {
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
				log.Printf("XML Parse Error %v", xmlError)

			} else {
				log.Println("Feed title:", feed.Title)
				for _, entry := range feed.Entries {
					log.Printf("%s - %s (%s)", entry.Title, entry.Author.Name, entry.Link)

					livestream, err := convertEntryToLivestream(entry)
					if err != nil {
						log.Println(err)
					} else {
						data.SaveLivestream(livestream)
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
