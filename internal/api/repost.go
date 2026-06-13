package api

import "fmt"

const (
	createRepostQueryID   = "ydMACa-dOjZx126SWo6q5A"
	createRepostOperation = "CreateRepostMutation"
	deleteRepostQueryID   = "w1Bo2Whh4f4lha5Djgnpvg"
	deleteRepostOperation = "DeleteRepostMutation"
)

var repostFeatures = map[string]any{
	"subscriptions_feature_can_gift_premium": true,
}

// CreateRepost reposts (retweets) a tweet.
func (c *Client) CreateRepost(tweetID string) error {
	vars := map[string]any{"tweet_id": tweetID}
	_, err := c.graphqlPost(createRepostQueryID, createRepostOperation, vars, repostFeatures)
	if err != nil {
		return fmt.Errorf("CreateRepost: %w", err)
	}
	return nil
}

// DeleteRepost removes a repost.
func (c *Client) DeleteRepost(tweetID string) error {
	vars := map[string]any{"tweet_id": tweetID}
	_, err := c.graphqlPost(deleteRepostQueryID, deleteRepostOperation, vars, repostFeatures)
	if err != nil {
		return fmt.Errorf("DeleteRepost: %w", err)
	}
	return nil
}
