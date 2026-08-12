package api

import (
	"encoding/json"
	"testing"
)

func TestGetNestedJSON_WalksKeys(t *testing.T) {
	raw := json.RawMessage(`{"data":{"user":{"result":{"id":"123"}}}}`)
	got, err := getNestedJSON(raw, "data", "user", "result")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]string
	if err := json.Unmarshal(got, &obj); err != nil {
		t.Fatalf("result not an object: %v", err)
	}
	if obj["id"] != "123" {
		t.Errorf("got id=%q, want 123", obj["id"])
	}
}

func TestGetNestedJSON_MissingKey(t *testing.T) {
	raw := json.RawMessage(`{"data":{}}`)
	if _, err := getNestedJSON(raw, "data", "user"); err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestGetNestedJSON_InvalidJSON(t *testing.T) {
	raw := json.RawMessage(`{not valid`)
	if _, err := getNestedJSON(raw, "data"); err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestGetNestedJSON_NoKeysReturnsInput(t *testing.T) {
	raw := json.RawMessage(`{"a":1}`)
	got, err := getNestedJSON(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != string(raw) {
		t.Errorf("got %q, want input unchanged", string(got))
	}
}
