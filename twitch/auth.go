package twitch

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/NekoFluff/hololive-livestream-notifier-go/utils"
)

func AccessToken() (string, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create the request body
	data := url.Values{}
	data.Set("client_id", utils.GetEnvVar("TWITCH_CLIENT_ID"))
	data.Set("client_secret", utils.GetEnvVar("TWITCH_CLIENT_SECRET"))
	data.Set("grant_type", "client_credentials")

	// Create the POST request
	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		// Handle error
		slog.Error("Error creating request: ", "error", err)
		return "", err
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		// Handle error
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		// Handle error
		return "", err
	}

	if result["status"] != nil {
		slog.Error("Failed to retrieve Twitch access token", "error", result["message"], "status", result["status"])
		return "", errors.New(result["message"].(string))
	}

	return result["access_token"].(string), nil
}
