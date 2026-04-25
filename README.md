# MCPGuard

Security scanner for MCP (Model Context Protocol) servers. Detects exposed
secrets, dangerous auto-approve patterns, suspicious URLs, and shell-injection
risks across Claude Desktop, Cursor, and other MCP clients — and tracks changes
over time so rug-pulls don't go unnoticed.

```bash
npx mcpguard scan
```

## Install

```bash
# Homebrew (macOS / Linux)
brew install withaxiom/tap/mcpguard

# Go
go install github.com/withaxiom/mcpguard@latest

# Or grab a binary from the releases page
# https://github.com/withaxiom/mcpguard/releases
```

## Quick start

```bash
# Auto-detect every MCP config on this machine and scan it
mcpguard scan

# Show what sources are supported and which ones exist on disk
mcpguard scan --list-sources

# Scan a specific client
mcpguard scan --source cursor
mcpguard scan --source claude-desktop

# Scan an explicit file
mcpguard scan -c ~/some/mcp.json

# Save a snapshot, then later see what changed
mcpguard scan --save
mcpguard diff
```

## Commands

### `scan`

Detects MCP configs, lists their servers, computes file + config fingerprints,
and runs every security rule.

| Flag | Description |
|---|---|
| `-s, --source` | Restrict to one source (`cursor`, `claude-desktop`) |
| `-c, --config` | Scan an explicit config file path |
| `--save` | Save a snapshot to `~/.mcpguard/<source>-latest.json` for later diffing |
| `--list-sources` | Print supported sources, default paths, and whether each exists |
| `--json` | Machine-readable output |
| `-v, --verbose` | Show redacted env vars and rule details |

### `diff`

Compares the live config against the last saved snapshot, or two snapshots
explicitly. Use this to catch silent changes to tools you've already approved.

```bash
mcpguard diff                                       # live vs. latest snapshot
mcpguard diff --before old.json --after new.json    # two snapshots
```

## Security rules

| Rule | Severity | What it catches |
|---|---|---|
| `exposed-secrets` | warning | Env var names that look like API keys, tokens, passwords, or credentials |
| `wildcard-autoapprove` | critical | `*` / `**` patterns that grant unrestricted tool access |
| `shell-injection` | critical / warning | Shell metacharacters (`&&`, `\|\|`, backtick, `$(`) in `command` or non-URL args |
| `suspicious-url` | info / warning | Remote endpoints (non-localhost) or unencrypted `http://` URLs |
| `disabled-servers` | info | Disabled servers still hanging around in config |

URL-shaped args (`postgresql://`, `https://`, `mongodb+srv://`, …) are skipped
by the shell-injection rule so connection strings with `&` query separators
don't false-positive.

## JSON output

Every command supports `--json` so you can pipe results into other tooling:

```bash
mcpguard scan --json | jq '.configs[].warnings[] | select(.severity == "critical")'
mcpguard scan --list-sources --json
mcpguard diff --json
```

## Project layout

This repo contains both the CLI and the marketing site:

- `cmd/`, `internal/`, `main.go` — the Go CLI
- `src/`, `public/`, `next.config.ts`, `package.json` — the Next.js landing page
  at [mcpguard.com](https://mcpguard.com)
- `supabase/` — schema for the waitlist table
- `homebrew-tap/` — Homebrew formula

### CLI development

```bash
go build ./...
go test ./...
go run . scan
```

### Landing page development

```bash
npm install
cp .env.example .env.local   # fill in Supabase / Resend / PostHog keys
npm run dev
```

Required env vars for the landing page:

| Variable | Description |
|---|---|
| `RESEND_API_KEY` | Resend API key for welcome emails |
| `NEXT_PUBLIC_SUPABASE_URL` | Supabase project URL |
| `SUPABASE_SERVICE_ROLE_KEY` | Supabase service role key |
| `NEXT_PUBLIC_POSTHOG_KEY` | PostHog project API key |
| `NEXT_PUBLIC_POSTHOG_HOST` | PostHog ingest host |
| `NEXT_PUBLIC_SITE_URL` | Production site URL |

Supabase schema:

```sql
CREATE TABLE waitlist (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  source TEXT DEFAULT 'landing_page',
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_waitlist_email ON waitlist(email);
```

## License

MIT — see [LICENSE](LICENSE).

Built by [AXIOM Collective](https://withaxiom.co).
