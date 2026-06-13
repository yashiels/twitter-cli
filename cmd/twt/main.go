package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yashiels/twitter-cli/internal/cmd"
	"github.com/yashiels/twitter-cli/internal/output"
)

// version is set at build time via -ldflags "-X main.version=<tag>".
var version = "0.2.0"

func main() {
	opts := output.DefaultOptions()
	var limit int

	root := &cobra.Command{
		Use:   "twt",
		Short: "twt — X/Twitter CLI",
		Long: `twt is a command-line client for X/Twitter using the private GraphQL API.

Authentication is via Chrome cookies on macOS (auto-extracted) or manually
provided auth_token + ct0 cookie values.`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			// Wire global flags into opts before each command runs.
			jsonFlag, _ := c.Flags().GetBool("json")
			plainFlag, _ := c.Flags().GetBool("plain")
			noColor, _ := c.Flags().GetBool("no-color")
			quiet, _ := c.Flags().GetBool("quiet")

			// Also check persistent flags on the root.
			if !jsonFlag {
				jsonFlag, _ = c.Root().PersistentFlags().GetBool("json")
			}
			if !plainFlag {
				plainFlag, _ = c.Root().PersistentFlags().GetBool("plain")
			}
			if !noColor {
				noColor, _ = c.Root().PersistentFlags().GetBool("no-color")
			}
			if !quiet {
				quiet, _ = c.Root().PersistentFlags().GetBool("quiet")
			}

			if jsonFlag {
				opts.Format = output.FormatJSON
			} else if plainFlag {
				opts.Format = output.FormatPlain
			}
			opts.NoColor = noColor || opts.NoColor
			opts.Quiet = quiet
		},
	}

	// Global persistent flags.
	root.PersistentFlags().Bool("json", false, "Output machine-readable JSON")
	root.PersistentFlags().Bool("plain", false, "Output stable tab-separated text")
	root.PersistentFlags().Bool("no-color", false, "Disable ANSI colors")
	root.PersistentFlags().BoolP("quiet", "q", false, "Suppress non-essential output")
	root.PersistentFlags().IntVarP(&limit, "limit", "n", 20, "Maximum number of results")

	// Register subcommands.
	root.AddCommand(
		// Auth
		cmd.NewAuthCmd(opts),
		// User
		cmd.NewUserCmd(opts),
		cmd.NewWhoamiCmd(opts),
		// Tweets
		cmd.NewTweetsCmd(opts, &limit),
		cmd.NewTweetCmd(opts),
		cmd.NewTimelineCmd(opts, &limit),
		cmd.NewSearchCmd(opts, &limit),
		// Post actions
		cmd.NewPostCmd(opts),
		cmd.NewDeleteCmd(opts),
		// Reposts
		cmd.NewRepostCmd(opts),
		cmd.NewUnrepostCmd(opts),
		// Likes
		cmd.NewLikeCmd(opts),
		cmd.NewUnlikeCmd(opts),
		cmd.NewLikesCmd(opts, &limit),
		// Bookmarks
		cmd.NewBookmarksCmd(opts, &limit),
		cmd.NewBookmarkCmd(opts),
		cmd.NewUnbookmarkCmd(opts),
		// Relationships
		cmd.NewFollowCmd(opts),
		cmd.NewUnfollowCmd(opts),
		cmd.NewFollowersCmd(opts, &limit),
		cmd.NewFollowingCmd(opts, &limit),
		// Moderation
		cmd.NewBlockCmd(opts),
		cmd.NewUnblockCmd(opts),
		cmd.NewMuteCmd(opts),
		cmd.NewUnmuteCmd(opts),
		// Mentions
		cmd.NewMentionsCmd(opts, &limit),
	)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
