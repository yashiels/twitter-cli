package api

import (
	"fmt"

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

	// Try multiple response paths — bookmarks timeline path varies.
	paths := [][]string{
		{"data", "bookmarks_timeline", "timeline", "instructions"},
		{"data", "bookmark_timeline_v2", "timeline", "instructions"},
		{"data", "timeline_response", "timeline", "instructions"},
	}

	for _, path := range paths {
		instructionsRaw, err := getNestedJSON(raw, path...)
		if err == nil {
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
