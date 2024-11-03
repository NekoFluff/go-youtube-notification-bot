package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/hololive-livestream-notifier-go/commands"
	"github.com/NekoFluff/hololive-livestream-notifier-go/pubsubhub"
	"github.com/NekoFluff/hololive-livestream-notifier-go/twitch"
	"github.com/NekoFluff/hololive-livestream-notifier-go/utils"
	"github.com/urfave/cli/v2"

	internalDiscord "github.com/NekoFluff/hololive-livestream-notifier-go/discord"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "twitch",
				Aliases: []string{"tw"},
				Usage:   "twitch related commands",
				Subcommands: []*cli.Command{
					{
						Name:     "users",
						Aliases:  []string{"u"},
						Usage:    "get a list of users from Twitch given a username",
						Category: "Twitch",
						Action: func(cCtx *cli.Context) error {
							users, err := twitch.GetUsers(cCtx.Args().First())

							if err != nil {
								slog.Error("Failed to retrieve Twitch users", "error", err)
								return err
							}

							slog.Info("Successfully retrieved Twitch users", "users", users)
							return nil
						},
					},
					{
						Name:     "token",
						Aliases:  []string{"t"},
						Usage:    "get a Twitch access token",
						Category: "Twitch",
						Action: func(cCtx *cli.Context) error {
							accessToken, err := twitch.AccessToken()

							if err != nil {
								slog.Error("Failed to retrieve Twitch access token", "error", err)
								return nil
							}

							fmt.Println(accessToken)
							return nil

						},
					},
				},
			},
			{
				Name:    "youtube",
				Aliases: []string{"yt"},
				Usage:   "youtube related commands",
				Subcommands: []*cli.Command{
					{
						Name:     "timestamp",
						Aliases:  []string{"ts"},
						Usage:    "parse a Youtube livestream timestamp for a given url",
						Category: "Youtube",
						Action: func(cCtx *cli.Context) error {
							ts, err := pubsubhub.GetLivestreamUnixTime(cCtx.Args().First())

							if err != nil {
								slog.Error("Failed to retrieve livestream start time", "error", err)
								return err
							}

							fmt.Println("Unix Timestamp:", ts.Unix())
							phoenixTime := ts.In(time.FixedZone("Phoenix", -7*60*60))
							fmt.Println("Phoenix Time:", phoenixTime.Format("2006-01-02 03:04:05 PM"))

							return nil
						},
					},
				},
			},
			{
				Name:    "commands",
				Aliases: []string{"cmd", "cmds"},
				Usage:   "manage discord bot commands",
				Subcommands: []*cli.Command{
					{
						Name:     "refresh",
						Usage:    "refresh the Discord bot commands",
						Category: "Commands",
						Action: func(cCtx *cli.Context) error {
							slog.Info("Refreshing Discord bot commands")

							token := utils.GetEnvVar("DISCORD_BOT_TOKEN")
							bot := discord.NewBot(token)
							defer bot.Stop()

							bot.ClearCommands(os.Getenv("DISCORD_GUILD_ID"))
							bot.AddCommands(
								commands.Ping(),
								commands.Subscription(),
							)
							bot.RegisterCommands(os.Getenv("DISCORD_GUILD_ID"))
							slog.Info("Refreshed Discord bot commands")

							return nil
						},
					},
					{
						Name:     "clear",
						Usage:    "clear the Discord bot commands",
						Category: "Commands",
						Action: func(cCtx *cli.Context) error {
							slog.Info("Clearing Discord bot commands")

							token := utils.GetEnvVar("DISCORD_BOT_TOKEN")
							bot := discord.NewBot(token)
							defer bot.Stop()

							bot.ClearCommands(os.Getenv("DISCORD_GUILD_ID"))
							slog.Info("Cleared Discord bot commands")

							return nil
						},
					},
				},
			},
			{
				Name:      "notify",
				Usage:     "test subscriber notifications",
				ArgsUsage: "<vtuber_name> <message>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.NArg() != 2 {
						return fmt.Errorf("<vtuber_name> and <message> arguments required")
					}

					token := utils.GetEnvVar("DISCORD_BOT_TOKEN")
					bot := discord.NewBot(token)
					defer bot.Stop()

					bot.AddCommands(
						commands.Ping(),
						commands.Subscription(),
					)
					bot.RegisterCommands(os.Getenv("DISCORD_GUILD_ID"))

					args := cCtx.Args()
					internalDiscord.SendSubscriberMessage(bot, []string{args.Get(0)}, args.Get(1))

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
