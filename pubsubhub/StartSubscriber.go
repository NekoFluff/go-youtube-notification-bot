package pubsubhub

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/NekoFluff/gobot/data"
	"github.com/dpup/gohubbub"
)

func StartSubscriber(webpage string, port int) {
	// Get the youtube channel feeds to subscribe to
	channelFeeds := data.GetFeeds()

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
				}
			}
		})
	}
	client.StartAndServe("", port)
}
