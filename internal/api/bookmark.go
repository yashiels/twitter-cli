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

var bookmarkFeatures = map[string]any{
	"subscriptions_feature_can_gift_premium": true,
}

// GetBookmarks fetches the authenticated user's bookmarks.
func (c *Client) GetBookmarks(limit int) ([]*types.Tweet, error) {
	if limit <= 0 {
		limit = 20
	}

	vars := map[string]any{"count": limit}
	// BookmarksTimelineQuery uses the same timeline features as home timeline.
	raw, err := c.graphqlGet(bookmarksQueryID, bookmarksOperation, vars, timelineFeatures)
	if err != nil {
		return nil, fmt.Errorf("GetBookmarks: %w", err)
	}

	// Try multiple response paths — APK response path needs live discovery.
	// APK schema uses camelCase "timelineResponse"
	instructionsRaw, err := getNestedJSON(raw, "data", "timelineResponse", "timeline", "instructions")
	if err != nil {
		instructionsRaw, err = getNestedJSON(raw, "data", "bookmarks_timeline", "timeline", "instructions")
		if err != nil {
			instructionsRaw, err = getNestedJSON(raw, "data", "timeline_response", "timeline", "instructions")
			if err != nil {
				return nil, fmt.Errorf("navigate bookmarks response: %w", err)
			}
		}
	}

	return parseTimelineInstructions(instructionsRaw)
}

// AddBookmark bookmarks a tweet.
func (c *Client) AddBookmark(tweetID string) error {
	vars := map[string]any{"tweet_id": tweetID}
	_, err := c.graphqlPost(bookmarkAddQueryID, bookmarkAddOperation, vars, bookmarkFeatures)
	if err != nil {
		return fmt.Errorf("AddBookmark: %w", err)
	}
	return nil
}

// RemoveBookmark removes a bookmark.
func (c *Client) RemoveBookmark(tweetID string) error {
	vars := map[string]any{"tweet_id": tweetID}
	_, err := c.graphqlPost(bookmarkRemQueryID, bookmarkRemOperation, vars, bookmarkFeatures)
	if err != nil {
		return fmt.Errorf("RemoveBookmark: %w", err)
	}
	return nil
}
