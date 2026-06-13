package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/yashiels/twitter-cli/internal/types"
)

// GetFollowers fetches a user's followers via REST v1.1.
// REST v1.1 endpoint: /1.1/followers/list.json
func (c *Client) GetFollowers(handle string, limit int) ([]*types.User, error) {
	if limit <= 0 {
		limit = 20
	}

	params := url.Values{
		"screen_name": {handle},
		"count":       {strconv.Itoa(limit)},
		"skip_status": {"true"},
	}
	raw, err := c.restGet("/1.1/followers/list.json", params)
	if err != nil {
		return nil, fmt.Errorf("GetFollowers: %w", err)
	}
	return parseRESTUsersResponse(raw)
}

// GetFollowing fetches who a user follows via REST v1.1.
// REST v1.1 endpoint: /1.1/friends/list.json
func (c *Client) GetFollowing(handle string, limit int) ([]*types.User, error) {
	if limit <= 0 {
		limit = 20
	}

	params := url.Values{
		"screen_name": {handle},
		"count":       {strconv.Itoa(limit)},
		"skip_status": {"true"},
	}
	raw, err := c.restGet("/1.1/friends/list.json", params)
	if err != nil {
		return nil, fmt.Errorf("GetFollowing: %w", err)
	}
	return parseRESTUsersResponse(raw)
}

// parseRESTUsersResponse parses a REST v1.1 users list response.
// Response shape: {"users": [...], "next_cursor_str": "..."}
func parseRESTUsersResponse(raw json.RawMessage) ([]*types.User, error) {
	var resp struct {
		Users []json.RawMessage `json:"users"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse REST users: %w", err)
	}

	users := make([]*types.User, 0, len(resp.Users))
	for _, u := range resp.Users {
		user, err := parseRESTUser(u)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

// restUserRaw mirrors a REST v1.1 user object.
// REST v1.1 uses flat fields (no nesting), unlike the GraphQL APK schema.
type restUserRaw struct {
	IDStr          string `json:"id_str"`
	ScreenName     string `json:"screen_name"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Location       string `json:"location"`
	FollowersCount int    `json:"followers_count"`
	FriendsCount   int    `json:"friends_count"`
	StatusesCount  int    `json:"statuses_count"`
	Verified       bool   `json:"verified"`
}

// parseRESTUser converts a REST v1.1 user JSON object to a types.User.
func parseRESTUser(raw json.RawMessage) (*types.User, error) {
	var u restUserRaw
	if err := json.Unmarshal(raw, &u); err != nil {
		return nil, err
	}
	if u.ScreenName == "" {
		return nil, fmt.Errorf("REST user missing screen_name")
	}
	return &types.User{
		RestID:         u.IDStr,
		ScreenName:     u.ScreenName,
		Name:           u.Name,
		Description:    u.Description,
		Location:       u.Location,
		FollowersCount: u.FollowersCount,
		FriendsCount:   u.FriendsCount,
		StatusesCount:  u.StatusesCount,
		Verified:       u.Verified,
	}, nil
}
