package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/hololive-livestream-notifier-go/commands"
	"github.com/NekoFluff/hololive-livestream-notifier-go/pubsubhub"
	"github.com/NekoFluff/hololive-livestream-notifier-go/twitch"
	"github.com/NekoFluff/hololive-livestream-notifier-go/utils"
)

// TODO: Tests
// TODO: Documentation

func main() {
	// Start up discord bot
	token := utils.GetEnvVar("DISCORD_BOT_TOKEN")
	bot := discord.NewBot(token)
	defer bot.Stop()

	if utils.GetEnvVar("DEVELOPER_MODE") == "ON" {
		bot.DeveloperIDs = strings.Split(utils.GetEnvVar("DEVELOPER_IDS"), ",")
	}

	bot.SendDeveloperMessage(fmt.Sprintf("%s is online!", bot.Session.State.User))

	// Generate Commands
	bot.AddCommands(
		commands.Ping(),
		commands.Subscription(),
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

	setupTwitchCallbackEndpoint(bot)

	// Start up new subscriber client
	go pubsubhub.StartSubscriber(webpage, portInt, bot)

	subscribeToTwitchWebhook(bot)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func setupTwitchCallbackEndpoint(bot *discord.Bot) {
	http.HandleFunc("/twitch/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read request body
		var requestBody map[string]interface{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Failed to read request body", "error", err)
			return
		}
		defer r.Body.Close()

		// Decode request body
		err = json.NewDecoder(bytes.NewReader(body)).Decode(&requestBody)
		if err != nil {
			slog.Error("Failed to decode request body", "error", err)
			return
		}

		// Verify the message signature
		message := getHmacMessage(r, body)
		hmac := "sha256=" + getHmac([]byte(utils.GetEnvVar("TWITCH_WEBHOOK_SECRET")), message)
		if !verifyTwitchMessage(hmac, r.Header["Twitch-Eventsub-Message-Signature"][0]) {
			slog.Error("Failed to verify message signature")
			http.Error(w, "Failed to verify message signature", http.StatusUnauthorized)
			return
		}

		// Handle the message
		messageType := r.Header["Twitch-Eventsub-Message-Type"][0]
		if messageType == "notification" {
			// Handle the notification
			subscription := requestBody["subscription"].(map[string]interface{})
			slog.Info("Event type", "type", subscription["type"])
			w.WriteHeader(http.StatusNoContent)

			bot.SendDeveloperMessage(fmt.Sprintf("Processing Twitch webhook notification of type %s", subscription["type"]))
			if subscription["type"] == "stream.online" {
				handleStreamOnlineEventType(subscription, bot)
			}
			bot.SendDeveloperMessage(fmt.Sprintf("Processed Twitch webhook notification of type %s", subscription["type"]))

		} else if messageType == "webhook_callback_verification" {
			// Respond with the challenge
			slog.Info("Responding with challenge", "challenge", requestBody["challenge"])
			challenge := requestBody["challenge"].(string)
			_, _ = w.Write([]byte(challenge)) // yolo. dont need error handling here

			bot.SendDeveloperMessage("Twitch webhook verified. Good to go!")

		} else if messageType == "revocation" {
			// Handle the revocation
			subscription := requestBody["subscription"].(map[string]interface{})
			slog.Error("Revoked subscription", "reason", subscription["status"], "condition", subscription["condition"])
			w.WriteHeader(http.StatusNoContent)

			bot.SendDeveloperMessage("@crazyfluff Twitch webhook access revoked!!!!")
			bot.SendDeveloperMessage("@crazyfluff Twitch webhook access revoked!!!!")
			bot.SendDeveloperMessage("@crazyfluff Twitch webhook access revoked!!!!")
		} else {
			slog.Error("Unknown message type", "type", messageType)
			w.WriteHeader(http.StatusNoContent)
		}
	})
}

func verifyTwitchMessage(hmac string, header string) bool {
	slog.Info("Comparing HMACs", "hmac", hmac, "header", header)
	return hmac == header
}

// Build the message used to get the HMAC.
func getHmacMessage(r *http.Request, body []byte) []byte {
	return []byte(r.Header["Twitch-Eventsub-Message-Id"][0] +
		r.Header["Twitch-Eventsub-Message-Timestamp"][0] +
		string(body))
}

// Get the HMAC.
func getHmac(secret []byte, message []byte) string {
	hmac := hmac.New(sha256.New, secret)
	hmac.Write([]byte(message))
	return hex.EncodeToString(hmac.Sum(nil))
}

func subscribeToTwitchWebhook(bot *discord.Bot) {
	secret := utils.GetEnvVar("TWITCH_WEBHOOK_SECRET")
	callback := utils.GetEnvVar("TWITCH_WEBHOOK_CALLBACK")
	bearerToken, err := twitch.AccessToken()

	if err != nil {
		slog.Error("Failed to get Twitch access token", "error", err)

		bot.SendDeveloperMessage("Failed to get Twitch access token")
	}

	// Create the request
	req, err := http.NewRequest("POST", "https://api.twitch.tv/helix/eventsub/subscriptions", nil)
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return
	}

	// Set the request headers
	req.Header.Set("Client-Id", utils.GetEnvVar("TWITCH_CLIENT_ID"))
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	// Create the request body
	body := map[string]interface{}{
		"type":    "stream.online",
		"version": "1",
		"condition": map[string]interface{}{
			"broadcaster_user_id": "734637782",
		},
		"transport": map[string]interface{}{
			"method":   "webhook",
			"callback": callback,
			"secret":   secret,
		},
	}

	// Marshal the request body
	jsonBody, err := json.Marshal(body)
	if err != nil {
		slog.Error("Failed to marshal request body", "error", err)
		return
	}

	// Set the request body
	req.Body = io.NopCloser(bytes.NewReader(jsonBody))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to send request", "error", err)
		return
	}

	// Close the response body
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body", "error", err)
		return
	}

	// Log the response
	slog.Info("Sent subscription request", "response", string(respBody))
}

func handleStreamOnlineEventType(subscription map[string]interface{}, bot *discord.Bot) {
	slog.Info("Stream online event", "subscription", subscription)

	bot.SendDeveloperMessage("@crazyfluff Nao is live on Twitch!")
	bot.SendDeveloperMessage("@crazyfluff Nao is live on Twitch!")
	bot.SendDeveloperMessage("@crazyfluff Nao is live on Twitch!")
	bot.SendDeveloperMessage("@crazyfluff Nao is live on Twitch!")
	bot.SendDeveloperMessage("@crazyfluff Nao is live on Twitch!")
	bot.SendDeveloperMessage("@crazyfluff Nao is live on Twitch!")
	bot.SendDeveloperMessage("@crazyfluff Nao is live on Twitch!")
	bot.SendDeveloperMessage("@crazyfluff Nao is live on Twitch!")
	bot.SendDeveloperMessage("@crazyfluff Nao is live on Twitch!")
	bot.SendDeveloperMessage("@crazyfluff Nao is live on Twitch!")
}
