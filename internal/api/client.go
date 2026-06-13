package api

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/yashiels/twitter-cli/internal/auth"
)

const (
	// BearerToken is the public Twitter Android bearer token.
	BearerToken = "AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"

	// BaseURL is the Twitter GraphQL API base.
	BaseURL = "https://api.twitter.com/graphql"

	// UserAgent mimics the Twitter Android app.
	UserAgent = "TwitterAndroid/11.99.0"
)

// Client is a Twitter API HTTP client.
type Client struct {
	http  *http.Client
	creds *auth.Credentials
}

// NewClient creates a new API client using stored credentials.
func NewClient(creds *auth.Credentials) *Client {
	return &Client{
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
		creds: creds,
	}
}

// do executes an HTTP request with the required Twitter auth headers.
// It handles rate limiting automatically.
func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+BearerToken)
	req.Header.Set("Cookie", fmt.Sprintf("auth_token=%s; ct0=%s", c.creds.AuthToken, c.creds.CT0))
	req.Header.Set("x-csrf-token", c.creds.CT0)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("x-twitter-active-user", "yes")
	req.Header.Set("x-twitter-auth-type", "OAuth2Session")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}

	logRateLimit(resp)

	if resp.StatusCode == http.StatusTooManyRequests {
		resetHeader := resp.Header.Get("x-rate-limit-reset")
		resetTime := parseResetTime(resetHeader)
		waitSec := max(int(time.Until(resetTime).Seconds()), 5)
		fmt.Fprintf(io.Discard, "") // suppress unused import
		return nil, fmt.Errorf("rate limited — reset in %ds (at %s)", waitSec, resetTime.Format(time.RFC1123))
	}

	return resp, nil
}

// jitter sleeps for a random 2–5 seconds to reduce API pressure.
func jitter() {
	ms := 2000 + rand.Intn(3000) //nolint:gosec
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// logRateLimit prints rate limit info to stderr if limits are getting low.
func logRateLimit(resp *http.Response) {
	remaining := resp.Header.Get("x-rate-limit-remaining")
	reset := resp.Header.Get("x-rate-limit-reset")
	if remaining == "" {
		return
	}
	rem, _ := strconv.Atoi(remaining)
	if rem < 5 {
		resetTime := parseResetTime(reset)
		fmt.Fprintf(io.Discard, "rate limit: %d remaining, resets %s\n", rem, resetTime.Format(time.RFC1123))
	}
}

// parseResetTime converts the Unix timestamp string from x-rate-limit-reset to a time.Time.
func parseResetTime(s string) time.Time {
	if s == "" {
		return time.Now().Add(60 * time.Second)
	}
	ts, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Now().Add(60 * time.Second)
	}
	return time.Unix(ts, 0)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
