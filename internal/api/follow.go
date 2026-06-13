package api

import "fmt"

const (
	followQueryID     = "44lRL9CTLTxi4aAMSqAmVw"
	followOperation   = "FollowUser"
	unfollowQueryID   = "zpWrwHHfa_6sKBQr6SGCwg"
	unfollowOperation = "UnfollowUser"
)

// followFeatures are the feature flags for follow/unfollow mutations.
var followFeatures = map[string]any{
	"subscriptions_feature_can_gift_premium": true,
}

// FollowUser follows the user with the given rest_id.
func (c *Client) FollowUser(userID string) error {
	vars := map[string]any{
		"rest_id": userID,
	}

	_, err := c.graphqlPost(followQueryID, followOperation, vars, followFeatures)
	if err != nil {
		return fmt.Errorf("FollowUser: %w", err)
	}
	return nil
}

// UnfollowUser unfollows the user with the given rest_id.
func (c *Client) UnfollowUser(userID string) error {
	vars := map[string]any{
		"rest_id": userID,
	}

	_, err := c.graphqlPost(unfollowQueryID, unfollowOperation, vars, followFeatures)
	if err != nil {
		return fmt.Errorf("UnfollowUser: %w", err)
	}
	return nil
}
