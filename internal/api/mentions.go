package api

import (
	"fmt"

	"github.com/yashiels/twitter-cli/internal/types"
)

const (
	mentionsQueryID   = "jsOzc8RhpUpH5InTskP6Yw"
	mentionsOperation = "NotificationTimelineQuery"
)

// GetMentions fetches the user's notification timeline (mentions and other notifications).
func (c *Client) GetMentions(userID string, limit int) ([]*types.Tweet, error) {
	if limit <= 0 {
		limit = 20
	}

	vars := map[string]any{
		"userId": userID,
		"count":  limit,
	}

	raw, err := c.graphqlGet(mentionsQueryID, mentionsOperation, vars, timelineFeatures)
	if err != nil {
		return nil, fmt.Errorf("GetMentions: %w", err)
	}

	// Try multiple response paths — notification timeline path needs live discovery.
	paths := [][]string{
		{"data", "notification_timeline", "timeline", "instructions"},
		{"data", "timeline_response", "timeline", "instructions"},
		{"data", "notifications_timeline", "timeline", "instructions"},
	}

	for _, path := range paths {
		instructionsRaw, err := getNestedJSON(raw, path...)
		if err == nil {
			return parseTimelineInstructions(instructionsRaw)
		}
	}

	return nil, fmt.Errorf("GetMentions: cannot navigate response — unknown notification timeline path")
}
