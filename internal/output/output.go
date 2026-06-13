package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/yashiels/twitter-cli/internal/types"
)

// Format controls how output is rendered.
type Format int

const (
	FormatHuman Format = iota
	FormatJSON
	FormatPlain
)

// Options holds global output configuration.
type Options struct {
	Format  Format
	Quiet   bool
	NoColor bool
}

// DefaultOptions returns sensible defaults based on the terminal.
func DefaultOptions() *Options {
	return &Options{
		Format:  FormatHuman,
		NoColor: !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()),
	}
}

// Printer writes formatted output.
type Printer struct {
	opts *Options
}

// New creates a new Printer with the given options.
func New(opts *Options) *Printer {
	if opts.NoColor {
		color.NoColor = true
	}
	return &Printer{opts: opts}
}

// JSON outputs any value as indented JSON.
func (p *Printer) JSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// PrintUser renders a user profile card.
func (p *Printer) PrintUser(u *types.User) error {
	switch p.opts.Format {
	case FormatJSON:
		return p.JSON(u)
	case FormatPlain:
		return p.printUserPlain(u)
	default:
		return p.printUserHuman(u)
	}
}

func (p *Printer) printUserHuman(u *types.User) error {
	bold := color.New(color.Bold)
	cyan := color.New(color.FgCyan)
	dim := color.New(color.Faint)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	// Header line: Name @handle [✓]
	verified := ""
	if u.Verified {
		verified = " " + yellow.Sprint("✓")
	}
	fmt.Printf("%s %s%s\n",
		bold.Sprint(u.Name),
		cyan.Sprint("@"+u.ScreenName),
		verified,
	)

	// Affiliation
	if u.Affiliation != "" {
		fmt.Printf("%s\n", dim.Sprint(u.Affiliation))
	}

	// Bio
	if u.Description != "" {
		wrapped := wordWrap(u.Description, 72)
		fmt.Println(wrapped)
	}

	// Location
	if u.Location != "" {
		fmt.Printf("📍 %s\n", u.Location)
	}

	fmt.Println()

	// Stats row
	fmt.Printf("%s followers  %s following  %s tweets\n",
		bold.Sprint(formatCount(u.FollowersCount)),
		bold.Sprint(formatCount(u.FriendsCount)),
		bold.Sprint(formatCount(u.StatusesCount)),
	)

	// Profile URL
	if u.ScreenName != "" {
		fmt.Printf("%s\n", dim.Sprint("https://x.com/"+u.ScreenName))
	}
	_ = green // used for status elsewhere
	return nil
}

func (p *Printer) printUserPlain(u *types.User) error {
	verified := "false"
	if u.Verified {
		verified = "true"
	}
	fields := []string{
		u.ScreenName,
		u.Name,
		u.Description,
		u.Location,
		fmt.Sprintf("%d", u.FollowersCount),
		fmt.Sprintf("%d", u.FriendsCount),
		fmt.Sprintf("%d", u.StatusesCount),
		verified,
		u.Affiliation,
	}
	fmt.Println(strings.Join(fields, "\t"))
	return nil
}

// PrintUsers renders a list of user profiles (for search --users results).
func (p *Printer) PrintUsers(users []*types.User) error {
	switch p.opts.Format {
	case FormatJSON:
		return p.JSON(users)
	case FormatPlain:
		for _, u := range users {
			if err := p.printUserPlain(u); err != nil {
				return err
			}
		}
		return nil
	default:
		dim := color.New(color.Faint)
		for i, u := range users {
			if i > 0 {
				fmt.Println(dim.Sprint(strings.Repeat("─", 72)))
			}
			if err := p.printUserHuman(u); err != nil {
				return err
			}
		}
		return nil
	}
}

// PrintTweets renders a list of tweets.
func (p *Printer) PrintTweets(tweets []*types.Tweet) error {
	switch p.opts.Format {
	case FormatJSON:
		return p.JSON(tweets)
	case FormatPlain:
		return p.printTweetsPlain(tweets)
	default:
		return p.printTweetsHuman(tweets)
	}
}

func (p *Printer) printTweetsHuman(tweets []*types.Tweet) error {
	dim := color.New(color.Faint)
	cyan := color.New(color.FgCyan)
	bold := color.New(color.Bold)
	yellow := color.New(color.FgYellow)

	for i, t := range tweets {
		if i > 0 {
			fmt.Println(dim.Sprint(strings.Repeat("─", 72)))
		}

		// Date and type indicator.
		dateStr := relativeTime(t.CreatedAt)
		typeStr := ""
		if t.IsRetweet {
			typeStr = " 🔁"
		} else if t.IsReply {
			typeStr = " 💬"
		}
		fmt.Printf("%s%s\n", dim.Sprint(dateStr), typeStr)

		// Tweet text.
		text := t.Text
		// Strip RT prefix for display.
		if strings.HasPrefix(text, "RT @") {
			parts := strings.SplitN(text, ": ", 2)
			if len(parts) == 2 {
				fmt.Printf("%s %s\n", dim.Sprint(parts[0]+":"), parts[1])
			} else {
				fmt.Println(text)
			}
		} else {
			fmt.Println(wordWrap(text, 72))
		}

		// Quoted tweet.
		if t.QuotedTweet != nil {
			qt := t.QuotedTweet
			qtHandle := ""
			if qt.AuthorHandle != "" {
				qtHandle = "@" + qt.AuthorHandle + ": "
			}
			qtText := wordWrap(qtHandle+qt.Text, 68)
			// Indent quoted tweet.
			for _, line := range strings.Split(qtText, "\n") {
				fmt.Printf("  %s\n", dim.Sprint("│ ")+line)
			}
		}

		// Engagement stats.
		stats := []string{}
		if t.FavoriteCount > 0 {
			stats = append(stats, yellow.Sprint("♥ ")+bold.Sprint(formatCount(t.FavoriteCount)))
		}
		if t.RetweetCount > 0 {
			stats = append(stats, "🔁 "+bold.Sprint(formatCount(t.RetweetCount)))
		}
		if t.ReplyCount > 0 {
			stats = append(stats, "💬 "+bold.Sprint(formatCount(t.ReplyCount)))
		}
		if t.ViewCount > 0 {
			stats = append(stats, dim.Sprint("👁 "+formatCount(t.ViewCount)))
		}
		if len(stats) > 0 {
			fmt.Println(strings.Join(stats, "  "))
		}

		// URL.
		if t.URL != "" {
			fmt.Println(cyan.Sprint(t.URL))
		}
	}
	fmt.Println()
	return nil
}

func (p *Printer) printTweetsPlain(tweets []*types.Tweet) error {
	for _, t := range tweets {
		fields := []string{
			t.ID,
			t.CreatedAt.Format(time.RFC3339),
			t.AuthorHandle,
			strings.ReplaceAll(t.Text, "\t", " "),
			strings.ReplaceAll(t.Text, "\n", " "),
			fmt.Sprintf("%d", t.FavoriteCount),
			fmt.Sprintf("%d", t.RetweetCount),
			fmt.Sprintf("%d", t.ReplyCount),
			fmt.Sprintf("%d", t.ViewCount),
			t.URL,
		}
		fmt.Println(strings.Join(fields, "\t"))
	}
	return nil
}

// Infof prints to stderr unless quiet.
func (p *Printer) Infof(format string, args ...any) {
	if !p.opts.Quiet {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}

// Errorf always prints to stderr.
func Errorf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
}

// formatCount formats large numbers as e.g. "12.3K" or "1.2M".
func formatCount(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}

// relativeTime returns a human-friendly relative time string.
func relativeTime(t time.Time) string {
	if t.IsZero() {
		return "unknown time"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	default:
		return t.Format("Jan 2, 2006")
	}
}

// wordWrap wraps text at word boundaries for the given column width.
func wordWrap(text string, width int) string {
	if utf8.RuneCountInString(text) <= width {
		return text
	}

	var lines []string
	for _, para := range strings.Split(text, "\n") {
		words := strings.Fields(para)
		if len(words) == 0 {
			lines = append(lines, "")
			continue
		}
		var line strings.Builder
		lineLen := 0
		for i, word := range words {
			wlen := utf8.RuneCountInString(word)
			if i == 0 {
				line.WriteString(word)
				lineLen = wlen
			} else if lineLen+1+wlen > width {
				lines = append(lines, line.String())
				line.Reset()
				line.WriteString(word)
				lineLen = wlen
			} else {
				line.WriteByte(' ')
				line.WriteString(word)
				lineLen += 1 + wlen
			}
		}
		if line.Len() > 0 {
			lines = append(lines, line.String())
		}
	}
	return strings.Join(lines, "\n")
}
