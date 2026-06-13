package api

import (
	"fmt"

	"github.com/yashiels/twitter-cli/internal/types"
)

const (
	likesQueryID   = "M34xxhtrHWGAxGofYSGclA"
	likesOperation = "UserProfileFavoritesTimelineQuery"
)

// GetUserLikes fetches tweets liked by a user (by their numeric rest_id).
func (c *Client) GetUserLikes(userID string, limit int) ([]*types.Tweet, error) {
	if limit <= 0 {
		limit = 20
	}

	vars := map[string]any{"rest_id": userID, "count": limit}
	raw, err := c.graphqlGet(likesQueryID, likesOperation, vars, timelineFeatures)
	if err != nil {
		return nil, fmt.Errorf("GetUserLikes: %w", err)
	}

	// Response path mirrors UserProfileOriginalsTimelineQuery.
	instructionsRaw, err := getNestedJSON(raw, "data", "user_result", "result", "timeline_response", "timeline", "instructions")
	if err != nil {
		// Alternate path
		instructionsRaw, err = getNestedJSON(raw, "data", "user_result", "result", "timeline", "timeline", "instructions")
		if err != nil {
			return nil, fmt.Errorf("navigate likes response: %w", err)
		}
	}

	return parseTimelineInstructions(instructionsRaw)
}
