# twt CLI Specification

## Name & Purpose
`twt` — X/Twitter from the terminal. Built on X's private GraphQL API, reverse-engineered from the Android APK.

## Usage
```
twt [global flags] <command> [args]
```

## Commands (Phase 1 — Read Operations)

| Command | Description |
|---------|-------------|
| `twt auth login` | Extract cookies from Chrome or paste manually |
| `twt auth status` | Show current session |
| `twt auth logout` | Remove credentials |
| `twt whoami` | Show your own authenticated profile |
| `twt user <handle>` | View user profile |
| `twt tweets <handle>` | Latest original tweets |
| `twt tweet <id>` | Single tweet detail |
| `twt timeline` | Home timeline (For You) |
| `twt timeline --latest` | Chronological home timeline (Following) |
| `twt search <query>` | Search tweets |
| `twt search <query> --users` | Search users |
| `twt followers <handle>` | List a user's followers |
| `twt following <handle>` | List who a user follows |
| `twt likes [handle]` | List liked tweets for self or another user |
| `twt bookmarks` | Your bookmarked tweets |
| `twt mentions` | Your mention/notification timeline |

## Commands (Phase 2 — Write Operations)

| Command | Description |
|---------|-------------|
| `twt post <text>` | Post a new tweet (prompts for confirmation) |
| `twt post <text> --reply <id>` | Reply to a tweet |
| `twt post <text> --quote <id>` | Quote tweet |
| `twt post <text> --yes` | Post without confirmation prompt |
| `twt delete <id>` | Delete a tweet (prompts for confirmation) |
| `twt delete <id> --yes` | Delete without confirmation prompt |
| `twt repost <id>` | Repost (retweet) a tweet |
| `twt unrepost <id>` | Remove a repost |
| `twt follow <handle>` | Follow a user |
| `twt unfollow <handle>` | Unfollow a user |
| `twt like <id>` | Like a tweet |
| `twt unlike <id>` | Unlike a tweet |
| `twt bookmark <id>` | Add a bookmark |
| `twt unbookmark <id>` | Remove a bookmark |
| `twt block <handle>` | Block a user |
| `twt unblock <handle>` | Unblock a user |
| `twt mute <handle>` | Mute a user |
| `twt unmute <handle>` | Unmute a user |

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

| Command | Protocol | Operation | Query ID / Endpoint |
|---------|----------|-----------|---------------------|
| `twt user <handle>` | GraphQL GET | `GetUserByScreenNameQuery` | `bbS0COK9SwcgdM7QCEqWDg` |
| `twt tweets <handle>` | GraphQL GET | `UserProfileOriginalsTimelineQuery` | `xlAB_H3dvYL4q1C-PzM_ag` |
| `twt tweet <id>` | GraphQL GET | `GetPostById` | `lOsezlo57Y40B-TLgWqxEA` |
| `twt timeline` | GraphQL GET | `HomeTimeline` | `t_sH369wuH1CO5lbW2qlYg` |
| `twt search` | GraphQL GET | `SearchTimelineQuery` | `rxBGDmZrc-NcrXfcRNUdMg` |
| `twt likes [handle]` | GraphQL GET | `UserProfileFavoritesTimelineQuery` | `M34xxhtrHWGAxGofYSGclA` |
| `twt bookmarks` | GraphQL GET | `BookmarksTimelineQuery` | `DN0j17CihaEo7QYmaZGkiw` |
| `twt mentions` | GraphQL GET | `NotificationTimelineQuery` | `jsOzc8RhpUpH5InTskP6Yw` |
| `twt followers` | REST v1.1 GET | — | `/1.1/followers/list.json` |
| `twt following` | REST v1.1 GET | — | `/1.1/friends/list.json` |
| `twt post` | GraphQL POST | `CreatePost` | `vMia9QJ2JVkCXuO5J4MTbw` |
| `twt delete` | GraphQL POST | `DeletePostMutation` | `1EVIme6zMCgTO7F95wuElA` |
| `twt repost` | GraphQL POST | `CreateRepostMutation` | `ydMACa-dOjZx126SWo6q5A` |
| `twt unrepost` | GraphQL POST | `DeleteRepostMutation` | `w1Bo2Whh4f4lha5Djgnpvg` |
| `twt follow` | GraphQL POST | `FollowUser` | `44lRL9CTLTxi4aAMSqAmVw` |
| `twt unfollow` | GraphQL POST | `UnfollowUser` | `zpWrwHHfa_6sKBQr6SGCwg` |
| `twt like` | GraphQL POST | `FavoriteMutation` | `awITBmMVajjvqY2wTL8DUw` |
| `twt unlike` | GraphQL POST | `UnfavoriteMutation` | — |
| `twt bookmark` | GraphQL POST | `BookmarkAddMutation` | `IjefskW4Kr2i-6XRdshQEg` |
| `twt unbookmark` | GraphQL POST | `BookmarkRemoveMutation` | `K5KIqVnds5iJ00WdHb8Nmw` |
| `twt block` | GraphQL POST | `BlockUser` | `8zl3cVULtte29uCoWREtBQ` |
| `twt unblock` | GraphQL POST | `UnblockUser` | `WtUZ-1fkiAJGXfN6gAwrLw` |
| `twt mute` | GraphQL POST | `MuteUser` | `LoZAfbPr53jnw9Y2FydOIQ` |
| `twt unmute` | GraphQL POST | `UnmuteUser` | `29vlsCe7kkuB4JKnQGeK5w` |

GraphQL requests go to:
```
https://api.twitter.com/graphql/{queryId}/{operationName}
```
Reads are GET; writes are POST with `application/json`.

REST v1.1 requests go to:
```
https://api.twitter.com{endpoint}
```

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
