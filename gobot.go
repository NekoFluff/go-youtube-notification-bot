package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/NekoFluff/gobot/discord"
	"github.com/NekoFluff/gobot/pubsubhub"
	"github.com/NekoFluff/gobot/utils"
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

// TODO: Tests
// TODO: Save and load subscriptions
// TODO: Subscription commnad (cmmands package)
// TODO: Documentation

func main() {
	// Load the .env file in the current directory
	godotenv.Load()

	// Start up discord bot
	token := utils.GetEnvVar("DISCORD_BOT_TOKEN")
	dg, err := discord.StartBot(token)
	if err != nil {
		log.Fatal(err)
	}
	defer discord.StopBot(dg)

	discord.SendChannelMessage(dg, "gobot", fmt.Sprintf("%s is online!", dg.State.User))

	// Load environment variables for pubsubhub subscriber
	webpage := utils.GetEnvVar("WEBPAGE")
	port := utils.GetEnvVar("PORT")

	// Translate port string into int
	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}

	// Start up new subscriber client
	pubsubhub.StartSubscriber(webpage, portInt, dg)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
