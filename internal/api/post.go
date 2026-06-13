package api

import (
	"encoding/json"
	"fmt"

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

// parseCreatePostResponse parses the CreatePost/CreateReply/CreateQuote API response.
// The current APK returns tweet_results as a flat object with rest_id directly (no "result" wrapper).
// We also try the classic nested path for forward-compat.
func parseCreatePostResponse(raw json.RawMessage) (*types.Tweet, error) {
	// Classic nested path: data -> create_tweet -> tweet_results -> result
	resultRaw, err := getNestedJSON(raw, "data", "create_tweet", "tweet_results", "result")
	if err == nil {
		return parseTweetResult(resultRaw)
	}

	// New APK format: tweet_results is flat — {"__typename":"TweetResults","rest_id":"..."}
	tweetResultsRaw, navErr := getNestedJSON(raw, "data", "create_tweet", "tweet_results")
	if navErr != nil {
		// Also try create_post key variant.
		tweetResultsRaw, navErr = getNestedJSON(raw, "data", "create_post", "tweet_results")
		if navErr != nil {
			return nil, fmt.Errorf("cannot navigate CreatePost response: %w", err)
		}
	}

	// Extract rest_id from the flat TweetResults object.
	var flat struct {
		RestID string `json:"rest_id"`
	}
	if jsonErr := json.Unmarshal(tweetResultsRaw, &flat); jsonErr != nil {
		return nil, fmt.Errorf("parse tweet_results: %w", jsonErr)
	}
	if flat.RestID == "" {
		return nil, fmt.Errorf("empty rest_id in CreatePost response")
	}

	// Return a minimal tweet; caller adds the author handle via stored credentials.
	return &types.Tweet{
		ID:  flat.RestID,
		URL: "https://x.com/i/status/" + flat.RestID,
	}, nil
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

	t, err := parseCreatePostResponse(raw)
	if err != nil {
		return nil, fmt.Errorf("CreatePost: %w", err)
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

	t, err := parseCreatePostResponse(raw)
	if err != nil {
		return nil, fmt.Errorf("CreateReply: %w", err)
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

	t, err := parseCreatePostResponse(raw)
	if err != nil {
		return nil, fmt.Errorf("CreateQuote: %w", err)
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
