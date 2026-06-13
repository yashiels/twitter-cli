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

// NewPostCmd returns the "post <text>" command with optional --reply and --quote flags.
func NewPostCmd(opts *output.Options) *cobra.Command {
	var yes bool
	var replyTo string
	var quoteTweet string

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

			client := api.NewClient(creds)
			p := output.New(opts)

			// Determine mode.
			var mode string
			switch {
			case replyTo != "":
				mode = fmt.Sprintf("reply to tweet %s", replyTo)
			case quoteTweet != "":
				mode = fmt.Sprintf("quote tweet %s", quoteTweet)
			default:
				mode = "new tweet"
			}

			// Confirmation prompt unless --yes.
			if !yes {
				fmt.Printf("Post %s?\n  %q\n[y/N]: ", mode, text)
				reader := bufio.NewReader(os.Stdin)
				line, _ := reader.ReadString('\n')
				line = strings.TrimSpace(strings.ToLower(line))
				if line != "y" && line != "yes" {
					fmt.Println("Aborted.")
					return nil
				}
			}

			var tweetID string

			switch {
			case replyTo != "":
				p.Infof("Posting reply...")
				t, err := client.CreateReply(text, replyTo)
				if err != nil {
					return fmt.Errorf("post reply: %w", err)
				}
				tweetID = t.ID
				handle := t.AuthorHandle
				if handle == "" {
					handle = creds.Handle
				}
				fmt.Printf("✓ Replied: https://x.com/%s/status/%s\n", handle, tweetID)

			case quoteTweet != "":
				p.Infof("Posting quote tweet...")
				t, err := client.CreateQuote(text, quoteTweet)
				if err != nil {
					return fmt.Errorf("post quote: %w", err)
				}
				tweetID = t.ID
				handle := t.AuthorHandle
				if handle == "" {
					handle = creds.Handle
				}
				fmt.Printf("✓ Quoted: https://x.com/%s/status/%s\n", handle, tweetID)

			default:
				p.Infof("Posting tweet...")
				t, err := client.CreatePost(text)
				if err != nil {
					return fmt.Errorf("post tweet: %w", err)
				}
				tweetID = t.ID
				handle := t.AuthorHandle
				if handle == "" {
					handle = creds.Handle
				}
				fmt.Printf("✓ Posted: https://x.com/%s/status/%s\n", handle, tweetID)
			}

			_ = tweetID
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")
	cmd.Flags().StringVar(&replyTo, "reply", "", "Reply to tweet ID")
	cmd.Flags().StringVar(&quoteTweet, "quote", "", "Quote tweet ID")
	return cmd
}
