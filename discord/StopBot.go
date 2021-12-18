package discord

import "github.com/bwmarrin/discordgo"

func StopBot(dg *discordgo.Session) {
	// Cleanly close down the Discord session.
	dg.Close()
}
