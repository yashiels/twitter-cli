package api

import "fmt"

const (
	blockQueryID     = "8zl3cVULtte29uCoWREtBQ"
	blockOperation   = "BlockUser"
	unblockQueryID   = "WtUZ-1fkiAJGXfN6gAwrLw"
	unblockOperation = "UnblockUser"
	muteQueryID      = "LoZAfbPr53jnw9Y2FydOIQ"
	muteOperation    = "MuteUser"
	unmuteQueryID    = "29vlsCe7kkuB4JKnQGeK5w"
	unmuteOperation  = "UnmuteUser"
)

// moderationFeatures are the feature flags for moderation mutations.
var moderationFeatures = map[string]any{
	"subscriptions_feature_can_gift_premium": true,
}

// BlockUser blocks the user with the given user ID.
func (c *Client) BlockUser(userID string) error {
	vars := map[string]any{
		"target_user_id": userID,
	}

	_, err := c.graphqlPost(blockQueryID, blockOperation, vars, moderationFeatures)
	if err != nil {
		return fmt.Errorf("BlockUser: %w", err)
	}
	return nil
}

// UnblockUser unblocks the user with the given user ID.
func (c *Client) UnblockUser(userID string) error {
	vars := map[string]any{
		"target_user_id": userID,
	}

	_, err := c.graphqlPost(unblockQueryID, unblockOperation, vars, moderationFeatures)
	if err != nil {
		return fmt.Errorf("UnblockUser: %w", err)
	}
	return nil
}

// MuteUser mutes the user with the given user ID.
func (c *Client) MuteUser(userID string) error {
	vars := map[string]any{
		"target_user_id": userID,
	}

	_, err := c.graphqlPost(muteQueryID, muteOperation, vars, moderationFeatures)
	if err != nil {
		return fmt.Errorf("MuteUser: %w", err)
	}
	return nil
}

// UnmuteUser unmutes the user with the given user ID.
func (c *Client) UnmuteUser(userID string) error {
	vars := map[string]any{
		"target_user_id": userID,
	}

	_, err := c.graphqlPost(unmuteQueryID, unmuteOperation, vars, moderationFeatures)
	if err != nil {
		return fmt.Errorf("UnmuteUser: %w", err)
	}
	return nil
}
