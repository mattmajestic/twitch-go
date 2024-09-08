package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"github.com/joho/godotenv"
	"mattmajestic/twitch-go/internal/models"
)

// LoadEnv loads environment variables from a .env file if available
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file, continuing with system environment variables.")
	}
}

// GetEnv retrieves an environment variable or returns a default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetAccessToken fetches an OAuth access token from the Twitch API
func GetAccessToken(clientID, clientSecret string) (string, error) {
	resp, err := http.PostForm("https://id.twitch.tv/oauth2/token", url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {"client_credentials"},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var accessTokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(body, &accessTokenResponse)
	if err != nil {
		return "", err
	}

	return accessTokenResponse.AccessToken, nil
}

// GetUserInfo fetches user information from the Twitch API
func GetUserInfo(clientID, accessToken, username string) (models.User, error) {
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		return models.User{}, err
	}

	q := req.URL.Query()
	q.Add("login", username)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.User{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.User{}, err
	}

	var usersResponse struct {
		Data []models.User `json:"data"`
	}
	err = json.Unmarshal(body, &usersResponse)
	if err != nil {
		return models.User{}, err
	}

	if len(usersResponse.Data) == 0 {
		return models.User{}, fmt.Errorf("no user found")
	}

	return usersResponse.Data[0], nil
}
