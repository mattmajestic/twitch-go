package models

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

// Follower represents a single follower of a Twitch user
type Follower struct {
	FromID     string `json:"from_id"`
	FromLogin  string `json:"from_login"`
	FromName   string `json:"from_name"`
	ToID       string `json:"to_id"`
	ToLogin    string `json:"to_login"`
	ToName     string `json:"to_name"`
	FollowedAt string `json:"followed_at"`
}

// FollowersResponse represents the structure of the response for fetching followers
type FollowersResponse struct {
	Data       []Follower `json:"data"`
	Total      int        `json:"total"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}