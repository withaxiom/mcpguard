# MCPGuard Landing Page

Landing page for [MCPGuard](https://mcpguard.com) — the security scanner for MCP servers.

## Stack

- **Framework:** Next.js 15 (App Router, Turbopack)
- **Styling:** Tailwind CSS v4
- **Animation:** Framer Motion
- **Email:** Resend
- **Database:** Supabase (waitlist)
- **Analytics:** PostHog
- **OG Images:** @vercel/og (Edge Runtime)
- **Deployment:** Vercel

## Setup

```bash
# Install dependencies
npm install

# Copy env file and fill in values
cp .env.example .env.local

# Run dev server
npm run dev
```

## Environment Variables

| Variable | Description |
|---|---|
| `RESEND_API_KEY` | Resend API key for welcome emails |
| `NEXT_PUBLIC_SUPABASE_URL` | Supabase project URL |
| `SUPABASE_SERVICE_ROLE_KEY` | Supabase service role key |
| `NEXT_PUBLIC_POSTHOG_KEY` | PostHog project API key |
| `NEXT_PUBLIC_POSTHOG_HOST` | PostHog ingest host |
| `NEXT_PUBLIC_SITE_URL` | Production site URL |

## Supabase Table

Create this table in your Supabase project:

```sql
CREATE TABLE waitlist (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  source TEXT DEFAULT 'landing_page',
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_waitlist_email ON waitlist(email);
```

## Deploy

```bash
vercel
```

---

Built by [AXIOM Collective](https://withaxiom.co)
# mcpguard
