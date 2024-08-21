package twitch

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/NekoFluff/hololive-livestream-notifier-go/utils"
)

type User struct {
	id          string
	login       string
	displayName string
}

func GetUsers(username string) ([]User, error) {
	bearerToken, err := AccessToken()
	if err != nil {
		return nil, err
	}

	// Create the request
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users?login="+username, nil)
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return nil, err
	}

	// Set the request headers
	req.Header.Set("Client-Id", utils.GetEnvVar("TWITCH_CLIENT_ID"))
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to send request", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if result["status"] != nil {
		slog.Error("Failed to retrieve Twitch users", "error", result["message"], "status", result["status"])
		return nil, errors.New(result["message"].(string))
	}

	users := make([]User, 0)
	for _, data := range result["data"].([]interface{}) {
		user := data.(map[string]interface{})
		users = append(users, User{
			id:          user["id"].(string),
			login:       user["login"].(string),
			displayName: user["display_name"].(string),
		})
	}

	return users, nil
}
