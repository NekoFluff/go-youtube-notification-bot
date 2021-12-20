package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/NekoFluff/gobot/discord"
	"github.com/NekoFluff/gobot/pubsubhub"
	"github.com/NekoFluff/gobot/utils"
)

// TODO: Tests
// TODO: Save and load subscriptions
// TODO: Subscription commnad (cmmands package)
// TODO: Documentation

func main() {

	// Start up discord bot
	token := utils.GetEnvVar("DISCORD_BOT_TOKEN")
	s, err := discord.StartBot(token)
	if err != nil {
		log.Fatal(err)
	}
	defer discord.StopBot(s)

	discord.SendChannelMessage(s, "gobot", fmt.Sprintf("%s is online!", s.State.User))

	// Load environment variables for pubsubhub subscriber
	webpage := utils.GetEnvVar("WEBPAGE")
	port := utils.GetEnvVar("PORT")

	// Translate port string into int
	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}

	// Start up new subscriber client
	pubsubhub.StartSubscriber(webpage, portInt, s)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
