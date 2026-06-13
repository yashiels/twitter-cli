# Changelog

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
