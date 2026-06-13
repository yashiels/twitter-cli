package types

import "time"

// Tweet represents a single tweet/post.
type Tweet struct {
	ID            string    `json:"id"`
	Text          string    `json:"text"`
	CreatedAt     time.Time `json:"created_at"`
	FavoriteCount int       `json:"favorite_count"`
	RetweetCount  int       `json:"retweet_count"`
	ReplyCount    int       `json:"reply_count"`
	ViewCount     int       `json:"view_count"`
	AuthorHandle  string    `json:"author_handle"`
	URL           string    `json:"url"`
	IsRetweet     bool      `json:"is_retweet"`
	IsReply       bool      `json:"is_reply"`
	QuotedTweet   *Tweet    `json:"quoted_tweet,omitempty"`
}

// TweetURL returns the canonical URL to the tweet.
func (t *Tweet) TweetURL() string {
	if t.AuthorHandle != "" && t.ID != "" {
		return "https://x.com/" + t.AuthorHandle + "/status/" + t.ID
	}
	return ""
}
