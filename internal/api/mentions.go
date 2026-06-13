package api

import (
	"fmt"

	"github.com/yashiels/twitter-cli/internal/types"
)

const (
	mentionsQueryID   = "jsOzc8RhpUpH5InTskP6Yw"
	mentionsOperation = "NotificationTimelineQuery"
)

// GetMentions fetches the authenticated user's notification timeline.
// userID is the numeric rest_id stored in credentials.
// Variable name "userId" (camelCase) is verified against the decompiled APK source.
func (c *Client) GetMentions(userID string, limit int) ([]*types.Tweet, error) {
	if limit <= 0 {
		limit = 20
	}

	vars := map[string]any{"userId": userID, "count": limit, "notificationTimelineType": "All"}
	raw, err := c.graphqlGet(mentionsQueryID, mentionsOperation, vars, timelineFeatures)
	if err != nil {
		return nil, fmt.Errorf("GetMentions: %w", err)
	}

	// Try multiple response paths — needs live discovery.
	// APK path: data -> user_result -> result -> timeline_response -> timeline -> instructions
	instructionsRaw, err := getNestedJSON(raw, "data", "user_result", "result", "timeline_response", "timeline", "instructions")
	if err != nil {
		instructionsRaw, err = getNestedJSON(raw, "data", "notification_timeline", "timeline", "instructions")
		if err != nil {
			return nil, fmt.Errorf("navigate mentions response: %w", err)
		}
	}

	// parseTimelineInstructions skips non-tweet entries automatically.
	return parseTimelineInstructions(instructionsRaw)
}
