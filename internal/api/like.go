package api

import "fmt"

const (
	likeQueryID     = "awITBmMVajjvqY2wTL8DUw"
	likeOperation   = "FavoriteMutation"
	unlikeQueryID   = "SgiqtNaXXanog3v96NJVQA"
	unlikeOperation = "UnfavoriteMutation"
)

// likeFeatures are the feature flags for like/unlike mutations.
var likeFeatures = map[string]any{
	"subscriptions_feature_can_gift_premium": true,
}

// LikeTweet adds a like (favourite) to the tweet with the given ID.
func (c *Client) LikeTweet(tweetID string) error {
	vars := map[string]any{
		"tweet_id": tweetID,
	}

	_, err := c.graphqlPost(likeQueryID, likeOperation, vars, likeFeatures)
	if err != nil {
		return fmt.Errorf("LikeTweet: %w", err)
	}
	return nil
}

// UnlikeTweet removes a like (favourite) from the tweet with the given ID.
func (c *Client) UnlikeTweet(tweetID string) error {
	vars := map[string]any{
		"tweet_id": tweetID,
	}

	_, err := c.graphqlPost(unlikeQueryID, unlikeOperation, vars, likeFeatures)
	if err != nil {
		return fmt.Errorf("UnlikeTweet: %w", err)
	}
	return nil
}
