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

// NewUserCmd returns the "user <handle>" command.
func NewUserCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "user <handle>",
		Short: "Look up a Twitter/X user profile",
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
			user, err := client.GetUserByScreenName(handle)
			if errors.Is(err, api.ErrUserNotFound) {
				output.Errorf("user @%s not found", handle)
				os.Exit(5)
			}
			if err != nil {
				return err
			}

			p := output.New(opts)
			return p.PrintUser(user)
		},
	}
}
