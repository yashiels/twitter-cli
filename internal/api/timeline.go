package api

import (
	"fmt"

	"github.com/yashiels/twitter-cli/internal/types"
)

const (
	homeTimelineQueryID   = "t_sH369wuH1CO5lbW2qlYg"
	homeTimelineOperation = "HomeTimeline"
	homeLatestQueryID     = "xFpaluPXjdeLR4A29HRmXw"
	homeLatestOperation   = "HomeTimelineLatest"
)

// timelineFeatures are the feature flags for home timeline queries.
var timelineFeatures = map[string]any{
	"android_home_timeline_status_injections_facepile_enabled": true,
	"grok_translations_post_auto_translation_is_enabled":       true,
	"reply_dislike_android_enabled":                            true,
	"subscriptions_feature_can_gift_premium":                   true,
	"x_lite_quick_promote_analytics_banner_enabled":            true,
}

// GetHomeTimeline fetches the authenticated user's home timeline.
// When latest is true it uses the "Following" / Latest tab (HomeTimelineLatest).
func (c *Client) GetHomeTimeline(limit int, latest bool) ([]*types.Tweet, error) {
	if limit <= 0 {
		limit = 20
	}

	queryID := homeTimelineQueryID
	operation := homeTimelineOperation
	if latest {
		queryID = homeLatestQueryID
		operation = homeLatestOperation
	}

	vars := map[string]any{
		"count": limit,
	}

	raw, err := c.graphqlGet(queryID, operation, vars, timelineFeatures)
	if err != nil {
		return nil, fmt.Errorf("GetHomeTimeline: %w", err)
	}

	// Navigate: data -> timeline_response -> timeline -> instructions (APK schema)
	instructionsRaw, err := getNestedJSON(raw, "data", "timeline_response", "timeline", "instructions")
	if err != nil {
		return nil, fmt.Errorf("navigate home timeline: %w", err)
	}

	return parseTimelineInstructions(instructionsRaw)
}
