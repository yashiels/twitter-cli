package api

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yashiels/twitter-cli/internal/types"
)

const (
	bookmarksQueryID     = "DN0j17CihaEo7QYmaZGkiw"
	bookmarksOperation   = "BookmarksTimelineQuery"
	bookmarkAddQueryID   = "IjefskW4Kr2i-6XRdshQEg"
	bookmarkAddOperation = "BookmarkAddMutation"
	bookmarkRemQueryID   = "K5KIqVnds5iJ00WdHb8Nmw"
	bookmarkRemOperation = "BookmarkRemoveMutation"
)

// bookmarkFeatures are the feature flags for bookmark operations.
var bookmarkFeatures = map[string]any{
	"subscriptions_feature_can_gift_premium": true,
}

// GetBookmarks fetches the authenticated user's bookmarks.
func (c *Client) GetBookmarks(limit int) ([]*types.Tweet, error) {
	if limit <= 0 {
		limit = 20
	}

	vars := map[string]any{
		"count": limit,
	}

	raw, err := c.graphqlGet(bookmarksQueryID, bookmarksOperation, vars, timelineFeatures)
	if err != nil {
		return nil, fmt.Errorf("GetBookmarks: %w", err)
	}

	if os.Getenv("DEBUG_TWT") != "" {
		var pretty any
		_ = json.Unmarshal(raw, &pretty)
		prettyJSON, _ := json.MarshalIndent(pretty, "", "  ")
		fmt.Fprintf(os.Stderr, "DEBUG GetBookmarks response:\n%s\n", prettyJSON)
	}

	// Try multiple response paths — bookmarks timeline path varies across APK versions.
	paths := [][]string{
		{"data", "timelineResponse", "timeline", "instructions"}, // current APK (camelCase)
		{"data", "bookmarks_timeline", "timeline", "instructions"},
		{"data", "bookmark_timeline_v2", "timeline", "instructions"},
		{"data", "timeline_response", "timeline", "instructions"},
	}

	for _, path := range paths {
		instructionsRaw, pathErr := getNestedJSON(raw, path...)
		if pathErr == nil {
			return parseTimelineInstructions(instructionsRaw)
		}
	}

	return nil, fmt.Errorf("GetBookmarks: cannot navigate response — unknown bookmarks timeline path")
}

// AddBookmark bookmarks a tweet.
func (c *Client) AddBookmark(tweetID string) error {
	vars := map[string]any{
		"tweet_id": tweetID,
	}

	_, err := c.graphqlPost(bookmarkAddQueryID, bookmarkAddOperation, vars, bookmarkFeatures)
	if err != nil {
		return fmt.Errorf("AddBookmark: %w", err)
	}
	return nil
}

// RemoveBookmark removes a bookmark.
func (c *Client) RemoveBookmark(tweetID string) error {
	vars := map[string]any{
		"tweet_id": tweetID,
	}

	_, err := c.graphqlPost(bookmarkRemQueryID, bookmarkRemOperation, vars, bookmarkFeatures)
	if err != nil {
		return fmt.Errorf("RemoveBookmark: %w", err)
	}
	return nil
}
