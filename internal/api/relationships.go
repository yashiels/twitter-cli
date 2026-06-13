package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/yashiels/twitter-cli/internal/types"
)

// parseRESTUser maps a REST v1.1 user object to types.User.
// REST v1.1 uses flat fields (screen_name, followers_count, etc.) not APK schema.
func parseRESTUser(raw json.RawMessage) (*types.User, error) {
	var u struct {
		IDStr          string `json:"id_str"`
		ScreenName     string `json:"screen_name"`
		Name           string `json:"name"`
		Description    string `json:"description"`
		Location       string `json:"location"`
		FollowersCount int    `json:"followers_count"`
		FriendsCount   int    `json:"friends_count"`
		StatusesCount  int    `json:"statuses_count"`
		Verified       bool   `json:"verified"`
	}

	if err := json.Unmarshal(raw, &u); err != nil {
		return nil, fmt.Errorf("parseRESTUser: %w", err)
	}

	return &types.User{
		RestID:         u.IDStr,
		ScreenName:     u.ScreenName,
		Name:           u.Name,
		Description:    u.Description,
		Location:       u.Location,
		FollowersCount: u.FollowersCount,
		FriendsCount:   u.FriendsCount,
		StatusesCount:  u.StatusesCount,
		Verified:       u.Verified,
	}, nil
}

// GetFollowers fetches a user's followers via REST v1.1.
func (c *Client) GetFollowers(handle string, limit int) ([]*types.User, error) {
	if limit <= 0 {
		limit = 20
	}

	params := url.Values{
		"screen_name": {handle},
		"count":       {strconv.Itoa(limit)},
		"skip_status": {"true"},
	}

	raw, err := c.restGet("/1.1/followers/list.json", params)
	if err != nil {
		return nil, fmt.Errorf("GetFollowers: %w", err)
	}

	return parseRESTUserList(raw)
}

// GetFollowing fetches who a user follows via REST v1.1.
func (c *Client) GetFollowing(handle string, limit int) ([]*types.User, error) {
	if limit <= 0 {
		limit = 20
	}

	params := url.Values{
		"screen_name": {handle},
		"count":       {strconv.Itoa(limit)},
		"skip_status": {"true"},
	}

	raw, err := c.restGet("/1.1/friends/list.json", params)
	if err != nil {
		return nil, fmt.Errorf("GetFollowing: %w", err)
	}

	return parseRESTUserList(raw)
}

// parseRESTTweetList parses a JSON array of REST v1.1 tweet objects (e.g. mentions_timeline).
func parseRESTTweetList(raw json.RawMessage) ([]*types.Tweet, error) {
	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("parseRESTTweetList: %w", err)
	}

	tweets := make([]*types.Tweet, 0, len(items))
	for _, item := range items {
		var t struct {
			IDStr                string `json:"id_str"`
			FullText             string `json:"full_text"`
			Text                 string `json:"text"` // fallback when tweet_mode!=extended
			CreatedAt            string `json:"created_at"`
			FavoriteCount        int    `json:"favorite_count"`
			RetweetCount         int    `json:"retweet_count"`
			ReplyCount           int    `json:"reply_count_v2"`
			InReplyToStatusIDStr string `json:"in_reply_to_status_id_str"`
			User                 struct {
				ScreenName string `json:"screen_name"`
			} `json:"user"`
		}
		if err := json.Unmarshal(item, &t); err != nil {
			continue
		}
		text := t.FullText
		if text == "" {
			text = t.Text
		}
		tweet := &types.Tweet{
			ID:            t.IDStr,
			Text:          text,
			FavoriteCount: t.FavoriteCount,
			RetweetCount:  t.RetweetCount,
			ReplyCount:    t.ReplyCount,
			AuthorHandle:  t.User.ScreenName,
			IsReply:       t.InReplyToStatusIDStr != "",
		}
		if t.CreatedAt != "" {
			// REST v1.1 date format: "Mon Jan 02 15:04:05 +0000 2006"
			if ts, err := time.Parse("Mon Jan 02 15:04:05 +0000 2006", t.CreatedAt); err == nil {
				tweet.CreatedAt = ts
			}
		}
		if t.User.ScreenName != "" && t.IDStr != "" {
			tweet.URL = "https://x.com/" + t.User.ScreenName + "/status/" + t.IDStr
		}
		tweets = append(tweets, tweet)
	}
	return tweets, nil
}

// parseRESTUserList parses the {"users": [...]} response from REST v1.1 list endpoints.
func parseRESTUserList(raw json.RawMessage) ([]*types.User, error) {
	var resp struct {
		Users []json.RawMessage `json:"users"`
	}

	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parseRESTUserList: %w", err)
	}

	users := make([]*types.User, 0, len(resp.Users))
	for _, userRaw := range resp.Users {
		u, err := parseRESTUser(userRaw)
		if err != nil || u == nil {
			continue
		}
		users = append(users, u)
	}

	return users, nil
}
