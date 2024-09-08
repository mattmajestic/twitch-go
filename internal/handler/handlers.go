package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

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

func FollowersHandler(w http.ResponseWriter, r *http.Request) {
	// Load environment variables
	services.LoadEnv()

	// Fetch the necessary environment variables
	clientID := services.GetEnv("TWITCH_CLIENT_ID", "")
	clientSecret := services.GetEnv("TWITCH_CLIENT_SECRET", "")
	username := services.GetEnv("TWITCH_USERNAME", "")

	// Get the access token
	accessToken, err := services.GetAccessToken(clientID, clientSecret)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get access token: %v", err), http.StatusInternalServerError)
		return
	}

	// Get user information (to find followers for the given user)
	user, err := services.GetUserInfo(clientID, accessToken, username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}

	// Display follower details (for simplicity, we’ll just return the username)
	response := map[string]string{
		"Follower": fmt.Sprintf("User %s has followers", user.DisplayName),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}