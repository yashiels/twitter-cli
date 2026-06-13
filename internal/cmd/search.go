package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yashiels/twitter-cli/internal/api"
	"github.com/yashiels/twitter-cli/internal/auth"
	"github.com/yashiels/twitter-cli/internal/output"
)

// NewSearchCmd returns the "search <query>" command.
func NewSearchCmd(opts *output.Options, limit *int) *cobra.Command {
	var searchUsers bool

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search tweets (or users with --users)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := strings.Join(args, " ")

			creds, err := auth.Load()
			if errors.Is(err, auth.ErrNotAuthenticated) {
				return fmt.Errorf("%w", auth.ErrNotAuthenticated)
			}
			if err != nil {
				return err
			}

			client := api.NewClient(creds)
			p := output.New(opts)

			if searchUsers {
				p.Infof("Searching users for %q (limit %d)...", query, *limit)
			} else {
				p.Infof("Searching tweets for %q (limit %d)...", query, *limit)
			}

			result, err := client.SearchTimeline(query, *limit, searchUsers)
			if err != nil {
				return fmt.Errorf("search: %w", err)
			}

			if searchUsers {
				if len(result.Users) == 0 {
					fmt.Printf("No users found for %q\n", query)
					return nil
				}
				if len(result.Users) > *limit {
					result.Users = result.Users[:*limit]
				}
				return p.PrintUsers(result.Users)
			}

			if len(result.Tweets) == 0 {
				fmt.Printf("No tweets found for %q\n", query)
				return nil
			}
			if len(result.Tweets) > *limit {
				result.Tweets = result.Tweets[:*limit]
			}
			return p.PrintTweets(result.Tweets)
		},
	}

	cmd.Flags().BoolVar(&searchUsers, "users", false, "Search for users instead of tweets")
	return cmd
}
