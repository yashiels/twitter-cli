package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// VerifyCredentials calls the v1.1 account/verify_credentials endpoint
// to confirm auth is working and return the current user's screen name.
func (c *Client) VerifyCredentials() (string, error) {
	handle, _, err := c.VerifyCredentialsWithID()
	return handle, err
}

// VerifyCredentialsWithID calls the v1.1 account/verify_credentials endpoint
// and returns both the screen_name and the numeric id_str.
func (c *Client) VerifyCredentialsWithID() (handle, userID string, err error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.twitter.com/1.1/account/verify_credentials.json",
		nil,
	)
	if err != nil {
		return "", "", err
	}

	resp, err := c.do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", "", fmt.Errorf("unauthorized — check auth_token and ct0")
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("verify_credentials: HTTP %d", resp.StatusCode)
	}

	var result struct {
		ScreenName string `json:"screen_name"`
		Name       string `json:"name"`
		ID         int64  `json:"id"`
		IDStr      string `json:"id_str"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", fmt.Errorf("decode verify response: %w", err)
	}

	idStr := result.IDStr
	if idStr == "" && result.ID != 0 {
		idStr = strconv.FormatInt(result.ID, 10)
	}

	return result.ScreenName, idStr, nil
}
