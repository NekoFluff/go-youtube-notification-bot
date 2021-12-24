package commands

import (
	"github.com/bwmarrin/discordgo"
)

var Unsubscribe = NewDiscordCommand(
	"unsubscribe",
	"Unsubscribe from some vtubers.",
	func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// args := strings.Split(m.Content, " ")[1:]

		// for _, arg := range args {
		// 	subscription := data.Subscription{
		// 		User:         m.Author.ID,
		// 		Subscription: arg,
		// 	}
		// 	data.DeleteSubscription(subscription)
		// }

		// _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Removed subscriptions from %s", args))
		// if err != nil {
		// 	log.Println(err)
		// }
	},
)
