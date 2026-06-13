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
	"github.com/yashiels/twitter-cli/internal/types"
)

// NewPostCmd returns the "post <text>" command.
func NewPostCmd(opts *output.Options) *cobra.Command {
	var yes bool
	var reply, quote string

	cmd := &cobra.Command{
		Use:   "post <text>",
		Short: "Post a new tweet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			text := args[0]

			creds, err := auth.Load()
			if errors.Is(err, auth.ErrNotAuthenticated) {
				return fmt.Errorf("%w", auth.ErrNotAuthenticated)
			}
			if err != nil {
				return err
			}

			if !yes {
				typeStr := "tweet"
				if reply != "" {
					typeStr = "reply"
				} else if quote != "" {
					typeStr = "quote tweet"
				}
				fmt.Printf("Post this %s?\n\n  %s\n\n[y/N]: ", typeStr, text)
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

			// Reject conflicting flags.
			if reply != "" && quote != "" {
				return fmt.Errorf("cannot use --reply and --quote together")
			}

			var tweet *types.Tweet
			switch {
			case reply != "":
				p.Infof("Posting reply to %s...", reply)
				tweet, err = client.CreateReply(text, reply)
				if err != nil {
					return fmt.Errorf("post reply: %w", err)
				}
			case quote != "":
				p.Infof("Posting quote tweet of %s...", quote)
				tweet, err = client.CreateQuote(text, quote)
				if err != nil {
					return fmt.Errorf("post quote: %w", err)
				}
			default:
				p.Infof("Posting tweet...")
				tweet, err = client.CreatePost(text)
				if err != nil {
					return fmt.Errorf("post: %w", err)
				}
			}

			if tweet != nil && tweet.URL != "" {
				fmt.Printf("✓ Posted: %s\n", tweet.URL)
			} else if tweet != nil {
				fmt.Printf("✓ Posted tweet %s\n", tweet.ID)
			} else {
				fmt.Println("✓ Posted.")
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&yes, "yes", false, "Skip confirmation prompt")
	cmd.Flags().StringVar(&reply, "reply", "", "Reply to tweet ID")
	cmd.Flags().StringVar(&quote, "quote", "", "Quote tweet ID")
	return cmd
}
