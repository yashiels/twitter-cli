package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// VerifyCredentials calls the v1.1 account/verify_credentials endpoint
// to confirm auth is working and return the current user's screen name.
func (c *Client) VerifyCredentials() (string, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.twitter.com/1.1/account/verify_credentials.json",
		nil,
	)
	if err != nil {
		return "", err
	}

	resp, err := c.do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("unauthorized — check auth_token and ct0")
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("verify_credentials: HTTP %d", resp.StatusCode)
	}

	var result struct {
		ScreenName string `json:"screen_name"`
		Name       string `json:"name"`
		ID         int64  `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode verify response: %w", err)
	}
	return result.ScreenName, nil
}
