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

// NewBlockCmd returns the "block <handle>" command.
func NewBlockCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "block <handle>",
		Short: "Block a user",
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

			p.Infof("Blocking @%s...", handle)
			if err := client.BlockUser(user.RestID); err != nil {
				return fmt.Errorf("block: %w", err)
			}

			fmt.Printf("✓ Blocked @%s\n", handle)
			return nil
		},
	}
}

// NewUnblockCmd returns the "unblock <handle>" command.
func NewUnblockCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "unblock <handle>",
		Short: "Unblock a user",
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

			p.Infof("Unblocking @%s...", handle)
			if err := client.UnblockUser(user.RestID); err != nil {
				return fmt.Errorf("unblock: %w", err)
			}

			fmt.Printf("✓ Unblocked @%s\n", handle)
			return nil
		},
	}
}

// NewMuteCmd returns the "mute <handle>" command.
func NewMuteCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "mute <handle>",
		Short: "Mute a user",
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

			p.Infof("Muting @%s...", handle)
			if err := client.MuteUser(user.RestID); err != nil {
				return fmt.Errorf("mute: %w", err)
			}

			fmt.Printf("✓ Muted @%s\n", handle)
			return nil
		},
	}
}

// NewUnmuteCmd returns the "unmute <handle>" command.
func NewUnmuteCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "unmute <handle>",
		Short: "Unmute a user",
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

			p.Infof("Unmuting @%s...", handle)
			if err := client.UnmuteUser(user.RestID); err != nil {
				return fmt.Errorf("unmute: %w", err)
			}

			fmt.Printf("✓ Unmuted @%s\n", handle)
			return nil
		},
	}
}
