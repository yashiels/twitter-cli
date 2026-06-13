package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yashiels/twitter-cli/internal/api"
	"github.com/yashiels/twitter-cli/internal/auth"
	"github.com/yashiels/twitter-cli/internal/output"
	"github.com/yashiels/twitter-cli/internal/types"
)

// NewWhoamiCmd returns the "whoami" command.
func NewWhoamiCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show your own profile",
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

			var user *types.User

			if creds.Handle != "" {
				p.Infof("Fetching profile for @%s...", creds.Handle)
				u, err := client.GetUserByScreenName(creds.Handle)
				if err == nil {
					user = u
				}
			}

			if user == nil && creds.UserID != "" {
				p.Infof("Fetching profile by user ID...")
				u, err := client.GetUserByID(creds.UserID)
				if err == nil {
					user = u
				}
			}

			if user == nil {
				return fmt.Errorf("could not resolve your profile — try: twt auth login to refresh credentials")
			}

			return p.PrintUser(user)
		},
	}
}
