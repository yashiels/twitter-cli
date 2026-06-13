# twt — X/Twitter CLI

A fast command-line client for X/Twitter using the private GraphQL API reverse-engineered from the Android APK. No API key required.

## Install

```sh
brew install yashiels/tap/twt
```

Or build from source:

```sh
go install github.com/yashiels/twitter-cli/cmd/twt@latest
```

Or clone and build:

```sh
git clone https://github.com/yashiels/twitter-cli.git
cd twitter-cli
make build
```

## Quick Start

```sh
twt auth login            # auto-extracts cookies from Chrome
twt whoami                # show your profile
twt timeline --latest     # your chronological feed
twt post "hello world"    # post a tweet (confirms before posting)
twt search "openclaw"     # search tweets
twt user steipete         # view another profile
```

## Authentication

```sh
# Auto-extract from Chrome (macOS — recommended)
twt auth login

# Manual entry
twt auth login --manual

# Check session
twt auth status

# Remove credentials
twt auth logout
```

Credentials are stored at `~/.config/twt/credentials.json` (mode 0600).

### Getting cookies manually

1. Open Chrome → **x.com** → log in
2. DevTools → Application → Cookies → `https://x.com`
3. Copy `auth_token` and `ct0`
4. `twt auth login --manual`

## Commands

### Identity

```sh
twt whoami                       # show your own profile
twt auth status                  # show session + stored user ID
```

### Profiles

```sh
twt user steipete
twt user steipete --json
```

### Tweets

```sh
twt tweets steipete              # latest original tweets
twt tweets steipete --limit 50   # more tweets
twt tweet 2065650561484267540    # single tweet by ID
```

### Post / Delete

```sh
twt post "hello world"                          # new tweet (confirms)
twt post "hello world" --yes                    # skip confirmation
twt post --reply 1234567890 "nice thread"       # reply to a tweet
twt post --quote 1234567890 "interesting take"  # quote-tweet
twt delete 1234567890                           # delete (confirms)
twt delete 1234567890 --yes                     # skip confirmation
```

### Repost

```sh
twt repost 1234567890
twt unrepost 1234567890
```

### Bookmarks

```sh
twt bookmarks                    # list your bookmarks
twt bookmark 1234567890          # save a tweet
twt unbookmark 1234567890        # remove a bookmark
```

### Follow / Unfollow / Followers / Following

```sh
twt follow steipete
twt unfollow steipete
twt followers steipete           # list steipete's followers
twt following steipete           # list who steipete follows
```

### Like / Unlike / Likes

```sh
twt like 1234567890
twt unlike 1234567890
twt likes                        # your liked tweets
twt likes steipete               # steipete's liked tweets
```

### Mentions

```sh
twt mentions                     # your notification timeline
twt mentions --limit 5
```

### Moderation

```sh
twt block someuser
twt unblock someuser
twt mute someuser
twt unmute someuser
```

### Search

```sh
twt search "openclaw"            # search tweets
twt search "steipete" --users    # search users
twt search "AI agents" --limit 5
```

### Timeline

```sh
twt timeline                     # For You feed
twt timeline --latest            # Following (chronological)
twt timeline --latest --limit 10
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--json` | | Machine-readable JSON output |
| `--plain` | | Tab-separated output for scripting |
| `--no-color` | | Disable ANSI colours |
| `--quiet` | `-q` | Suppress progress messages |
| `--limit` | `-n` | Max results (default 20) |
| `--version` | `-v` | Print version |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid usage |
| 3 | Not authenticated |
| 4 | Rate limited |
| 5 | Not found |

## Rate Limiting

Built-in safety to avoid bans:

- **2–5s random jitter** between requests
- **Auto-retry** on HTTP 429 with exponential backoff
- **Header respect** — reads `x-rate-limit-remaining` / `x-rate-limit-reset`
- **Warnings** on stderr when ≤10 requests remaining

## Development

```sh
# Setup hooks (runs lint/fmt/vet before each commit)
make hooks

# Build
make build

# Test
make test

# Lint
make lint

# Format
make fmt
```

### Architecture

```
cmd/twt/main.go              # Cobra root, command wiring
internal/
  api/
    client.go                # HTTP client, auth headers, jitter, 429 retry, restGet()
    graphql.go               # GraphQL GET/POST helpers
    user.go                  # GetUserByScreenNameQuery, GetUserQuery (by ID)
    tweets.go                # UserProfileOriginalsTimelineQuery
    tweet.go                 # GetPostById
    post.go                  # CreatePost, CreateReply, CreateQuote, DeletePost
    repost.go                # CreateRepostMutation, DeleteRepostMutation
    bookmark.go              # BookmarksTimelineQuery, BookmarkAdd/RemoveMutation
    relationships.go         # GetFollowers, GetFollowing (REST v1.1)
    moderation.go            # BlockUser, UnblockUser, MuteUser, UnmuteUser
    likes.go                 # UserProfileFavoritesTimelineQuery
    mentions.go              # NotificationTimelineQuery
    follow.go                # FollowUser / UnfollowUser mutations
    like.go                  # FavoriteMutation / UnfavoriteMutation
    search.go                # SearchTimelineQuery
    timeline.go              # HomeTimeline / HomeTimelineLatest
    verify.go                # v1.1 verify_credentials (returns handle + user ID)
  auth/
    store.go                 # Credential storage (auth_token, ct0, user_id, handle)
    chrome.go                # macOS Chrome cookie extraction (Keychain + AES)
  cmd/                       # Cobra command implementations
  output/output.go           # Human / JSON / plain formatters
  types/                     # User and Tweet domain structs
```

### How it works

The CLI uses Twitter's private Android GraphQL API (`api.twitter.com/graphql/`), reverse-engineered from the X Android APK v11.99.0. Authentication is via browser session cookies (`auth_token` + `ct0`), not OAuth app credentials.

## License

MIT
