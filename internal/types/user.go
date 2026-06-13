package types

// User represents a Twitter/X user profile.
type User struct {
	ID          string `json:"id"`
	RestID      string `json:"rest_id"`
	ScreenName  string `json:"screen_name"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Verified    bool   `json:"is_blue_verified"`
	Affiliation string `json:"affiliation,omitempty"`

	FollowersCount int `json:"followers_count"`
	FriendsCount   int `json:"friends_count"`
	StatusesCount  int `json:"statuses_count"`
}
