package api

import (
	"encoding/json"
	"testing"
)

// TestParseTweetResult_ViewsCountNumberOrString ensures a single field-type
// variance on views.count (observed as both a quoted string and a bare number
// across the APK/web schemas) does not discard the whole tweet.
func TestParseTweetResult_ViewsCountNumberOrString(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want int
	}{
		{"string", `{"rest_id":"1","details":{"full_text":"hi"},"views":{"count":"1234"}}`, 1234},
		{"number", `{"rest_id":"2","details":{"full_text":"hi"},"views":{"count":5678}}`, 5678},
		{"empty", `{"rest_id":"3","details":{"full_text":"hi"},"views":{"count":""}}`, 0},
		{"missing", `{"rest_id":"4","details":{"full_text":"hi"}}`, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tw, err := parseTweetResult(json.RawMessage(tc.raw))
			if err != nil {
				t.Fatalf("parseTweetResult error: %v", err)
			}
			if tw == nil {
				t.Fatalf("tweet was dropped, want parsed")
			}
			if tw.ViewCount != tc.want {
				t.Errorf("ViewCount = %d, want %d", tw.ViewCount, tc.want)
			}
		})
	}
}

// TestParseTweetResult_CreatedAtMsNumberOrString ensures created_at_ms parses
// whether delivered as a number or a numeric string.
func TestParseTweetResult_CreatedAtMsNumberOrString(t *testing.T) {
	cases := []struct {
		name string
		raw  string
	}{
		{"number", `{"rest_id":"1","details":{"full_text":"hi","created_at_ms":1700000000000}}`},
		{"string", `{"rest_id":"2","details":{"full_text":"hi","created_at_ms":"1700000000000"}}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tw, err := parseTweetResult(json.RawMessage(tc.raw))
			if err != nil {
				t.Fatalf("parseTweetResult error: %v", err)
			}
			if tw == nil {
				t.Fatalf("tweet was dropped, want parsed")
			}
			if tw.CreatedAt.UnixMilli() != 1700000000000 {
				t.Errorf("CreatedAt = %v (%d ms), want 1700000000000 ms", tw.CreatedAt, tw.CreatedAt.UnixMilli())
			}
		})
	}
}

// TestParseTweetResult_QuoteAndBookmarkCounts ensures quote/bookmark counts are
// surfaced onto the domain Tweet from both the APK (counts.*) and web (legacy.*)
// schemas.
func TestParseTweetResult_QuoteAndBookmarkCounts(t *testing.T) {
	cases := []struct {
		name              string
		raw               string
		wantQuote, wantBk int
	}{
		{
			"apk counts",
			`{"rest_id":"1","details":{"full_text":"hi"},"counts":{"quote_count":7,"bookmark_count":9}}`,
			7, 9,
		},
		{
			"legacy fallback",
			`{"rest_id":"2","legacy":{"full_text":"hi","quote_count":3,"bookmark_count":4}}`,
			3, 4,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tw, err := parseTweetResult(json.RawMessage(tc.raw))
			if err != nil {
				t.Fatalf("parseTweetResult error: %v", err)
			}
			if tw == nil {
				t.Fatalf("tweet was dropped, want parsed")
			}
			if tw.QuoteCount != tc.wantQuote {
				t.Errorf("QuoteCount = %d, want %d", tw.QuoteCount, tc.wantQuote)
			}
			if tw.BookmarkCount != tc.wantBk {
				t.Errorf("BookmarkCount = %d, want %d", tw.BookmarkCount, tc.wantBk)
			}
		})
	}
}

// TestParseTimelineInstructions_ModuleEntry ensures TimelineTimelineModule
// entries (threads/clusters carried under content.items[]) are not silently
// dropped.
func TestParseTimelineInstructions_ModuleEntry(t *testing.T) {
	instructions := `[
	  {"type":"TimelineAddEntries","entries":[
	    {"entryId":"cluster-1","content":{
	      "__typename":"TimelineTimelineModule",
	      "items":[
	        {"item":{"itemContent":{"tweet_results":{"result":{"rest_id":"111","details":{"full_text":"first"}}}}}},
	        {"item":{"itemContent":{"tweet_results":{"result":{"rest_id":"222","details":{"full_text":"second"}}}}}}
	      ]
	    }}
	  ]}
	]`

	tweets, err := parseTimelineInstructions(json.RawMessage(instructions))
	if err != nil {
		t.Fatalf("parseTimelineInstructions error: %v", err)
	}
	if len(tweets) != 2 {
		t.Fatalf("got %d tweets, want 2 (module items dropped)", len(tweets))
	}
	if tweets[0].ID != "111" || tweets[1].ID != "222" {
		t.Errorf("module tweet IDs = %q, %q; want 111, 222", tweets[0].ID, tweets[1].ID)
	}
}

// TestParseTimelineInstructions_SingleEntry is a regression guard that ordinary
// single-tweet entries still parse alongside the new module handling.
func TestParseTimelineInstructions_SingleEntry(t *testing.T) {
	instructions := `[
	  {"type":"TimelineAddEntries","entries":[
	    {"entryId":"tweet-1","content":{"content":{"tweet_results":{"result":{"rest_id":"999","details":{"full_text":"solo"}}}}}}
	  ]}
	]`

	tweets, err := parseTimelineInstructions(json.RawMessage(instructions))
	if err != nil {
		t.Fatalf("parseTimelineInstructions error: %v", err)
	}
	if len(tweets) != 1 {
		t.Fatalf("got %d tweets, want 1", len(tweets))
	}
	if tweets[0].ID != "999" {
		t.Errorf("tweet ID = %q, want 999", tweets[0].ID)
	}
}
