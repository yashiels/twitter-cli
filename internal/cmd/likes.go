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
// Without an argument, shows the authenticated user's likes.
// With a handle argument, shows that user's likes.
func NewLikesCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "likes [handle]",
		Short: "List liked tweets",
		Args:  cobra.MaximumNArgs(1),
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

			var userID string

			if len(args) == 0 {
				// Use the authenticated user's stored ID.
				if creds.UserID == "" {
					return fmt.Errorf("user ID not stored — run: twt auth login to refresh credentials")
				}
				userID = creds.UserID
				p.Infof("Fetching your likes...")
			} else {
				// Resolve the provided handle to a user ID.
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
				p.Infof("Fetching likes for @%s...", handle)
			}

			tweets, err := client.GetUserLikes(userID, limit)
			if err != nil {
				return fmt.Errorf("likes: %w", err)
			}

			return p.PrintTweets(tweets)
		},
	}
}
