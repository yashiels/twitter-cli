package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yashiels/twitter-cli/internal/api"
	"github.com/yashiels/twitter-cli/internal/auth"
	"github.com/yashiels/twitter-cli/internal/output"
	"github.com/yashiels/twitter-cli/internal/types"
)

// NewTweetCmd returns the "tweet <id>" command.
func NewTweetCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "tweet <id>",
		Short: "Show a single tweet by ID",
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

			p.Infof("Fetching tweet %s...", tweetID)
			t, err := client.GetTweetByID(tweetID)
			if err != nil {
				output.Errorf("%v", err)
				os.Exit(5)
			}

			return p.PrintTweets([]*types.Tweet{t})
		},
	}
}
