package main

import (
	"encoding/xml"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/NekoFluff/gobot/commands"
	"github.com/NekoFluff/gobot/youtube/pubsubhub"
	"github.com/bwmarrin/discordgo"
	"github.com/dpup/gohubbub"
	"github.com/joho/godotenv"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	// Load the .env file in the current directory
	godotenv.Load()

	port, err := strconv.Atoi(os.Getenv("PORT"))

	if err != nil {
		log.Fatal(err)
	}

	if port == 0 {
		log.Fatal("$PORT must be set")
	}

	topicURL := "https://www.youtube.com/xml/feeds/videos.xml?channel_id=UCPem6W8TYuoSs0cIAnkKy6Q"
	client := gohubbub.NewClient("yt-notifier-bot.herokuapp.com", "YT Notifier")
	client.DiscoverAndSubscribe(topicURL, func(contentType string, body []byte) {
		// Handle update notification.
		var feed pubsubhub.Feed
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
	client.StartAndServe("", port)

	// // Create a new Discord session using the provided bot token.
	// dg, err := discordgo.New("Bot " + Token)
	// if err != nil {
	// 	fmt.Println("error creating Discord session,", err)
	// 	return
	// }

	// // Register the messageCreate func as a callback for MessageCreate events.
	// dg.AddHandler(messageCreate)

	// // In this example, we only care about receiving message events.
	// dg.Identify.Intents = discordgo.IntentsGuildMessages

	// // Open a websocket connection to Discord and begin listening.
	// err = dg.Open()
	// if err != nil {
	// 	fmt.Println("error opening connection,", err)
	// 	return
	// }

	// // Wait here until CTRL-C or other term signal is received.
	// fmt.Println("Bot is now running. Press CTRL-C to exit.")
	// sc := make(chan os.Signal, 1)
	// signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	// <-sc

	// // Cleanly close down the Discord session.
	// dg.Close()
}

type Gopher struct {
	Name string `json: "name"`
}

var AllCommands = []*commands.DiscordCommand{
	commands.Help,
	commands.GoRoutine,
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	for _, c := range AllCommands {
		if m.Content == c.Command {
			c.Execute(s, m)
		}
	}

	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}
