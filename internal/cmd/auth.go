package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/yashiels/twitter-cli/internal/api"
	"github.com/yashiels/twitter-cli/internal/auth"
	"github.com/yashiels/twitter-cli/internal/output"
)

// NewAuthCmd returns the "auth" parent command with login/logout/status subcommands.
func NewAuthCmd(opts *output.Options) *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
		Long:  "Authenticate with X/Twitter. Extracts cookies from Chrome by default.",
	}

	authCmd.AddCommand(
		newAuthLoginCmd(opts),
		newAuthStatusCmd(opts),
		newAuthLogoutCmd(opts),
	)

	return authCmd
}

func newAuthLoginCmd(opts *output.Options) *cobra.Command {
	var manual bool

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to X/Twitter",
		Long: `Extract auth_token and ct0 cookies from Chrome (macOS).
Use --manual to enter cookies interactively instead.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p := output.New(opts)

			var authToken, ct0 string

			if !manual {
				p.Infof("Extracting cookies from Chrome...")
				cookies, err := auth.ExtractFromChrome()
				if err != nil {
					p.Infof("Chrome extraction failed: %v", err)
					p.Infof("Falling back to manual entry.")
					manual = true
				} else {
					authToken = cookies.AuthToken
					ct0 = cookies.CT0
					p.Infof("Chrome cookies extracted successfully.")
				}
			}

			if manual {
				var err error
				authToken, ct0, err = promptCredentials()
				if err != nil {
					return err
				}
			}

			creds := &auth.Credentials{
				AuthToken: authToken,
				CT0:       ct0,
			}

			// Verify credentials and resolve identity.
			p.Infof("Verifying credentials...")
			client := api.NewClient(creds)
			identity, err := verifyAndGetIdentity(client)
			if err != nil {
				return fmt.Errorf("credentials invalid or API error: %w", err)
			}
			creds.Handle = identity.Handle
			creds.UserID = identity.UserID

			if err := auth.Save(creds); err != nil {
				return fmt.Errorf("save credentials: %w", err)
			}

			if creds.Handle != "" {
				fmt.Printf("✓ Logged in as @%s\n", creds.Handle)
			} else {
				fmt.Println("✓ Logged in (handle unavailable)")
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&manual, "manual", false, "Skip Chrome auto-detection and enter cookies manually")
	return cmd
}

func newAuthStatusCmd(opts *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current auth state",
		RunE: func(cmd *cobra.Command, args []string) error {
			creds, err := auth.Load()
			if errors.Is(err, auth.ErrNotAuthenticated) {
				fmt.Println("Not authenticated. Run: twt auth login")
				os.Exit(3)
			}
			if err != nil {
				return err
			}

			handle := creds.Handle
			if handle == "" {
				handle = "(unknown)"
			}

			p := output.New(opts)
			if opts.Format == output.FormatJSON {
				return p.JSON(map[string]any{
					"authenticated": true,
					"handle":        handle,
					"user_id":       creds.UserID,
					"saved_at":      creds.SavedAt,
				})
			}

			fmt.Printf("Logged in as @%s\n", handle)
			if creds.UserID != "" {
				fmt.Printf("User ID: %s\n", creds.UserID)
			}
			if !creds.SavedAt.IsZero() {
				fmt.Printf("Token saved: %s\n", creds.SavedAt.Format("2006-01-02 15:04:05"))
			}
			return nil
		},
	}
}

func newAuthLogoutCmd(_ *output.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove stored credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := auth.Delete(); err != nil {
				return err
			}
			fmt.Println("Logged out. Credentials removed.")
			return nil
		},
	}
}

// promptCredentials reads auth_token and ct0 interactively.
func promptCredentials() (authToken, ct0 string, err error) {
	fmt.Fprint(os.Stderr, "auth_token (from .x.com cookies): ")
	raw, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		// Fallback for non-terminal.
		fmt.Fscanln(os.Stdin, &authToken)
	} else {
		fmt.Fprintln(os.Stderr)
		authToken = strings.TrimSpace(string(raw))
	}

	fmt.Fprint(os.Stderr, "ct0 (CSRF token from .x.com cookies): ")
	raw, err = term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fscanln(os.Stdin, &ct0)
	} else {
		fmt.Fprintln(os.Stderr)
		ct0 = strings.TrimSpace(string(raw))
	}

	if authToken == "" || ct0 == "" {
		return "", "", fmt.Errorf("both auth_token and ct0 are required")
	}
	return authToken, ct0, nil
}

// identityResult holds the resolved identity from verify.
type identityResult struct {
	Handle string
	UserID string
}

// verifyAndGetIdentity calls verify_credentials to confirm auth and resolve identity.
func verifyAndGetIdentity(client *api.Client) (*identityResult, error) {
	result, err := client.VerifyCredentials()
	if err != nil {
		// Fall back to confirming API works by checking a known user.
		_, ferr := client.GetUserByScreenName("twitter")
		if ferr != nil {
			return nil, fmt.Errorf("API verification failed: %w", err)
		}
		return &identityResult{}, nil // credentials work but we can't get identity
	}
	return &identityResult{
		Handle: result.ScreenName,
		UserID: result.UserID,
	}, nil
}
