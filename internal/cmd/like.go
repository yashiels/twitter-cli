package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yashiels/twitter-cli/internal/api"
	"github.com/yashiels/twitter-cli/internal/auth"
	"github.com/yashiels/twitter-cli/internal/output"
)

// NewLikeCmd returns the "like <tweet-id>" command.
func NewLikeCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "like <tweet-id>",
		Short: "Like a tweet",
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

			p.Infof("Liking tweet %s...", tweetID)
			if err := client.LikeTweet(tweetID); err != nil {
				return fmt.Errorf("like: %w", err)
			}

			fmt.Printf("✓ Liked tweet %s\n", tweetID)
			return nil
		},
	}
}

// NewUnlikeCmd returns the "unlike <tweet-id>" command.
func NewUnlikeCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "unlike <tweet-id>",
		Short: "Unlike a tweet",
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

			p.Infof("Unliking tweet %s...", tweetID)
			if err := client.UnlikeTweet(tweetID); err != nil {
				return fmt.Errorf("unlike: %w", err)
			}

			fmt.Printf("✓ Unliked tweet %s\n", tweetID)
			return nil
		},
	}
}
