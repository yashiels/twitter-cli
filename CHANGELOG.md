# Changelog

## [0.2.0] - 2026-06-14

### Added
- `twt post <text>` — post a new tweet with confirmation prompt (`--yes` to skip)
- `twt post --reply <id> <text>` — reply to a tweet
- `twt post --quote <id> <text>` — quote-tweet a tweet
- `twt delete <id>` — delete a tweet with confirmation (`--yes` to skip)
- `twt repost <id>` — repost (retweet) a tweet
- `twt unrepost <id>` — remove a repost
- `twt bookmarks` — list your bookmarked tweets
- `twt bookmark <id>` — bookmark a tweet
- `twt unbookmark <id>` — remove a bookmark
- `twt followers <handle>` — list a user's followers (REST v1.1)
- `twt following <handle>` — list who a user follows (REST v1.1)
- `twt mentions` — show recent mentions/notifications
- `twt likes [handle]` — list liked tweets (your own or another user's)
- `twt block <handle>` — block a user
- `twt unblock <handle>` — unblock a user
- `twt mute <handle>` — mute a user
- `twt unmute <handle>` — unmute a user
- `twt whoami` — show your own profile

### Changed
- `twt auth login` now stores user ID and handle in credentials for commands that need them
- `twt auth status` now shows stored user ID

## [0.1.0] - 2026-06-13

### Added
- Initial release
- `twt auth login` — auto-extract Chrome cookies on macOS, manual fallback
- `twt auth status` — show current session
- `twt auth logout` — remove credentials
- `twt user <handle>` — profile lookup with stats, bio, verified badge
- `twt tweets <handle>` — timeline with dates, engagement stats, view counts
- `--json`, `--plain`, `--no-color`, `--quiet`, `--limit` global flags
- Built-in rate limiting: 2-5s jitter, auto-retry on 429
- APK-schema response parsing (relationship_counts, tweet_counts, details.*)
