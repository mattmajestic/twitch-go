package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"text/template"
	"mattmajestic/twitch-go/internal/models"
	"mattmajestic/twitch-go/internal/services"
)

// HomeHandler serves the home page with Twitch user information.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	services.LoadEnv() // Load environment variables from .env

	clientID := services.GetEnv("TWITCH_CLIENT_ID", "")
	clientSecret := services.GetEnv("TWITCH_CLIENT_SECRET", "")
	username := services.GetEnv("TWITCH_USERNAME", "")

	accessToken, err := services.GetAccessToken(clientID, clientSecret)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get access token: %v", err), http.StatusInternalServerError)
		return
	}

	user, err := services.GetUserInfo(clientID, accessToken, username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("static/template.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse template: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute template: %v", err), http.StatusInternalServerError)
	}
}

// FollowersHandler handles the /followers route and returns follower information
func FollowersHandler(w http.ResponseWriter, r *http.Request) {
	// Load environment variables
	services.LoadEnv()

	// Fetch necessary environment variables
	clientID := services.GetEnv("TWITCH_CLIENT_ID", "")
	clientSecret := services.GetEnv("TWITCH_CLIENT_SECRET", "")
	username := services.GetEnv("TWITCH_USERNAME", "")

	// Get access token
	accessToken, err := getAccessToken(clientID, clientSecret)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get access token: %v", err), http.StatusInternalServerError)
		return
	}

	// Get user info to obtain user ID
	user, err := getUserInfo(clientID, accessToken, username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}

	// Fetch followers using the user's ID
	followers, err := getFollowers(clientID, accessToken, user.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get followers: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the followers list as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

// getAccessToken fetches an OAuth access token from the Twitch API
func getAccessToken(clientID, clientSecret string) (string, error) {
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

	var accessTokenResponse models.AccessTokenResponse
	err = json.Unmarshal(body, &accessTokenResponse)
	if err != nil {
		return "", err
	}

	return accessTokenResponse.AccessToken, nil
}

// getUserInfo fetches user information from the Twitch API
func getUserInfo(clientID, accessToken, username string) (models.User, error) {
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		return models.User{}, err
	}

	// Set query parameters and headers
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

	var usersResponse models.UsersResponse
	err = json.Unmarshal(body, &usersResponse)
	if err != nil {
		return models.User{}, err
	}

	if len(usersResponse.Data) == 0 {
		return models.User{}, fmt.Errorf("no user found")
	}

	return usersResponse.Data[0], nil
}

// getFollowers fetches the list of followers for the given user
func getFollowers(clientID, accessToken, userID string) ([]models.Follower, error) {
	// Define the Twitch API endpoint to get the followers
	url := fmt.Sprintf("https://api.twitch.tv/helix/users/follows?to_id=%s", userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set required headers for the request
	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON response into the appropriate model
	var followersResponse models.FollowersResponse
	err = json.Unmarshal(body, &followersResponse)
	if err != nil {
		return nil, err
	}

	return followersResponse.Data, nil
}
