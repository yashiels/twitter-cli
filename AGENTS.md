# AGENTS.md — twt (twitter-cli)

## Project Purpose
`twt` is a Go CLI for X/Twitter built on the **private GraphQL API**, reverse-engineered from the Android APK. It extracts credentials directly from Chrome on macOS (no developer API key needed).

## Module
`github.com/yashiels/twitter-cli`

## Binary
`twt`

## Key Commands
```sh
go build -o twt ./cmd/twt   # build
go test ./...                # test
golangci-lint run --timeout=5m   # lint
```

## Architecture

```
cmd/twt/          — Cobra entry point (main.go)
internal/
  api/            — HTTP client, GraphQL request/response, rate limiting
  auth/           — Chrome cookie extraction (macOS Keychain), credential store
  cmd/            — Cobra sub-commands (auth, user, tweets, …)
  output/         — Formatters: pretty, --json, --plain
  types/          — Domain structs (User, Tweet, Credentials, …)
```

### Package responsibilities

| Package | Role |
|---------|------|
| `internal/api` | Builds and fires GraphQL requests; handles jitter, 429 retry, and x-rate-limit-* headers |
| `internal/auth` | Decrypts Chrome's `Cookies` SQLite DB via macOS Keychain; stores `~/.config/twt/credentials.json` |
| `internal/cmd` | Cobra command tree; delegates to api + output |
| `internal/output` | ANSI-rich table, `--json` (raw struct), `--plain` (tab-separated) |
| `internal/types` | Shared value types — no external dependencies |

## Auth Model
- **Primary**: auto-extract `auth_token` + `ct0` from Chrome profile on macOS (Keychain decryption)
- **Fallback**: `twt auth login --manual` — paste tokens directly
- **Storage**: `~/.config/twt/credentials.json` (mode `0600`)
- **Env override**: `TWL_AUTH_TOKEN`, `TWL_CT0`

## API
- **Base URL**: `https://api.twitter.com/graphql/{queryId}/{operationName}`
- All reads are GET; all writes are POST with `application/json` body
- Required headers mimic the Android APK: `Authorization: Bearer <bearer>`, `x-twitter-auth-type`, `x-twitter-client-language`, etc.

## Rate Limiting
- 2–5 s random jitter between requests
- Auto-retry on HTTP 429 (up to 2 retries)
- Respects `x-rate-limit-remaining` / `x-rate-limit-reset`
- Warns to stderr when ≤ 10 requests remain
- Gives up if the reset window is > 5 minutes

## Response Parsing
Responses follow the APK schema:
- User stats in `relationship_counts` / `tweet_counts`
- User detail fields under `details.*`
- Tweets in `timeline_v2.timeline.instructions[].entries[].content.itemContent.tweet_results.result`

## Commit Rules
- Follow **Conventional Commits**: `type(scope): summary #issue`
- Every commit body must include the trailer:
  ```
  Co-authored-by: Yashiel Sookdeo <yashiel@skyner.co.za>
  ```
- **Never commit** secrets, API keys, bearer tokens, `credentials.json`, or cookie database files.
- Stage specific files only — never `git add .`

## What Not To Do
- Do not add a Twitter developer API key or OAuth 2.0 flow — this tool uses the private GraphQL API.
- Do not modify the Chrome Cookies DB directly.
- Do not introduce CGO dependencies beyond what is already in go.mod.
- Do not silently swallow 429 errors — surface them via the rate-limit warning path.
