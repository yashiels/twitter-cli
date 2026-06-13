package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
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
// It adds jitter between requests and handles rate limiting with auto-retry.
func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+BearerToken)
	req.Header.Set("Cookie", fmt.Sprintf("auth_token=%s; ct0=%s", c.creds.AuthToken, c.creds.CT0))
	req.Header.Set("x-csrf-token", c.creds.CT0)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("x-twitter-active-user", "yes")
	req.Header.Set("x-twitter-auth-type", "OAuth2Session")

	// Random jitter before each request to mimic human behavior.
	jitter()

	var resp *http.Response
	var err error

	// Retry up to 2 times on rate limit (429).
	// Capture GetBody so POST bodies can be reset between retries.
	getBody := req.GetBody
	for attempt := 0; attempt < 3; attempt++ {
		// Reset body on retries — after the first Do() the transport drains it.
		if attempt > 0 && getBody != nil {
			newBody, bodyErr := getBody()
			if bodyErr != nil {
				return nil, fmt.Errorf("reset request body: %w", bodyErr)
			}
			req.Body = newBody
		}
		resp, err = c.http.Do(req)
		if err != nil {
			return nil, fmt.Errorf("http request: %w", err)
		}

		logRateLimit(resp)

		if resp.StatusCode != http.StatusTooManyRequests {
			break
		}

		// Rate limited — wait until the reset window.
		resp.Body.Close()
		resetTime := parseResetTime(resp.Header.Get("x-rate-limit-reset"))
		wait := time.Until(resetTime)
		if wait <= 0 {
			wait = time.Duration(10*(attempt+1)) * time.Second
		}
		if wait > 5*time.Minute {
			return nil, fmt.Errorf("rate limited — reset too far away (%s), giving up", wait.Round(time.Second))
		}
		fmt.Fprintf(os.Stderr, "⏳ Rate limited, waiting %s...\n", wait.Round(time.Second))
		time.Sleep(wait)
	}

	if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
		resp.Body.Close()
		return nil, fmt.Errorf("rate limited after retries — try again later")
	}

	return resp, nil
}

// jitter sleeps for a random 2–5 seconds to reduce API pressure.
func jitter() {
	ms := 2000 + rand.Intn(3000) //nolint:gosec
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// logRateLimit prints rate limit info to stderr when limits are getting low.
func logRateLimit(resp *http.Response) {
	remaining := resp.Header.Get("x-rate-limit-remaining")
	reset := resp.Header.Get("x-rate-limit-reset")
	if remaining == "" {
		return
	}
	rem, _ := strconv.Atoi(remaining)
	if rem <= 10 {
		resetTime := parseResetTime(reset)
		fmt.Fprintf(os.Stderr, "⚠️  Rate limit: %d requests remaining, resets %s\n", rem, resetTime.Format(time.Kitchen))
	}
}

// restGet executes a Twitter REST v1.1 GET request.
// path is relative to https://api.twitter.com (e.g., "/1.1/followers/list.json").
// params are query parameters.
func (c *Client) restGet(path string, params url.Values) (json.RawMessage, error) {
	fullURL := "https://api.twitter.com" + path + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: HTTP %d", resp.StatusCode)
	}

	var raw json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return raw, nil
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
