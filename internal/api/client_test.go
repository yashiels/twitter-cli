package api

import (
	"testing"
	"time"
)

func TestParseResetTime_ValidUnix(t *testing.T) {
	got := parseResetTime("1700000000")
	want := time.Unix(1700000000, 0)
	if !got.Equal(want) {
		t.Errorf("parseResetTime = %v, want %v", got, want)
	}
}

func TestParseResetTime_EmptyFallsBackToFuture(t *testing.T) {
	before := time.Now()
	got := parseResetTime("")
	// Empty input backs off ~60s into the future.
	if !got.After(before.Add(50 * time.Second)) {
		t.Errorf("empty input = %v, expected ~60s in the future from %v", got, before)
	}
}

func TestParseResetTime_InvalidFallsBackToFuture(t *testing.T) {
	before := time.Now()
	got := parseResetTime("not-a-number")
	if !got.After(before.Add(50 * time.Second)) {
		t.Errorf("invalid input = %v, expected ~60s in the future from %v", got, before)
	}
}
