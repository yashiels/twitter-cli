package api

import (
	"encoding/json"
	"fmt"

	"github.com/yashiels/twitter-cli/internal/types"
)

const (
	searchQueryID   = "rxBGDmZrc-NcrXfcRNUdMg"
	searchOperation = "SearchTimelineQuery"
)

// searchFeatures are the feature flags for SearchTimelineQuery.
var searchFeatures = map[string]any{
	"android_home_timeline_status_injections_facepile_enabled": true,
	"grok_translations_post_auto_translation_is_enabled":       true,
	"reply_dislike_android_enabled":                            true,
	"subscriptions_feature_can_gift_premium":                   true,
	"x_lite_quick_promote_analytics_banner_enabled":            true,
}

// SearchResult holds either tweet results or user results depending on the search mode.
type SearchResult struct {
	Tweets []*types.Tweet
	Users  []*types.User
}

// SearchTimeline searches tweets or users.
// When searchUsers is true, product is "People" and results contain users.
func (c *Client) SearchTimeline(query string, limit int, searchUsers bool) (*SearchResult, error) {
	if limit <= 0 {
		limit = 20
	}

	product := "Top"
	if searchUsers {
		product = "People"
	}

	vars := map[string]any{
		"query":         query,
		"count":         limit,
		"query_source":  "typed_query",
		"timeline_type": product,
	}

	raw, err := c.graphqlGet(searchQueryID, searchOperation, vars, searchFeatures)
	if err != nil {
		return nil, fmt.Errorf("SearchTimeline: %w", err)
	}

	// Navigate: data -> search -> timeline_response -> timeline -> instructions (APK schema)
	instructionsRaw, err := getNestedJSON(raw, "data", "search", "timeline_response", "timeline", "instructions")
	if err != nil {
		return nil, fmt.Errorf("navigate search timeline: %w", err)
	}

	if searchUsers {
		users, err := parseUserTimelineInstructions(instructionsRaw)
		if err != nil {
			return nil, err
		}
		return &SearchResult{Users: users}, nil
	}

	tweets, err := parseTimelineInstructions(instructionsRaw)
	if err != nil {
		return nil, err
	}
	return &SearchResult{Tweets: tweets}, nil
}

// parseUserTimelineInstructions extracts users from a timeline instructions array.
// Used for People search results.
func parseUserTimelineInstructions(instructionsRaw json.RawMessage) ([]*types.User, error) {
	var instructions []json.RawMessage
	if err := json.Unmarshal(instructionsRaw, &instructions); err != nil {
		return nil, fmt.Errorf("parse instructions: %w", err)
	}

	var users []*types.User
	for _, instr := range instructions {
		var inst struct {
			Entries []json.RawMessage `json:"entries"`
		}
		var instAlt struct {
			Entries []json.RawMessage `json:"entries"`
		}
		_ = json.Unmarshal(instr, &inst)
		_ = json.Unmarshal(instr, &instAlt)

		entries := inst.Entries
		if len(entries) == 0 {
			entries = instAlt.Entries
		}

		for _, entry := range entries {
			u, err := parseUserTimelineEntry(entry)
			if err != nil || u == nil {
				continue
			}
			users = append(users, u)
		}
	}
	return users, nil
}

// parseUserTimelineEntry extracts a User from a timeline entry.
func parseUserTimelineEntry(entry json.RawMessage) (*types.User, error) {
	// Try content -> content -> user_results -> result
	resultRaw, err := getNestedJSON(entry, "content", "content", "user_results", "result")
	if err != nil {
		// Try content -> itemContent -> user_results -> result
		resultRaw, err = getNestedJSON(entry, "content", "itemContent", "user_results", "result")
		if err != nil {
			return nil, nil // not a user entry
		}
	}
	return parseUserResult(resultRaw)
}
