package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yashiels/twitter-cli/internal/api"
	"github.com/yashiels/twitter-cli/internal/auth"
	"github.com/yashiels/twitter-cli/internal/output"
)

// NewRepostCmd returns the "repost <tweet-id>" command.
func NewRepostCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "repost <tweet-id>",
		Short: "Repost (retweet) a tweet",
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

			p.Infof("Reposting tweet %s...", tweetID)
			if err := client.CreateRepost(tweetID); err != nil {
				return fmt.Errorf("repost: %w", err)
			}

			fmt.Printf("✓ Reposted tweet %s\n", tweetID)
			return nil
		},
	}
}

// NewUnrepostCmd returns the "unrepost <tweet-id>" command.
func NewUnrepostCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "unrepost <tweet-id>",
		Short: "Remove a repost (undo retweet)",
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

			p.Infof("Removing repost of tweet %s...", tweetID)
			if err := client.DeleteRepost(tweetID); err != nil {
				return fmt.Errorf("unrepost: %w", err)
			}

			fmt.Printf("✓ Unreposted tweet %s\n", tweetID)
			return nil
		},
	}
}
