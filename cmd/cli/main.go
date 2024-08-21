package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/NekoFluff/hololive-livestream-notifier-go/pubsubhub"
	"github.com/NekoFluff/hololive-livestream-notifier-go/twitch"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "twitch",
				Aliases: []string{"tw"},
				Usage:   "options for twitch",
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
				Usage:   "options for youtube",
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
