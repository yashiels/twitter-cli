package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// VerifyResult holds the current user's identity from verify_credentials.
type VerifyResult struct {
	ScreenName string
	UserID     string
}

// VerifyCredentials calls the v1.1 account/verify_credentials endpoint
// to confirm auth is working and return the current user's identity.
func (c *Client) VerifyCredentials() (*VerifyResult, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.twitter.com/1.1/account/verify_credentials.json",
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized — check auth_token and ct0")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("verify_credentials: HTTP %d", resp.StatusCode)
	}

	var result struct {
		ScreenName string `json:"screen_name"`
		Name       string `json:"name"`
		IDStr      string `json:"id_str"`
		ID         int64  `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode verify response: %w", err)
	}

	userID := result.IDStr
	if userID == "" && result.ID != 0 {
		userID = strconv.FormatInt(result.ID, 10)
	}

	return &VerifyResult{
		ScreenName: result.ScreenName,
		UserID:     userID,
	}, nil
}
