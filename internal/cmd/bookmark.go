package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yashiels/twitter-cli/internal/api"
	"github.com/yashiels/twitter-cli/internal/auth"
	"github.com/yashiels/twitter-cli/internal/output"
)

// NewBookmarksCmd returns the "bookmarks" command.
func NewBookmarksCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "bookmarks",
		Short: "List your bookmarked tweets",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			limit, _ := cmd.Root().PersistentFlags().GetInt("limit")

			creds, err := auth.Load()
			if errors.Is(err, auth.ErrNotAuthenticated) {
				return fmt.Errorf("%w", auth.ErrNotAuthenticated)
			}
			if err != nil {
				return err
			}

			client := api.NewClient(creds)
			p := output.New(opts)

			p.Infof("Fetching bookmarks...")
			tweets, err := client.GetBookmarks(limit)
			if err != nil {
				return fmt.Errorf("bookmarks: %w", err)
			}

			return p.PrintTweets(tweets)
		},
	}
}

// NewBookmarkCmd returns the "bookmark <tweet-id>" command.
func NewBookmarkCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "bookmark <tweet-id>",
		Short: "Bookmark a tweet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tweetID := args[0]

			creds, err := auth.Load()
			if errors.Is(err, auth.ErrNotAuthenticated) {
				return fmt.Errorf("%w", auth.ErrNotAuthenticated)
			}
			if err != nil {
				return err
			}

			client := api.NewClient(creds)
			p := output.New(opts)

			p.Infof("Bookmarking tweet %s...", tweetID)
			if err := client.AddBookmark(tweetID); err != nil {
				return fmt.Errorf("bookmark: %w", err)
			}

			fmt.Printf("✓ Bookmarked tweet %s\n", tweetID)
			return nil
		},
	}
}

// NewUnbookmarkCmd returns the "unbookmark <tweet-id>" command.
func NewUnbookmarkCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "unbookmark <tweet-id>",
		Short: "Remove a bookmark",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tweetID := args[0]

			creds, err := auth.Load()
			if errors.Is(err, auth.ErrNotAuthenticated) {
				return fmt.Errorf("%w", auth.ErrNotAuthenticated)
			}
			if err != nil {
				return err
			}

			client := api.NewClient(creds)
			p := output.New(opts)

			p.Infof("Removing bookmark for tweet %s...", tweetID)
			if err := client.RemoveBookmark(tweetID); err != nil {
				return fmt.Errorf("unbookmark: %w", err)
			}

			fmt.Printf("✓ Removed bookmark %s\n", tweetID)
			return nil
		},
	}
}
