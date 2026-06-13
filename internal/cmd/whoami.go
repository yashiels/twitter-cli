package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yashiels/twitter-cli/internal/api"
	"github.com/yashiels/twitter-cli/internal/auth"
	"github.com/yashiels/twitter-cli/internal/output"
)

// NewWhoamiCmd returns the "whoami" command.
// It shows the authenticated user's profile using stored credentials.
func NewWhoamiCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show your profile",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			creds, err := auth.Load()
			if errors.Is(err, auth.ErrNotAuthenticated) {
				return fmt.Errorf("%w", auth.ErrNotAuthenticated)
			}
			if err != nil {
				return err
			}

			if creds.Handle == "" && creds.UserID == "" {
				return fmt.Errorf("no identity stored — re-run: twt auth login")
			}

			client := api.NewClient(creds)
			p := output.New(opts)

			// Prefer handle lookup; fall back to ID if handle is stale or missing.
			if creds.Handle != "" {
				user, lookupErr := client.GetUserByScreenName(creds.Handle)
				if lookupErr == nil {
					return p.PrintUser(user)
				}
				// Handle might be stale — fall through to ID lookup if available.
				if creds.UserID == "" {
					return fmt.Errorf("whoami: %w", lookupErr)
				}
				p.Infof("Handle @%s lookup failed, trying by user ID...", creds.Handle)
			}

			user, err := client.GetUserByID(creds.UserID)
			if errors.Is(err, api.ErrUserNotFound) {
				return fmt.Errorf("profile not found — re-run: twt auth login")
			}
			if err != nil {
				return fmt.Errorf("whoami: %w", err)
			}
			return p.PrintUser(user)
		},
	}
}
