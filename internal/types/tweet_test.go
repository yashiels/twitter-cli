package types

import "testing"

func TestTweetURL(t *testing.T) {
	cases := []struct {
		name  string
		tweet Tweet
		want  string
	}{
		{"handle and id", Tweet{ID: "123", AuthorHandle: "jack"}, "https://x.com/jack/status/123"},
		{"missing handle", Tweet{ID: "123"}, ""},
		{"missing id", Tweet{AuthorHandle: "jack"}, ""},
		{"empty", Tweet{}, ""},
	}
	for _, tc := range cases {
		if got := tc.tweet.TweetURL(); got != tc.want {
			t.Errorf("%s: TweetURL() = %q, want %q", tc.name, got, tc.want)
		}
	}
}
