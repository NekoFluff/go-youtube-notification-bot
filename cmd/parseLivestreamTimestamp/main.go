package main

import (
	"log/slog"

	"github.com/NekoFluff/hololive-livestream-notifier-go/pubsubhub"
)

func main() {
	url := "https://www.youtube.com/watch?v=7J4yZBIvxdw"

	result, err := pubsubhub.GetLivestreamUnixTime(url)

	slog.Info("Successfully retrieved livestream start time", "result", result, "error", err, "timestamp", result.Unix())
}
