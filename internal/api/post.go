package api

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yashiels/twitter-cli/internal/types"
)

const (
	createPostQueryID   = "vMia9QJ2JVkCXuO5J4MTbw"
	createPostOperation = "CreatePost"
	deletePostQueryID   = "1EVIme6zMCgTO7F95wuElA"
	deletePostOperation = "DeletePostMutation"
)

// postFeatures are the feature flags for post mutations.
var postFeatures = map[string]any{
	"subscriptions_feature_can_gift_premium": true,
}

// CreatePost creates a new tweet with the given text.
func (c *Client) CreatePost(text string) (*types.Tweet, error) {
	vars := map[string]any{
		"post_text": text,
	}

	raw, err := c.graphqlPost(createPostQueryID, createPostOperation, vars, postFeatures)
	if err != nil {
		return nil, fmt.Errorf("CreatePost: %w", err)
	}

	// Debug: print raw response if DEBUG_TWT is set.
	if os.Getenv("DEBUG_TWT") != "" {
		var pretty any
		_ = json.Unmarshal(raw, &pretty)
		prettyJSON, _ := json.MarshalIndent(pretty, "", "  ")
		fmt.Fprintf(os.Stderr, "DEBUG CreatePost response:\n%s\n", prettyJSON)
	}

	// Response path: data -> create_tweet -> tweet_results -> result
	// Fallback: data -> create_post -> tweet_results -> result (some APK versions)
	resultRaw, err := getNestedJSON(raw, "data", "create_tweet", "tweet_results", "result")
	if err != nil {
		resultRaw, err = getNestedJSON(raw, "data", "create_post", "tweet_results", "result")
		if err != nil {
			return nil, fmt.Errorf("CreatePost: navigate response: %w", err)
		}
	}

	t, err := parseTweetResult(resultRaw)
	if err != nil {
		return nil, fmt.Errorf("CreatePost: parse tweet: %w", err)
	}
	return t, nil
}

// CreateReply creates a reply to an existing tweet.
func (c *Client) CreateReply(text, replyToID string) (*types.Tweet, error) {
	vars := map[string]any{
		"post_text": text,
		"reply": map[string]any{
			"in_reply_to_tweet_id":   replyToID,
			"exclude_reply_user_ids": []string{},
		},
	}

	raw, err := c.graphqlPost(createPostQueryID, createPostOperation, vars, postFeatures)
	if err != nil {
		return nil, fmt.Errorf("CreateReply: %w", err)
	}

	// Same path as CreatePost — fallback for alternate APK response key.
	resultRaw, err := getNestedJSON(raw, "data", "create_tweet", "tweet_results", "result")
	if err != nil {
		resultRaw, err = getNestedJSON(raw, "data", "create_post", "tweet_results", "result")
		if err != nil {
			return nil, fmt.Errorf("CreateReply: navigate response: %w", err)
		}
	}

	t, err := parseTweetResult(resultRaw)
	if err != nil {
		return nil, fmt.Errorf("CreateReply: parse tweet: %w", err)
	}
	return t, nil
}

// CreateQuote creates a quote tweet.
func (c *Client) CreateQuote(text, quotedTweetID string) (*types.Tweet, error) {
	vars := map[string]any{
		"post_text":      text,
		"attachment_url": "https://x.com/i/status/" + quotedTweetID,
	}

	raw, err := c.graphqlPost(createPostQueryID, createPostOperation, vars, postFeatures)
	if err != nil {
		return nil, fmt.Errorf("CreateQuote: %w", err)
	}

	// Same path as CreatePost — fallback for alternate APK response key.
	resultRaw, err := getNestedJSON(raw, "data", "create_tweet", "tweet_results", "result")
	if err != nil {
		resultRaw, err = getNestedJSON(raw, "data", "create_post", "tweet_results", "result")
		if err != nil {
			return nil, fmt.Errorf("CreateQuote: navigate response: %w", err)
		}
	}

	t, err := parseTweetResult(resultRaw)
	if err != nil {
		return nil, fmt.Errorf("CreateQuote: parse tweet: %w", err)
	}
	return t, nil
}

// DeletePost deletes a tweet by ID.
func (c *Client) DeletePost(tweetID string) error {
	vars := map[string]any{
		"tweet_id": tweetID,
	}

	_, err := c.graphqlPost(deletePostQueryID, deletePostOperation, vars, postFeatures)
	if err != nil {
		return fmt.Errorf("DeletePost: %w", err)
	}
	return nil
}
