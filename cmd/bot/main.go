package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/go-hololive-notification-bot/commands"
	"github.com/NekoFluff/go-hololive-notification-bot/pubsubhub"
	"github.com/NekoFluff/go-hololive-notification-bot/utils"
)

// TODO: Tests
// TODO: Documentation

func main() {
	// Start up discord bot
	token := utils.GetEnvVar("DISCORD_BOT_TOKEN")
	bot := discord.NewBot(token)
	defer bot.Stop()

	bot.SendChannelMessage("hololive-notifications", fmt.Sprintf("%s is online!", bot.Session.State.User))

	// Generate Commands
	bot.AddCommands(
		commands.Ping(),
		commands.Subscribe(),
		commands.Unsubscribe(),
	)
	bot.RegisterCommands()

	// Load environment variables for pubsubhub subscriber
	webpage := utils.GetEnvVar("WEBPAGE")
	port := utils.GetEnvVar("PORT")

	// Translate port string into int
	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}

	// Start up new subscriber client
	pubsubhub.StartSubscriber(webpage, portInt, bot)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
