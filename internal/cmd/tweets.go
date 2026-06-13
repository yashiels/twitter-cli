package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yashiels/twitter-cli/internal/api"
	"github.com/yashiels/twitter-cli/internal/auth"
	"github.com/yashiels/twitter-cli/internal/output"
)

// NewTweetsCmd returns the "tweets <handle>" command.
func NewTweetsCmd(opts *output.Options, limit *int) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tweets <handle>",
		Short: "Show recent tweets from a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handle := strings.TrimPrefix(args[0], "@")

			creds, err := auth.Load()
			if errors.Is(err, auth.ErrNotAuthenticated) {
				return fmt.Errorf("%w", auth.ErrNotAuthenticated)
			}
			if err != nil {
				return err
			}

			client := api.NewClient(creds)
			p := output.New(opts)

			// Step 1: resolve handle → user ID.
			p.Infof("Resolving @%s...", handle)
			user, err := client.GetUserByScreenName(handle)
			if errors.Is(err, api.ErrUserNotFound) {
				output.Errorf("user @%s not found", handle)
				os.Exit(5)
			}
			if err != nil {
				return fmt.Errorf("lookup user: %w", err)
			}

			if user.RestID == "" {
				output.Errorf("could not resolve user ID for @%s", handle)
				os.Exit(5)
			}

			// Step 2: fetch tweets.
			p.Infof("Fetching tweets for @%s (limit %d)...", handle, *limit)
			tweets, err := client.GetUserTweets(user.RestID, *limit)
			if err != nil {
				return fmt.Errorf("fetch tweets: %w", err)
			}

			// Set author handle on tweets that are missing it (fallback).
			for _, t := range tweets {
				if t.AuthorHandle == "" {
					t.AuthorHandle = user.ScreenName
					t.URL = t.TweetURL()
				}
			}

			if len(tweets) == 0 {
				fmt.Printf("No tweets found for @%s\n", handle)
				return nil
			}

			// Trim to limit.
			if len(tweets) > *limit {
				tweets = tweets[:*limit]
			}

			return p.PrintTweets(tweets)
		},
	}

	return cmd
}
