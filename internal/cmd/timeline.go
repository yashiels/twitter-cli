package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yashiels/twitter-cli/internal/api"
	"github.com/yashiels/twitter-cli/internal/auth"
	"github.com/yashiels/twitter-cli/internal/output"
)

// NewTimelineCmd returns the "timeline" command.
func NewTimelineCmd(opts *output.Options, limit *int) *cobra.Command {
	var latest bool

	cmd := &cobra.Command{
		Use:   "timeline",
		Short: "Show your home timeline",
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

			tab := "For You"
			if latest {
				tab = "Following"
			}
			p.Infof("Fetching %s timeline (limit %d)...", tab, *limit)

			tweets, err := client.GetHomeTimeline(*limit, latest)
			if err != nil {
				return fmt.Errorf("timeline: %w", err)
			}

			if len(tweets) == 0 {
				fmt.Println("No tweets in timeline.")
				return nil
			}

			if len(tweets) > *limit {
				tweets = tweets[:*limit]
			}

			return p.PrintTweets(tweets)
		},
	}

	cmd.Flags().BoolVar(&latest, "latest", false, "Show Following (latest) tab instead of For You")
	return cmd
}
