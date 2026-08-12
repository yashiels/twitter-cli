package output

import (
	"strings"
	"testing"
	"time"
)

func TestFormatCount(t *testing.T) {
	cases := []struct {
		in   int
		want string
	}{
		{0, "0"},
		{42, "42"},
		{999, "999"},
		{1000, "1.0K"},
		{1500, "1.5K"},
		{999_999, "1000.0K"},
		{1_000_000, "1.0M"},
		{2_500_000, "2.5M"},
	}
	for _, tc := range cases {
		if got := formatCount(tc.in); got != tc.want {
			t.Errorf("formatCount(%d) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestRelativeTime(t *testing.T) {
	if got := relativeTime(time.Time{}); got != "unknown time" {
		t.Errorf("zero time = %q, want %q", got, "unknown time")
	}
	now := time.Now()
	cases := []struct {
		name string
		t    time.Time
		want string
	}{
		{"seconds", now.Add(-10 * time.Second), "just now"},
		{"minutes", now.Add(-30 * time.Minute), "30m ago"},
		{"hours", now.Add(-5 * time.Hour), "5h ago"},
		{"days", now.Add(-3 * 24 * time.Hour), "3d ago"},
	}
	for _, tc := range cases {
		if got := relativeTime(tc.t); got != tc.want {
			t.Errorf("%s: relativeTime = %q, want %q", tc.name, got, tc.want)
		}
	}
	// Older than a week falls back to an absolute date.
	old := now.Add(-30 * 24 * time.Hour)
	if got := relativeTime(old); got != old.Format("Jan 2, 2006") {
		t.Errorf("old time = %q, want absolute date %q", got, old.Format("Jan 2, 2006"))
	}
}

func TestWordWrap_ShortTextUnchanged(t *testing.T) {
	in := "hello world"
	if got := wordWrap(in, 40); got != in {
		t.Errorf("short text = %q, want unchanged %q", got, in)
	}
}

func TestWordWrap_WrapsAtBoundary(t *testing.T) {
	got := wordWrap("the quick brown fox jumps", 9)
	for _, line := range strings.Split(got, "\n") {
		if len([]rune(line)) > 9 {
			t.Errorf("line %q exceeds width 9", line)
		}
	}
	// Wrapping must not drop or reorder words.
	if strings.Join(strings.Fields(got), " ") != "the quick brown fox jumps" {
		t.Errorf("wrapped text lost content: %q", got)
	}
}

func TestWordWrap_PreservesParagraphs(t *testing.T) {
	got := wordWrap("line one here\nline two here", 40)
	if !strings.Contains(got, "\n") {
		t.Errorf("expected paragraph break preserved, got %q", got)
	}
}
