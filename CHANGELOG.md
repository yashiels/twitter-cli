# Changelog

## [0.2.0] - 2026-06-14

### Added
- `twt post <text>` — post a new tweet with confirmation prompt; `--yes` skips it
- `twt post --reply <id>` — reply to a tweet
- `twt post --quote <id>` — quote tweet
- `twt delete <id>` — delete a tweet (confirmation prompt; `--yes` skips it)
- `twt repost <id>` / `twt unrepost <id>` — repost and undo repost
- `twt bookmarks` — list your bookmarked tweets
- `twt bookmark <id>` / `twt unbookmark <id>` — add and remove bookmarks
- `twt followers <handle>` / `twt following <handle>` — list followers and following
- `twt likes [handle]` — list liked tweets for yourself or another user
- `twt mentions` — notification/mention timeline
- `twt block <handle>` / `twt unblock <handle>` — block and unblock users
- `twt mute <handle>` / `twt unmute <handle>` — mute and unmute users
- `twt whoami` — show your own authenticated profile
- `twt follow <handle>` / `twt unfollow <handle>` — follow and unfollow users
- `twt like <id>` / `twt unlike <id>` — like and unlike tweets
- `twt search <query>` — search tweets and users
- `twt timeline` / `twt timeline --latest` — home timeline (For You / Following)
- `twt tweet <id>` — single tweet detail by ID

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
