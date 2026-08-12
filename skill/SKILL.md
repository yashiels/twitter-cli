---
name: twt
description: Read and act on X/Twitter from the terminal — timeline, search, tweets, profiles, mentions. Posting, liking, reposting, following, blocking and deleting are confirm-gated write actions.
---

# twt

`twt` is a fast command-line client for X/Twitter that talks directly to Twitter's private Android GraphQL API (reverse-engineered from the X Android APK) — no developer API key or OAuth app needed. It authenticates with your existing browser session cookies (`auth_token` + `ct0`), which it can auto-extract from Chrome on macOS or take by manual entry. Reads (timeline, search, profiles, tweets, mentions) are safe to run anytime; writes (post, like, repost, follow, block, mute, bookmark, delete) act on the live account and must only run when the user explicitly asks.

## Install

```sh
brew install yashiels/tap/twt
```

Or build from source (Go): `go install github.com/yashiels/twitter-cli/cmd/twt@latest`.

## Auth

Authentication is via X/Twitter session cookies, **not** OAuth app credentials. `twt` needs two cookie values: `auth_token` and `ct0` (the CSRF token). There are two ways to supply them:

- **Chrome auto-extraction (macOS, default)** — `twt auth login` reads Chrome's encrypted cookie store (`~/Library/Application Support/Google/Chrome/Default/Cookies`, then `Profile 1`, then Chromium), decrypting with the AES key from the "Chrome Safe Storage" macOS Keychain entry. Requires a logged-in `x.com` session in a local Chrome/Chromium profile.
- **Manual entry** — `twt auth login --manual` prompts interactively (hidden input) for `auth_token` and `ct0`. Get them from Chrome DevTools → Application → Cookies → `https://x.com`.

On login, `twt` verifies the cookies against the API, resolves your handle/user ID, and stores everything at `~/.config/twt/credentials.json` (mode `0600`). There are **no environment variables or flags** for passing tokens directly — credentials come only from `auth login` or the on-disk file.

| Command | Description |
|---------|-------------|
| `twt auth login` | Auto-extract cookies from Chrome (macOS) and log in |
| `twt auth login --manual` | Skip Chrome; enter `auth_token` + `ct0` interactively |
| `twt auth status` | Show current auth state (handle, token save time) |
| `twt auth logout` | Remove stored credentials |
| `twt whoami` | Show your own authenticated profile (`--json` supported) |

## Read commands (safe, unattended-friendly)

| Command | Description |
|---------|-------------|
| `twt timeline` | Your home "For You" feed |
| `twt timeline --latest` | Following feed, chronological |
| `twt search <query>` | Search tweets |
| `twt search <query> --users` | Search users instead of tweets |
| `twt user <handle>` | Look up a user profile |
| `twt tweets <handle>` | Recent original tweets from a user |
| `twt tweet <tweet-id>` | Show a single tweet by ID |
| `twt mentions` | Your mentions / notification timeline |
| `twt likes [handle]` | Liked tweets (yours, or another user's) |
| `twt bookmarks` | List your bookmarked tweets |
| `twt followers <handle>` | List a user's followers |
| `twt following <handle>` | List who a user follows |

## Write / action commands (confirm-gated)

These mutate the live account. `post` and `delete` prompt for confirmation and take `--yes` to skip it; the rest execute immediately with no built-in prompt, so an agent must gate them on explicit user intent.

| Command | Description | Built-in prompt |
|---------|-------------|-----------------|
| `twt post <text>` | New tweet; `--reply <id>` / `--quote <id>`; `--yes` skips prompt | Yes (`--yes`) |
| `twt delete <tweet-id>` | Delete one of your tweets; `--yes` skips prompt | Yes (`--yes`) |
| `twt like <tweet-id>` | Like a tweet | No |
| `twt unlike <tweet-id>` | Remove a like | No |
| `twt repost <tweet-id>` | Repost (retweet) | No |
| `twt unrepost <tweet-id>` | Remove a repost | No |
| `twt bookmark <tweet-id>` | Add a bookmark | No |
| `twt unbookmark <tweet-id>` | Remove a bookmark | No |
| `twt follow <handle>` | Follow a user | No |
| `twt unfollow <handle>` | Unfollow a user | No |
| `twt block <handle>` | Block a user | No |
| `twt unblock <handle>` | Unblock a user | No |
| `twt mute <handle>` | Mute a user | No |
| `twt unmute <handle>` | Unmute a user | No |

## Global flags

| Flag | Short | Description |
|------|-------|-------------|
| `--json` | | Machine-readable JSON output |
| `--plain` | | Stable tab-separated output for scripting |
| `--no-color` | | Disable ANSI colours |
| `--quiet` | `-q` | Suppress non-essential output |
| `--limit` | `-n` | Max results (default 20) |
| `--version` | `-v` | Print version |

Exit codes: `0` success · `1` general error · `2` invalid usage · `3` not authenticated · `4` rate limited · `5` not found.

## Headless / agent usage

- **Reads are safe to run unattended.** `timeline`, `search`, `user`, `tweets`, `tweet`, `mentions`, `likes`, `bookmarks`, `followers`, `following`, `whoami`, and `auth status` never change state — run these freely, ideally with `--json` for parsing.
- **Writes are confirm-gated by policy.** Only run `post`, `delete`, `like`/`unlike`, `repost`/`unrepost`, `bookmark`/`unbookmark`, `follow`/`unfollow`, `block`/`unblock`, `mute`/`unmute` when the user has explicitly asked for that specific action. Note that only `post` and `delete` have an interactive confirmation prompt (pass `--yes` to skip); every other write executes immediately with no prompt, so the agent itself is the only guard — do not run them speculatively.
- **Auth cannot be bootstrapped remotely via Chrome.** `twt auth login` (Chrome auto-extraction) needs a logged-in `x.com` session in a local Chrome/Chromium profile plus macOS Keychain access — it will not work on a headless box with no browser profile. The headless fallback is the manual `auth_token` + `ct0` path.
- **Seeding credentials headless.** There are no env vars/flags for tokens. Either (a) run `twt auth login` / `--manual` once interactively on the machine so `~/.config/twt/credentials.json` is populated, or (b) place a valid `credentials.json` (JSON `{"auth_token": "...", "ct0": "..."}`, mode `0600`) at `~/.config/twt/` directly. After that, every command works unattended until the cookies expire — at which point re-auth requires the user. If a command exits `3` (not authenticated), surface that rather than retrying.

## Typical flow

```sh
# 1. Authenticate (once; auto-extracts from Chrome on macOS)
twt auth login
twt auth status                 # confirm handle

# 2. Read (safe anytime)
twt timeline --latest --limit 10
twt user steipete
twt tweets steipete --limit 20
twt search "openclaw" --json | jq '.[0]'
twt mentions --limit 10

# 3. Write (only on explicit user request)
twt post "Hello from the CLI!"                       # prompts to confirm
twt post "Reply text" --reply 2065650561484267540 --yes
twt like 2065650561484267540
twt follow steipete
```
