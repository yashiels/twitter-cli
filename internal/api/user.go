package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/yashiels/twitter-cli/internal/types"
)

const (
	getUserQueryID   = "bbS0COK9SwcgdM7QCEqWDg"
	getUserOperation = "GetUserByScreenNameQuery"

	getUserByIDQueryID   = "q5op2HlD7t5RvWB_SrFVfQ"
	getUserByIDOperation = "GetUserQuery"
)

// getUserFeatures are the feature flags required by GetUserByScreenNameQuery.
var getUserFeatures = map[string]any{
	"subscriptions_feature_can_gift_premium": true,
}

// ErrUserNotFound is returned when a user does not exist.
var ErrUserNotFound = errors.New("user not found")

// GetUserByID fetches a user by numeric rest_id.
func (c *Client) GetUserByID(userID string) (*types.User, error) {
	vars := map[string]any{
		"rest_id": userID,
	}

	raw, err := c.graphqlGet(getUserByIDQueryID, getUserByIDOperation, vars, getUserFeatures)
	if err != nil {
		return nil, fmt.Errorf("GetUserByID: %w", err)
	}

	// Path: data -> user_result -> result
	result, err := getNestedJSON(raw, "data", "user_result", "result")
	if err != nil {
		return nil, ErrUserNotFound
	}
	return parseUserResult(result)
}

// GetUserByScreenName fetches a Twitter user's profile by their @handle.
func (c *Client) GetUserByScreenName(handle string) (*types.User, error) {
	vars := map[string]any{
		"screen_name": handle,
	}

	raw, err := c.graphqlGet(getUserQueryID, getUserOperation, vars, getUserFeatures)
	if err != nil {
		return nil, fmt.Errorf("GetUserByScreenName: %w", err)
	}

	return parseUserResponse(raw)
}

// parseUserResponse parses the APK-schema user response.
// The APK schema is flat (not nested under "legacy" like the web API).
func parseUserResponse(raw json.RawMessage) (*types.User, error) {
	// Path: data.user_result.result
	result, err := getNestedJSON(raw, "data", "user_result", "result")
	if err != nil {
		return nil, ErrUserNotFound
	}
	return parseUserResult(result)
}

// parseUserResult converts a raw user result JSON object into a User.
// This is called both from parseUserResponse and from timeline entry parsers.
func parseUserResult(result json.RawMessage) (*types.User, error) {
	// The result object has typename
	var typed struct {
		Typename string `json:"__typename"`
	}
	_ = json.Unmarshal(result, &typed)
	if typed.Typename == "UserUnavailable" {
		return nil, ErrUserNotFound
	}

	// APK flat schema — try flat first.
	var flat struct {
		RestID string `json:"rest_id"`
		Core   struct {
			ScreenName string `json:"screen_name"`
			Name       string `json:"name"`
		} `json:"core"`
		Location struct {
			Location string `json:"location"`
		} `json:"location"`
		ProfileBio struct {
			Description string `json:"description"`
		} `json:"profile_bio"`
		// APK schema: relationship_counts and tweet_counts
		RelationshipCounts struct {
			Followers int `json:"followers"`
			Following int `json:"following"`
		} `json:"relationship_counts"`
		TweetCounts struct {
			Tweets int `json:"tweets"`
		} `json:"tweet_counts"`
		Verification struct {
			IsBlueVerified bool `json:"is_blue_verified"`
		} `json:"verification"`
		AffiliatesHighlightedLabel struct {
			Label struct {
				Description string `json:"description"`
			} `json:"label"`
		} `json:"affiliates_highlighted_label"`
		// Legacy fallback (web schema)
		IsBlueVerified bool `json:"is_blue_verified"`
		Legacy         *struct {
			ScreenName     string `json:"screen_name"`
			Name           string `json:"name"`
			Description    string `json:"description"`
			Location       string `json:"location"`
			FollowersCount int    `json:"followers_count"`
			FriendsCount   int    `json:"friends_count"`
			StatusesCount  int    `json:"statuses_count"`
		} `json:"legacy"`
	}

	if err := json.Unmarshal(result, &flat); err != nil {
		return nil, fmt.Errorf("parse user result: %w", err)
	}

	user := &types.User{
		RestID:      flat.RestID,
		Verified:    flat.Verification.IsBlueVerified || flat.IsBlueVerified,
		Affiliation: flat.AffiliatesHighlightedLabel.Label.Description,
	}

	// Prefer APK flat schema fields, fall back to legacy.
	if flat.Core.ScreenName != "" {
		user.ScreenName = flat.Core.ScreenName
		user.Name = flat.Core.Name
	}
	if flat.Location.Location != "" {
		user.Location = flat.Location.Location
	}
	if flat.ProfileBio.Description != "" {
		user.Description = flat.ProfileBio.Description
	}
	// APK schema: relationship_counts + tweet_counts
	user.FollowersCount = flat.RelationshipCounts.Followers
	user.FriendsCount = flat.RelationshipCounts.Following
	user.StatusesCount = flat.TweetCounts.Tweets

	// Fallback: web API legacy schema
	if user.ScreenName == "" && flat.Legacy != nil {
		user.ScreenName = flat.Legacy.ScreenName
		user.Name = flat.Legacy.Name
		user.Description = flat.Legacy.Description
		user.Location = flat.Legacy.Location
		user.FollowersCount = flat.Legacy.FollowersCount
		user.FriendsCount = flat.Legacy.FriendsCount
		user.StatusesCount = flat.Legacy.StatusesCount
	}

	if user.RestID == "" && user.ScreenName == "" {
		return nil, ErrUserNotFound
	}

	return user, nil
}
