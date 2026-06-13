package api

import (
	"fmt"

	"github.com/yashiels/twitter-cli/internal/types"
)

const (
	createPostQueryID   = "vMia9QJ2JVkCXuO5J4MTbw"
	createPostOperation = "CreatePost"
	deletePostQueryID   = "1EVIme6zMCgTO7F95wuElA"
	deletePostOperation = "DeletePostMutation"
)

var postFeatures = map[string]any{
	"subscriptions_feature_can_gift_premium": true,
}

// CreatePost creates a new tweet.
// Variable name "post_text" is verified against the decompiled APK source.
func (c *Client) CreatePost(text string) (*types.Tweet, error) {
	vars := map[string]any{"post_text": text}
	raw, err := c.graphqlPost(createPostQueryID, createPostOperation, vars, postFeatures)
	if err != nil {
		return nil, fmt.Errorf("CreatePost: %w", err)
	}
	// Response path: data -> create_tweet -> tweet_results -> result
	result, err := getNestedJSON(raw, "data", "create_tweet", "tweet_results", "result")
	if err != nil {
		// Alternate path in some APK versions
		result, err = getNestedJSON(raw, "data", "create_post", "tweet_results", "result")
		if err != nil {
			return nil, fmt.Errorf("navigate create_post response: %w", err)
		}
	}
	return parseTweetResult(result)
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
	result, err := getNestedJSON(raw, "data", "create_tweet", "tweet_results", "result")
	if err != nil {
		result, err = getNestedJSON(raw, "data", "create_post", "tweet_results", "result")
		if err != nil {
			return nil, fmt.Errorf("navigate create_reply response: %w", err)
		}
	}
	return parseTweetResult(result)
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
	result, err := getNestedJSON(raw, "data", "create_tweet", "tweet_results", "result")
	if err != nil {
		result, err = getNestedJSON(raw, "data", "create_post", "tweet_results", "result")
		if err != nil {
			return nil, fmt.Errorf("navigate create_quote response: %w", err)
		}
	}
	return parseTweetResult(result)
}

// DeletePost deletes a tweet by ID.
func (c *Client) DeletePost(tweetID string) error {
	vars := map[string]any{"tweet_id": tweetID}
	_, err := c.graphqlPost(deletePostQueryID, deletePostOperation, vars, postFeatures)
	if err != nil {
		return fmt.Errorf("DeletePost: %w", err)
	}
	return nil
}
