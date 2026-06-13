# twt — X/Twitter CLI

A fast command-line client for X/Twitter using the private GraphQL API reverse-engineered from the Android APK.

## Install

```sh
go install github.com/yashiels/twitter-cli/cmd/twt@latest
# or clone and build:
make build
```

## Authentication

```sh
# Auto-extract from Chrome (macOS — recommended)
twt auth login

# Manual entry (enter auth_token and ct0 interactively)
twt auth login --manual

# Check current auth state (exit 3 if not logged in)
twt auth status

# Remove stored credentials
twt auth logout
```

Credentials are stored at `~/.config/twt/credentials.json` (mode 0600).

### Getting cookies manually

1. Open Chrome, go to **x.com** and log in.
2. Open DevTools → Application → Cookies → `https://x.com`.
3. Copy the values of `auth_token` and `ct0`.
4. Run `twt auth login --manual` and paste them when prompted.

## Commands

### User profile

```sh
twt user steipete
twt user @elonmusk

# JSON output
twt user steipete --json

# Tab-separated (for scripting)
twt user steipete --plain
```

### Tweets

```sh
twt tweets steipete
twt tweets @steipete --limit 50

# JSON output
twt tweets steipete --json

# Plain tab-separated
twt tweets steipete --plain
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--json` | | Machine-readable JSON output |
| `--plain` | | Tab-separated output for scripting |
| `--no-color` | | Disable ANSI colours |
| `--quiet` | `-q` | Suppress progress/info messages |
| `--limit` | `-n` | Max results (default 20) |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 3 | Not authenticated |
| 5 | User not found |

## Architecture

```
twitter-cli/
├── cmd/twt/main.go           # Entry point, cobra root
├── internal/
│   ├── api/
│   │   ├── client.go         # HTTP client, auth headers, rate limiting
│   │   ├── graphql.go        # GraphQL GET request helper
│   │   ├── user.go           # GetUserByScreenNameQuery
│   │   ├── tweets.go         # UserProfileOriginalsTimelineQuery
│   │   └── verify.go         # v1.1 verify_credentials
│   ├── auth/
│   │   ├── store.go          # Credential storage (~/.config/twt/)
│   │   └── chrome.go         # macOS Chrome cookie extraction
│   ├── cmd/
│   │   ├── auth.go           # auth login/status/logout
│   │   ├── user.go           # user command
│   │   └── tweets.go         # tweets command
│   ├── output/
│   │   └── output.go         # human/json/plain formatters
│   └── types/
│       ├── user.go           # User struct
│       └── tweet.go          # Tweet struct
```

## API Details

Uses Twitter's private Android GraphQL API:

- Bearer token: public Android client bearer (same for all clients)
- Auth: `auth_token` + `ct0` cookies extracted from Chrome or entered manually
- Rate limiting: reads `x-rate-limit-remaining` / `x-rate-limit-reset` headers, adds jitter between requests

## License

MIT
