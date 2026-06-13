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

// NewLikesCmd returns the "likes [handle]" command.
// With no argument it uses the stored user ID from credentials.
func NewLikesCmd(opts *output.Options, limit *int) *cobra.Command {
	return &cobra.Command{
		Use:   "likes [handle]",
		Short: "List liked tweets",
		Args:  cobra.MaximumNArgs(1),
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

			var userID, displayHandle string

			if len(args) == 0 {
				// Use stored user ID from credentials.
				if creds.UserID == "" {
					return fmt.Errorf("user ID not stored — re-run: twt auth login")
				}
				userID = creds.UserID
				displayHandle = creds.Handle
				if displayHandle == "" {
					displayHandle = userID
				}
			} else {
				handle := strings.TrimPrefix(args[0], "@")
				p.Infof("Resolving @%s...", handle)
				user, err := client.GetUserByScreenName(handle)
				if errors.Is(err, api.ErrUserNotFound) {
					output.Errorf("user @%s not found", handle)
					os.Exit(5)
				}
				if err != nil {
					return fmt.Errorf("lookup user: %w", err)
				}
				userID = user.RestID
				displayHandle = user.ScreenName
			}

			p.Infof("Fetching likes for @%s (limit %d)...", displayHandle, *limit)
			tweets, err := client.GetUserLikes(userID, *limit)
			if err != nil {
				return fmt.Errorf("likes: %w", err)
			}

			if len(tweets) == 0 {
				fmt.Println("No likes found.")
				return nil
			}

			if len(tweets) > *limit {
				tweets = tweets[:*limit]
			}

			return p.PrintTweets(tweets)
		},
	}
}
