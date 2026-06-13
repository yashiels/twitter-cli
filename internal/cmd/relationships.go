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

// NewFollowersCmd returns the "followers <handle>" command.
func NewFollowersCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "followers <handle>",
		Short: "List a user's followers",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handle := strings.TrimPrefix(args[0], "@")
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

			p.Infof("Fetching followers of @%s...", handle)
			users, err := client.GetFollowers(handle, limit)
			if err != nil {
				return fmt.Errorf("followers: %w", err)
			}

			return p.PrintUsers(users)
		},
	}
}

// NewFollowingCmd returns the "following <handle>" command.
func NewFollowingCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "following <handle>",
		Short: "List users that a user follows",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handle := strings.TrimPrefix(args[0], "@")
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

			p.Infof("Fetching accounts @%s follows...", handle)
			users, err := client.GetFollowing(handle, limit)
			if err != nil {
				return fmt.Errorf("following: %w", err)
			}

			return p.PrintUsers(users)
		},
	}
}
