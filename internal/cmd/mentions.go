package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yashiels/twitter-cli/internal/api"
	"github.com/yashiels/twitter-cli/internal/auth"
	"github.com/yashiels/twitter-cli/internal/output"
)

// NewMentionsCmd returns the "mentions" command.
func NewMentionsCmd(opts *output.Options, limit *int) *cobra.Command {
	return &cobra.Command{
		Use:   "mentions",
		Short: "Show your mentions / notification timeline",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			creds, err := auth.Load()
			if errors.Is(err, auth.ErrNotAuthenticated) {
				return fmt.Errorf("%w", auth.ErrNotAuthenticated)
			}
			if err != nil {
				return err
			}

			client := api.NewClient(creds)
			p := output.New(opts)

			p.Infof("Fetching mentions (limit %d)...", *limit)
			tweets, err := client.GetMentions(*limit)
			if err != nil {
				return fmt.Errorf("mentions: %w", err)
			}

			if len(tweets) == 0 {
				fmt.Println("No mentions found.")
				return nil
			}

			if len(tweets) > *limit {
				tweets = tweets[:*limit]
			}

			return p.PrintTweets(tweets)
		},
	}
}
