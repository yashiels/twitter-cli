package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/yashiels/twitter-cli/internal/types"
)

// GetMentions fetches the authenticated user's mentions via the REST v1.1 endpoint.
// The GraphQL NotificationTimelineQuery (jsOzc8RhpUpH5InTskP6Yw) returns HTTP 422
// with "notificationTimelineType must be defined" for unknown enum values, and
// HTTP 500 for known values — behaviour differs across APK versions. The REST
// statuses/mentions_timeline.json endpoint is stable and returns the same data.
func (c *Client) GetMentions(limit int) ([]*types.Tweet, error) {
	if limit <= 0 {
		limit = 20
	}

	params := url.Values{
		"count":            {strconv.Itoa(limit)},
		"tweet_mode":       {"extended"},
		"include_entities": {"true"},
	}

	raw, err := c.restGet("/1.1/statuses/mentions_timeline.json", params)
	if err != nil {
		return nil, fmt.Errorf("GetMentions: %w", err)
	}

	return parseRESTTweetList(raw)
}

// parseRESTTweetList parses a JSON array of REST v1.1 tweet objects.
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
			Text                 string `json:"text"`
			CreatedAt            string `json:"created_at"`
			FavoriteCount        int    `json:"favorite_count"`
			RetweetCount         int    `json:"retweet_count"`
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
			AuthorHandle:  t.User.ScreenName,
			IsReply:       t.InReplyToStatusIDStr != "",
		}
		if t.CreatedAt != "" {
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
