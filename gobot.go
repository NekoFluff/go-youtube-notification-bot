package main

import (
	"context"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/NekoFluff/gobot/commands"
	"github.com/NekoFluff/gobot/youtube/pubsubhub"
	"github.com/bwmarrin/discordgo"
	"github.com/dpup/gohubbub"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

type ChannelFeed struct {
	FirstName  string
	LastName   string
	TopicURL   string
	Group      string
	Generation int
}

func loadMongo() []ChannelFeed {
	uri := os.Getenv("MONGO_CONNECTION_URI")

	if uri == "" {
		log.Fatal("$MONGO_CONNECTION_URI must be set")
	}

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")

	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()
	// client.Database("<db>").Collection("<collection>").InsertOne(ctx, bson.D{{"x",1}})

	collection := client.Database("hololive-en").Collection("feeds")
	cur, err := collection.Find(context.Background(), bson.D{})
	var results []ChannelFeed
	if err = cur.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%#v", results)
	return results
}

func main() {
	// Load the .env file in the current directory
	godotenv.Load()

	channelFeeds := loadMongo()

	port, err := strconv.Atoi(os.Getenv("PORT"))

	if err != nil {
		log.Fatal(err)
	}

	if port == 0 {
		log.Fatal("$PORT must be set")
	}

	webpage := os.Getenv("WEBPAGE")
	if webpage == "" {
		log.Fatal("$WEBPAGE must be set")
	}

	for _, channelFeed := range channelFeeds {
		topicURL := channelFeed.TopicURL
		client := gohubbub.NewClient(webpage, "YT Notifier")
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
	}

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
