# twt CLI Specification

## Name & Purpose
`twt` — X/Twitter from the terminal. Built on X's private GraphQL API, reverse-engineered from the Android APK.

## Usage
```
twt [global flags] <command> [args]
```

## Commands (Phase 1 — Implemented)

| Command | Description |
|---------|-------------|
| `twt auth login` | Extract cookies from Chrome or paste manually |
| `twt auth status` | Show current session |
| `twt auth logout` | Remove credentials |
| `twt user <handle>` | View user profile |
| `twt tweets <handle>` | Latest original tweets |

## Commands (Future)

| Command | Description |
|---------|-------------|
| `twt tweet <id>` | Single tweet detail |
| `twt follow <handle>` | Follow user |
| `twt unfollow <handle>` | Unfollow user |
| `twt like <tweet-id>` | Like tweet |
| `twt unlike <tweet-id>` | Unlike tweet |
| `twt search <query>` | Search tweets |
| `twt timeline` | Home timeline |
| `twt timeline --latest` | Chronological home timeline |
| `twt bookmarks` | Your bookmarks |
| `twt bookmark <id>` | Add bookmark |

## Global Flags

| Flag | Description |
|------|-------------|
| `-h, --help` | Show help |
| `--version` | Print version |
| `-q, --quiet` | Suppress non-essential output |
| `--json` | Machine-readable JSON |
| `--plain` | Stable tab-separated text |
| `--no-color` | Disable ANSI colors |
| `-n, --limit <int>` | Max results (default: 20) |

## Auth Model

- **Primary**: auto-extract `auth_token` + `ct0` from Chrome (macOS Keychain decryption of the Cookies SQLite DB)
- **Fallback**: manual entry via `twt auth login --manual`
- **Storage**: `~/.config/twt/credentials.json` (mode `0600`)
- **Env override**: `TWL_AUTH_TOKEN`, `TWL_CT0`

Credential priority order: env vars → `credentials.json` → Chrome extraction.

## API Mappings

| Command | GraphQL Operation | Query ID |
|---------|-------------------|----------|
| `twt user <handle>` | `GetUserByScreenNameQuery` | `bbS0COK9SwcgdM7QCEqWDg` |
| `twt tweets <handle>` | `UserProfileOriginalsTimelineQuery` | `xlAB_H3dvYL4q1C-PzM_ag` |
| `twt tweet <id>` | `GetPostById` | `lOsezlo57Y40B-TLgWqxEA` |
| `twt follow` | `FollowUser` | `44lRL9CTLTxi4aAMSqAmVw` |
| `twt unfollow` | `UnfollowUser` | `zpWrwHHfa_6sKBQr6SGCwg` |
| `twt like` | `FavoriteMutation` | `awITBmMVajjvqY2wTL8DUw` |
| `twt search` | `SearchTimelineQuery` | `rxBGDmZrc-NcrXfcRNUdMg` |
| `twt timeline` | `HomeTimeline` | `t_sH369wuH1CO5lbW2qlYg` |
| `twt bookmarks` | `BookmarksTimelineQuery` | `DN0j17CihaEo7QYmaZGkiw` |

All requests go to:
```
https://api.twitter.com/graphql/{queryId}/{operationName}
```
Reads are GET; writes are POST with `application/json`.

## Rate Limiting Strategy

1. **Jitter**: 2–5 s random sleep between consecutive requests
2. **Auto-retry**: on HTTP 429, wait for `x-rate-limit-reset` and retry (up to 2 retries)
3. **Header respect**: read `x-rate-limit-remaining` and `x-rate-limit-reset` after every response
4. **Warning**: print to stderr when `x-rate-limit-remaining` ≤ 10
5. **Give up**: if `x-rate-limit-reset` is more than 5 minutes in the future, exit with code 4

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Generic error |
| `2` | Invalid usage / bad flags |
| `3` | Auth required or expired |
| `4` | Rate limited |
| `5` | User or tweet not found |

## Example Invocations

```sh
# Authenticate using Chrome cookies
twt auth login

# Check session
twt auth status

# Look up a profile
twt user jack

# Get last 10 tweets as JSON
twt tweets jack --limit 10 --json

# Plain tab-separated output (good for piping to awk/cut)
twt tweets jack --plain

# Quiet mode — suppress banners and warnings
twt user jack --quiet

# No colour (for CI / log capture)
twt tweets jack --no-color
```

## Response Schema Notes

The APK GraphQL schema differs from the public v2 REST API:

- Follower/following counts: `relationship_counts.followers` / `relationship_counts.following`
- Tweet/like counts: `tweet_counts.tweets` / `tweet_counts.likes`
- Profile details: `details.name`, `details.description`, `details.location`, `details.url`
- Verified badge: `details.is_blue_verified`
- Tweet entries: nested under `timeline_v2.timeline.instructions[].entries[].content.itemContent.tweet_results.result`
- View counts: `views.count` on the tweet result object
