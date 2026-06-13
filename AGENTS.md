# AGENTS.md

Telegraph style. Imperative. Min tokens.

## Project

- `twt` — Go CLI for X/Twitter private GraphQL API. No developer API key.
- Module: `github.com/yashiels/twitter-cli`
- Binary: `twt`
- License: MIT

## Build / Test / Lint

```sh
go build -o twt ./cmd/twt
go test ./...
golangci-lint run --timeout=5m
```

Pre-commit hook runs gofumpt + goimports + go vet + golangci-lint. Activate: `make hooks`.

## Architecture

```
cmd/twt/main.go        — Cobra root, command wiring
internal/api/          — HTTP client, GraphQL GET/POST, rate limiter, all endpoint methods
internal/auth/         — Chrome cookie extraction (macOS Keychain + AES-128-CBC), credential store
internal/cmd/          — Cobra commands (auth, user, tweets, tweet, follow, like, search, timeline)
internal/output/       — Human / JSON / plain formatters, color, TTY detection
internal/types/        — User, Tweet domain structs
```

## API Contract

- Base: `https://api.twitter.com/graphql/{queryId}/{operationName}`
- Auth: `auth_token` + `ct0` cookies from Chrome, NOT OAuth app credentials
- Reads: GET with `variables` + `features` query params
- Writes: POST with JSON body `{"variables": ..., "features": ...}`
- User-Agent: `TwitterAndroid/11.99.0`
- APK schema, NOT web schema. `legacy` is mostly empty. Use `details.*`, `counts.*`, `relationship_counts.*`, `tweet_counts.*`

## Rate Limiting

- 2–5s random jitter before every request
- Auto-retry on 429 with `x-rate-limit-reset` wait (up to 2 retries)
- Warn stderr when ≤10 requests remaining
- Give up if reset >5 min away

## Commits

- Conventional Commits: `type(scope): summary`
- Trailer required: `Co-authored-by: Yashiel Sookdeo <yashiel@skyner.co.za>`
- Stage specific files. Never `git add .`
- Run `make lint` before committing if hooks aren't active

## Prohibited

- No Twitter developer API keys or OAuth 2.0 flows
- No CGO dependencies beyond existing go.mod
- No secrets, tokens, or cookie data in commits
- No modifying Chrome's Cookies DB
- No silently swallowing 429 errors
