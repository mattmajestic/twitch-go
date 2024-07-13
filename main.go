package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
}

type UsersResponse struct {
	Data []User `json:"data"`
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

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

	var accessTokenResponse AccessTokenResponse
	err = json.Unmarshal(body, &accessTokenResponse)
	if err != nil {
		return "", err
	}

	return accessTokenResponse.AccessToken, nil
}

func getUserInfo(clientID, accessToken, username string) (User, error) {
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		return User{}, err
	}

	q := req.URL.Query()
	q.Add("login", username)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return User{}, err
	}

	var usersResponse UsersResponse
	err = json.Unmarshal(body, &usersResponse)
	if err != nil {
		return User{}, err
	}

	if len(usersResponse.Data) == 0 {
		return User{}, fmt.Errorf("no user found")
	}

	return usersResponse.Data[0], nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	loadEnv()

	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	username := os.Getenv("TWITCH_USERNAME")

	accessToken, err := getAccessToken(clientID, clientSecret)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get access token: %s", err), http.StatusInternalServerError)
		return
	}

	user, err := getUserInfo(clientID, accessToken, username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %s", err), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse template: %s", err), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, user)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
