package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// graphqlGet executes a Twitter GraphQL GET request.
// queryID and operationName form the path; variables and features are JSON-encoded query params.
func (c *Client) graphqlGet(queryID, operationName string, variables, features map[string]any) (json.RawMessage, error) {
	varsJSON, err := json.Marshal(variables)
	if err != nil {
		return nil, fmt.Errorf("marshal variables: %w", err)
	}

	featJSON, err := json.Marshal(features)
	if err != nil {
		return nil, fmt.Errorf("marshal features: %w", err)
	}

	endpoint := fmt.Sprintf("%s/%s/%s", BaseURL, queryID, operationName)

	params := url.Values{}
	params.Set("variables", string(varsJSON))
	params.Set("features", string(featJSON))

	fullURL := endpoint + "?" + params.Encode()

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

	// Check for API-level errors in the response.
	var errCheck struct {
		Errors []struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(raw, &errCheck); err == nil && len(errCheck.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s (code %d)", errCheck.Errors[0].Message, errCheck.Errors[0].Code)
	}

	return raw, nil
}

// getNestedJSON extracts a value from a JSON blob by following a dot-separated path.
func getNestedJSON(data json.RawMessage, keys ...string) (json.RawMessage, error) {
	current := data
	for _, key := range keys {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(current, &obj); err != nil {
			return nil, fmt.Errorf("unmarshal at key %q: %w", key, err)
		}
		val, ok := obj[key]
		if !ok {
			return nil, fmt.Errorf("key %q not found in response", key)
		}
		current = val
	}
	return current, nil
}
