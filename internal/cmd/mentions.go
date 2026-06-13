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
func NewMentionsCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "mentions",
		Short: "Show your recent mentions and notifications",
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

			p.Infof("Fetching mentions...")
			tweets, err := client.GetMentions(limit)
			if err != nil {
				return fmt.Errorf("mentions: %w", err)
			}

			return p.PrintTweets(tweets)
		},
	}
}
