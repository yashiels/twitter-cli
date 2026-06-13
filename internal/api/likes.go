package api

import (
	"fmt"

	"github.com/yashiels/twitter-cli/internal/types"
)

const (
	likesQueryID   = "M34xxhtrHWGAxGofYSGclA"
	likesOperation = "UserProfileFavoritesTimelineQuery"
)

// GetUserLikes fetches tweets liked by a user (identified by numeric user ID).
func (c *Client) GetUserLikes(userID string, limit int) ([]*types.Tweet, error) {
	if limit <= 0 {
		limit = 20
	}

	vars := map[string]any{
		"rest_id": userID,
		"count":   limit,
	}

	raw, err := c.graphqlGet(likesQueryID, likesOperation, vars, timelineFeatures)
	if err != nil {
		return nil, fmt.Errorf("GetUserLikes: %w", err)
	}

	// Try response paths — likes timeline mirrors tweets timeline structure.
	paths := [][]string{
		{"data", "user_result", "result", "timeline_response", "timeline", "instructions"},
		{"data", "user_result", "result", "timeline", "timeline", "instructions"},
	}

	for _, path := range paths {
		instructionsRaw, err := getNestedJSON(raw, path...)
		if err == nil {
			return parseTimelineInstructions(instructionsRaw)
		}
	}

	return nil, fmt.Errorf("GetUserLikes: cannot navigate response — unknown likes timeline path")
}
