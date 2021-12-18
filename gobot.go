package main

import (
	"flag"
	"log"
	"strconv"

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

// TODO: Schedule Cron Jobs (cron package)
// TODO: Save and load schedules livestreams (data package)
// TODO: Server channel messaging on discord (discord package)
// TODO: DM for debugging
// TODO: env variable for developers
// TODO: evn variable for developer mode
// TODO: Tests
// TODO: Save and load subscriptions
// TODO: Subscription commnad (cmmands package)
// TODO: Documentation

func main() {
	// Load the .env file in the current directory
	godotenv.Load()

	// Load environment variables
	webpage := utils.GetEnvVar("WEBPAGE")
	port := utils.GetEnvVar("PORT")
	token := utils.GetEnvVar("DISCORD_BOT_TOKEN")

	// Translate port string into int
	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}

	// Start up new subscriber client
	pubsubhub.StartSubscriber(webpage, portInt)

	// Start up discord bot
	dg, err := discord.StartBot(token)
	if err != nil {
		log.Fatal(err)
	}
	defer discord.StopBot(dg)
}
