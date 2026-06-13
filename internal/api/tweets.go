package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/yashiels/twitter-cli/internal/types"
)

const (
	tweetsQueryID   = "xlAB_H3dvYL4q1C-PzM_ag"
	tweetsOperation = "UserProfileOriginalsTimelineQuery"
)

// tweetsFeatures are the feature flags for UserProfileOriginalsTimelineQuery.
var tweetsFeatures = map[string]any{
	"android_home_timeline_status_injections_facepile_enabled": true,
	"grok_translations_post_auto_translation_is_enabled":       true,
	"reply_dislike_android_enabled":                            true,
	"subscriptions_feature_can_gift_premium":                   true,
	"x_lite_quick_promote_analytics_banner_enabled":            true,
}

// GetUserTweets fetches recent tweets for the given user ID (rest_id).
func (c *Client) GetUserTweets(userID string, limit int) ([]*types.Tweet, error) {
	if limit <= 0 {
		limit = 20
	}

	vars := map[string]any{
		"rest_id": userID,
		"count":   limit,
	}

	raw, err := c.graphqlGet(tweetsQueryID, tweetsOperation, vars, tweetsFeatures)
	if err != nil {
		return nil, fmt.Errorf("GetUserTweets: %w", err)
	}

	return parseTweetsResponse(raw)
}

// rawTweetResult is used for JSON parsing of tweet entries.
type rawTweetResult struct {
	Typename string `json:"__typename"`

	// Direct tweet fields (TweetResults)
	RestID string `json:"rest_id"`
	Legacy struct {
		FullText             string `json:"full_text"`
		CreatedAt            string `json:"created_at"`
		FavoriteCount        int    `json:"favorite_count"`
		RetweetCount         int    `json:"retweet_count"`
		ReplyCount           int    `json:"reply_count"`
		InReplyToScreenName  string `json:"in_reply_to_screen_name"`
		RetweetedStatusIDStr string `json:"retweeted_status_id_str"`
		UserIDStr            string `json:"user_id_str"`
	} `json:"legacy"`
	NoteTweet struct {
		NoteTweetResults struct {
			Result struct {
				Text string `json:"text"`
			} `json:"result"`
		} `json:"note_tweet_results"`
	} `json:"note_tweet"`
	Views struct {
		Count string `json:"count"`
	} `json:"views"`
	Core struct {
		UserResults struct {
			Result struct {
				Core struct {
					ScreenName string `json:"screen_name"`
				} `json:"core"`
				Legacy *struct {
					ScreenName string `json:"screen_name"`
				} `json:"legacy"`
			} `json:"result"`
		} `json:"user_results"`
	} `json:"core"`

	// Wrapped tweet (TweetWithVisibilityResults)
	Tweet *rawTweetResult `json:"tweet"`

	// Quoted tweet
	QuotedStatusResult *struct {
		Result *rawTweetResult `json:"result"`
	} `json:"quoted_status_result"`
}

// parseTweetsResponse parses the timeline response into a slice of Tweets.
func parseTweetsResponse(raw json.RawMessage) ([]*types.Tweet, error) {
	// Navigate: data -> user_result -> result -> timeline_response -> timeline -> instructions
	timelineRaw, err := getNestedJSON(raw, "data", "user_result", "result", "timeline_response", "timeline", "instructions")
	if err != nil {
		// Try alternate path
		timelineRaw, err = getNestedJSON(raw, "data", "user_result", "result", "timeline", "timeline", "instructions")
		if err != nil {
			return nil, fmt.Errorf("navigate timeline: %w", err)
		}
	}

	var instructions []json.RawMessage
	if err := json.Unmarshal(timelineRaw, &instructions); err != nil {
		return nil, fmt.Errorf("parse instructions: %w", err)
	}

	var tweets []*types.Tweet
	for _, instr := range instructions {
		var inst struct {
			Type    string            `json:"__typename"`
			Entries []json.RawMessage `json:"entries"`
		}
		// Also try "type" field
		var instAlt struct {
			Type    string            `json:"type"`
			Entries []json.RawMessage `json:"entries"`
		}
		_ = json.Unmarshal(instr, &inst)
		_ = json.Unmarshal(instr, &instAlt)

		entries := inst.Entries
		if len(entries) == 0 {
			entries = instAlt.Entries
		}

		for _, entry := range entries {
			t, err := parseTimelineEntry(entry)
			if err != nil || t == nil {
				continue
			}
			tweets = append(tweets, t)
		}
	}

	return tweets, nil
}

// parseTimelineEntry extracts a Tweet from a timeline entry.
func parseTimelineEntry(entry json.RawMessage) (*types.Tweet, error) {
	// Structure: entry -> content -> content -> tweet_results -> result
	tweetResultRaw, err := getNestedJSON(entry, "content", "content", "tweet_results", "result")
	if err != nil {
		// Try alternate: entry -> content -> itemContent -> tweet_results -> result
		tweetResultRaw, err = getNestedJSON(entry, "content", "itemContent", "tweet_results", "result")
		if err != nil {
			return nil, nil // not a tweet entry
		}
	}

	return parseTweetResult(tweetResultRaw)
}

// parseTweetResult converts a raw tweet result JSON into a Tweet.
func parseTweetResult(raw json.RawMessage) (*types.Tweet, error) {
	var result rawTweetResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}

	// Unwrap TweetWithVisibilityResults
	if result.Typename == "TweetWithVisibilityResults" && result.Tweet != nil {
		return convertRawTweet(result.Tweet), nil
	}

	return convertRawTweet(&result), nil
}

// convertRawTweet converts a rawTweetResult to a types.Tweet.
func convertRawTweet(r *rawTweetResult) *types.Tweet {
	if r == nil || r.RestID == "" {
		return nil
	}

	// Prefer note_tweet text (long tweets), fall back to legacy.full_text.
	text := r.NoteTweet.NoteTweetResults.Result.Text
	if text == "" {
		text = r.Legacy.FullText
	}
	if text == "" {
		return nil
	}

	// Parse created_at ("Mon Jan 02 15:04:05 +0000 2006" format).
	var createdAt time.Time
	if r.Legacy.CreatedAt != "" {
		parsed, err := time.Parse("Mon Jan 02 15:04:05 +0000 2006", r.Legacy.CreatedAt)
		if err == nil {
			createdAt = parsed
		}
	}

	// Extract view count (stored as a string).
	var viewCount int
	if r.Views.Count != "" {
		viewCount, _ = strconv.Atoi(r.Views.Count)
	}

	// Determine author handle.
	handle := r.Core.UserResults.Result.Core.ScreenName
	if handle == "" && r.Core.UserResults.Result.Legacy != nil {
		handle = r.Core.UserResults.Result.Legacy.ScreenName
	}

	tweet := &types.Tweet{
		ID:            r.RestID,
		Text:          text,
		CreatedAt:     createdAt,
		FavoriteCount: r.Legacy.FavoriteCount,
		RetweetCount:  r.Legacy.RetweetCount,
		ReplyCount:    r.Legacy.ReplyCount,
		ViewCount:     viewCount,
		AuthorHandle:  handle,
		IsRetweet:     r.Legacy.RetweetedStatusIDStr != "",
		IsReply:       r.Legacy.InReplyToScreenName != "",
	}
	tweet.URL = tweet.TweetURL()

	// Quoted tweet.
	if r.QuotedStatusResult != nil && r.QuotedStatusResult.Result != nil {
		if qt := convertRawTweet(r.QuotedStatusResult.Result); qt != nil {
			tweet.QuotedTweet = qt
		}
	}

	return tweet
}
