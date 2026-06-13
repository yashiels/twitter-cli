package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yashiels/twitter-cli/internal/api"
	"github.com/yashiels/twitter-cli/internal/auth"
	"github.com/yashiels/twitter-cli/internal/output"
)

// NewDeleteCmd returns the "delete <tweet-id>" command.
func NewDeleteCmd(opts *output.Options) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <tweet-id>",
		Short: "Delete a tweet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tweetID := args[0]

			creds, err := auth.Load()
			if errors.Is(err, auth.ErrNotAuthenticated) {
				return fmt.Errorf("%w", auth.ErrNotAuthenticated)
			}
			if err != nil {
				return err
			}

			if !yes {
				fmt.Printf("Delete tweet %s? [y/N]: ", tweetID)
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
				if answer != "y" && answer != "yes" {
					fmt.Println("Cancelled.")
					return nil
				}
			}

			client := api.NewClient(creds)
			p := output.New(opts)

			p.Infof("Deleting tweet %s...", tweetID)
			if err := client.DeletePost(tweetID); err != nil {
				return fmt.Errorf("delete: %w", err)
			}

			fmt.Printf("✓ Deleted tweet %s\n", tweetID)
			return nil
		},
	}

	cmd.Flags().BoolVar(&yes, "yes", false, "Skip confirmation prompt")
	return cmd
}
