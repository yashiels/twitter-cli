package api

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/yashiels/twitter-cli/internal/types"
)

// GetMentions fetches the authenticated user's mentions via the REST v1.1 endpoint.
// The GraphQL NotificationTimelineQuery has variable requirements that differ across APK
// versions; the REST endpoint is stable and well-suited for this purpose.
func (c *Client) GetMentions(limit int) ([]*types.Tweet, error) {
	if limit <= 0 {
		limit = 20
	}

	params := url.Values{
		"count":            {strconv.Itoa(limit)},
		"tweet_mode":       {"extended"},
		"include_entities": {"true"},
	}

	raw, err := c.restGet("/1.1/statuses/mentions_timeline.json", params)
	if err != nil {
		return nil, fmt.Errorf("GetMentions: %w", err)
	}

	return parseRESTTweetList(raw)
}
