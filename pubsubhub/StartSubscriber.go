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

var sentMessages map[string]time.Time = make(map[string]time.Time)

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
			feed, xmlError := ParseXML(string(body))

			if xmlError != nil {
				errorMsg := fmt.Sprintf("XML Parse Error %v", xmlError)
				log.Println(errorMsg)
				discord.SendDeveloperMessage(dg, errorMsg)
			} else {
				processFeed(dg, feed)
			}
		})
	}
	client.StartAndServe("", port)
}

func processFeed(dg *discordgo.Session, feed Feed) {
	discord.SendDeveloperMessage(dg, fmt.Sprintf("Processing feed: %#v", feed))
	for _, entry := range feed.Entries {
		log.Printf("%s - %s (%s)", entry.Title, entry.Author.Name, entry.Link)

		livestream, err := convertEntryToLivestream(entry)
		if err != nil {
			log.Println(err)
			discord.SendDeveloperMessage(dg, fmt.Sprintf("%v", err))

		} else {
			data.SaveLivestream(livestream)
			discord.ScheduleLivestreamNotifications(dg, livestream.Url, livestream.Date)
			discord.SendDeveloperMessage(dg, fmt.Sprintf("Processed livestream: %s", livestream.Url))

			// Only send a `will livestream on` message if the time the livestream starts has changed
			if sentMessages[livestream.Url] != livestream.Date {
				// load PST time zone
				loc, err := time.LoadLocation("PST") // use other time zones such as MST, IST
				if err != nil {
					log.Println(err)
				}

				// e.g. [Flare Ch. 不知火フレア] Livestream on Mon, 02 Jan 2006 15:04:05 PST
				message := fmt.Sprintf("[%s] will livestream on [%s] - [%s]", livestream.Author, livestream.Date.In(loc).Format(time.RFC1123), livestream.Url)
				log.Println(message)
				discord.SendChannelMessage(dg, "gobot", message)
				sentMessages[livestream.Url] = livestream.Date
			}
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
