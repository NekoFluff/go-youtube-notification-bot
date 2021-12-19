package discord

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func SendDeveloperMessage(s *discordgo.Session, message string) {
	developer_mode := os.Getenv("DEVELOPER_MODE")
	if developer_mode != "ON" && developer_mode != "1" {
		return
	}

	developerIds := getDeveloperIds()

	for _, developerId := range developerIds {
		ch, err := s.UserChannelCreate(developerId)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = s.ChannelMessageSend(ch.ID, message)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func getDeveloperIds() []string {
	ids := os.Getenv("DEVELOPER_IDS")
	if ids == "" {
		return []string{}
	}
	return strings.Split(ids, ",")
}
