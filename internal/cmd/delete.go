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

			// Confirmation prompt unless --yes.
			if !yes {
				fmt.Printf("Delete tweet %s? [y/N]: ", tweetID)
				reader := bufio.NewReader(os.Stdin)
				line, _ := reader.ReadString('\n')
				line = strings.TrimSpace(strings.ToLower(line))
				if line != "y" && line != "yes" {
					fmt.Println("Aborted.")
					return nil
				}
			}

			client := api.NewClient(creds)
			p := output.New(opts)

			p.Infof("Deleting tweet %s...", tweetID)
			if err := client.DeletePost(tweetID); err != nil {
				return fmt.Errorf("delete tweet: %w", err)
			}

			fmt.Printf("✓ Deleted tweet %s\n", tweetID)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")
	return cmd
}
