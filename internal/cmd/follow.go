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

// NewFollowCmd returns the "follow <handle>" command.
func NewFollowCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "follow <handle>",
		Short: "Follow a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handle := strings.TrimPrefix(args[0], "@")

			creds, err := auth.Load()
			if errors.Is(err, auth.ErrNotAuthenticated) {
				return fmt.Errorf("%w", auth.ErrNotAuthenticated)
			}
			if err != nil {
				return err
			}

			client := api.NewClient(creds)
			p := output.New(opts)

			p.Infof("Resolving @%s...", handle)
			user, err := client.GetUserByScreenName(handle)
			if errors.Is(err, api.ErrUserNotFound) {
				output.Errorf("user @%s not found", handle)
				os.Exit(5)
			}
			if err != nil {
				return fmt.Errorf("lookup user: %w", err)
			}

			p.Infof("Following @%s...", handle)
			if err := client.FollowUser(user.RestID); err != nil {
				return fmt.Errorf("follow: %w", err)
			}

			fmt.Printf("✓ Now following @%s\n", handle)
			return nil
		},
	}
}

// NewUnfollowCmd returns the "unfollow <handle>" command.
func NewUnfollowCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "unfollow <handle>",
		Short: "Unfollow a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handle := strings.TrimPrefix(args[0], "@")

			creds, err := auth.Load()
			if errors.Is(err, auth.ErrNotAuthenticated) {
				return fmt.Errorf("%w", auth.ErrNotAuthenticated)
			}
			if err != nil {
				return err
			}

			client := api.NewClient(creds)
			p := output.New(opts)

			p.Infof("Resolving @%s...", handle)
			user, err := client.GetUserByScreenName(handle)
			if errors.Is(err, api.ErrUserNotFound) {
				output.Errorf("user @%s not found", handle)
				os.Exit(5)
			}
			if err != nil {
				return fmt.Errorf("lookup user: %w", err)
			}

			p.Infof("Unfollowing @%s...", handle)
			if err := client.UnfollowUser(user.RestID); err != nil {
				return fmt.Errorf("unfollow: %w", err)
			}

			fmt.Printf("✓ Unfollowed @%s\n", handle)
			return nil
		},
	}
}
