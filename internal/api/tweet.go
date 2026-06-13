package api

import (
	"fmt"

	"github.com/yashiels/twitter-cli/internal/types"
)

const (
	tweetQueryID   = "lOsezlo57Y40B-TLgWqxEA"
	tweetOperation = "GetPostById"
)

// tweetFeatures are the feature flags for GetPostById.
var tweetFeatures = map[string]any{
	"grok_translations_post_auto_translation_is_enabled": true,
	"reply_dislike_android_enabled":                      true,
	"subscriptions_feature_can_gift_premium":             true,
	"x_lite_quick_promote_analytics_banner_enabled":      true,
}

// GetTweetByID fetches a single tweet by its ID.
func (c *Client) GetTweetByID(tweetID string) (*types.Tweet, error) {
	vars := map[string]any{
		"rest_id": tweetID,
	}

	raw, err := c.graphqlGet(tweetQueryID, tweetOperation, vars, tweetFeatures)
	if err != nil {
		return nil, fmt.Errorf("GetTweetByID: %w", err)
	}

	// Navigate: data -> tweet_result_by_rest_id -> result (APK schema)
	resultRaw, err := getNestedJSON(raw, "data", "tweet_result_by_rest_id", "result")
	if err != nil {
		return nil, fmt.Errorf("navigate tweet result: %w", err)
	}

	t, err := parseTweetResult(resultRaw)
	if err != nil {
		return nil, fmt.Errorf("parse tweet: %w", err)
	}
	if t == nil {
		return nil, fmt.Errorf("tweet %s not found", tweetID)
	}
	return t, nil
}
