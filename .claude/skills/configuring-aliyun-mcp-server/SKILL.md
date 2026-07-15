---
name: configuring-aliyun-mcp-server
description: Use when configuring the Alibaba Cloud OpenAPI MCP Server in Claude (Claude Code / Desktop) — adding a server, running multiple accounts side by side (test vs prod / 多账号), selecting credentials by aliyun CLI profile, or troubleshooting MCP auth failures like 401 "Failed to exchange token", "consent required", or an AccessKey/profile that is not being picked up.
---

# Configuring the Alibaba Cloud OpenAPI MCP Server in Claude

## Overview

There are **two independent ways** to connect Claude to the Alibaba Cloud OpenAPI MCP Server. They authenticate differently and are **not interchangeable** — the right one depends on how the account's MCP server was provisioned.

| Approach | Command | Auth | Account selection |
|---|---|---|---|
| **A. mcp-remote (OAuth)** | `npx mcp-remote-alibaba-cloud <server-url>` | Browser OAuth authorization-code, token cached in `~/.mcp-auth` per URL | Which account you **log in as** in the browser, per server URL |
| **B. uvx proxy + profile** | `uvx alibabacloud.mcp-proxy@latest` | Local AK/STS via the Credentials SDK chain → IMS `GenerateAccessToken` | `ALIBABA_CLOUD_PROFILE` env var (aliyun CLI profile) |

**Core rule — Approach B needs a server that supports token exchange:** Approach A (OAuth) works on **every** MCP server. Approach B (profile, browser-free) only works on **newer Core MCP servers** that accept the `GenerateAccessToken` internal token exchange. Older servers accept OAuth tokens (A) but reject B's token with `401 {"error":"Failed to exchange token"}`. **This is a server-generation property, NOT a consent/permission thing** — a fresh full OAuth login (which does grant consent) does **not** make an old server accept B; verified by forcing a fresh authorization-code consent and still getting the 401. **Decision:** try B with `scripts/smoke-test.py <profile>`; if it lists tools, use B for that account; if it returns `Failed to exchange token`, that account's server is old — use A, or create a new Core server at [api.aliyun.com/mcp](https://api.aliyun.com/mcp). A mixed fleet (new accounts on B, old on A) is normal.

## Multi-account: one named server per account

For test vs prod (or any N accounts), register **one Claude MCP server per account**, named distinctly. The server name is the isolation boundary — Claude sees `aliyun-test__*` vs `aliyun-prod__*` and cannot cross them.

- Prefer **user scope** for convenience, or **project scope** to pin prod to only the deploy repo.
- Do NOT try to switch accounts inside a single server (`x_assume_account_id` multi-account MCP needs cross-account RAM trust — wrong tool for isolated test/prod).

## Approach A — mcp-remote (OAuth), the reliable default

```bash
claude mcp add aliyun-test --scope user -- npx mcp-remote-alibaba-cloud <TEST_SERVER_URL>
claude mcp add aliyun-prod --scope user -- npx mcp-remote-alibaba-cloud <PROD_SERVER_URL>
```

First connect to each triggers a browser OAuth. **Log into the matching account per URL, using a separate browser profile / incognito per account** so session cookies don't cross accounts. Token caches are keyed per URL, so accounts coexist.

Get each account's `<SERVER_URL>` from [api.aliyun.com/mcp](https://api.aliyun.com/mcp) while logged into that account, or via discovery (Approach B's proxy logs `Discovered MCP server URL:` per profile — see smoke-test.py).

## Approach B — uvx proxy + aliyun CLI profile

Prereqs: `uv`/`uvx` installed; each profile's identity has `ram:GenerateAccessToken` + `openapiexplorer:*`; **and the account's MCP server must support the internal token exchange** (see Core rule — smoke-test first; older servers return `Failed to exchange token` and can't use B).

```json
{
  "mcpServers": {
    "aliyun-test": {
      "command": "uvx", "args": ["alibabacloud.mcp-proxy@latest"],
      "env": { "ALIBABA_CLOUD_PROFILE": "<test-profile>",
               "ALIBABA_CLOUD_ACCESS_KEY_ID": "", "ALIBABA_CLOUD_ACCESS_KEY_SECRET": "" }
    }
  }
}
```

**The empty-string AK is load-bearing.** The Credentials SDK default chain is: env AK/SK → OIDC → **CLI profile (`~/.aliyun/config.json`)** → credentials.ini → ECS role. If a global `ALIBABA_CLOUD_ACCESS_KEY_ID` is exported in the shell, it is inherited by the child process and **wins over your profile.** Setting it to `""` makes the env provider raise "cannot be empty"; the chain catches that and falls through to the profile named by `ALIBABA_CLOUD_PROFILE`. Extra flags: `--site-type INTL` (default CN), `--safety-policy` / `--allow-tools` to gate a prod server, `--scope` / `--client-id` / `--server-url` overrides.

## Verify

**Approach B (profile)** — run the smoke test: it sets `ALIBABA_CLOUD_PROFILE=<profile>`, spawns the uvx proxy, and drives a real MCP `initialize` + `tools/list`:

```bash
python3 scripts/smoke-test.py <profile>   # profile only — it cannot verify a server URL
```

**Approach A (server URL)** — verify with the raw curl below, or connect it in Claude and run `/mcp`.

A raw check for a known token/URL: `curl -sS -X POST <url> -H "Authorization: Bearer <tok>" -H 'Content-Type: application/json' -H 'Accept: application/json, text/event-stream' --data '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"c","version":"0"}}}'` → HTTP 200 = good.

## Troubleshooting

| Symptom | Cause | Fix |
|---|---|---|
| `401 {"error":"Failed to exchange token"}` (Approach B), persists after a fresh OAuth login | That account's MCP server is an older generation that doesn't accept `GenerateAccessToken` tokens | Use **Approach A** (mcp-remote) for that account, or create a **new** Core server at api.aliyun.com/mcp. Do NOT keep re-consenting — it won't help. |
| `IMS GenerateAccessToken failed ... Application does not contain any required scope, consent required` | Overrode `--scope` to a value the app isn't authorized for | Leave `--scope` at default; if B still fails, it's the server-generation issue above |
| Wrong account is used despite `ALIBABA_CLOUD_PROFILE` | Global `ALIBABA_CLOUD_ACCESS_KEY_ID` in env wins the chain | Set `ALIBABA_CLOUD_ACCESS_KEY_ID=""` + `..._SECRET=""` in that server's `env` |
| Both accounts hit the same server URL | Same AccessKey shared by two profiles → discovery returns same server | Check `aliyun configure list`; distinct accounts have distinct discovered URLs |
| `401 Unauthorized` on first connect (Approach A), no OAuth prompt | Token cache stale / logged into wrong account | Re-run OAuth in incognito as the correct account; caches live in `~/.mcp-auth` |

## Key facts

- New MCP servers load only after **restarting Claude Code**; check with `/mcp`.
- Machine-specific mappings (which server id = which account) belong in memory/notes, **not** this skill.
- `~/.claude.json` holds `claude mcp add` entries; back it up before bulk edits.
