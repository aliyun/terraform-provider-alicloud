# Multi-worker deployment

Phase 1 of Aone #84478029: horizontally scale headless Task execution across
multiple hosts by splitting the bridge into two roles.

## Architecture

```
                  ┌────────────────────────────────────────────────┐
                  │  Control plane (remote HTTP)                    │
                  │  https://pre-agent.aliyun-inc.com               │
                  │  Task / Session / Worker / fence / operations   │
                  └────────────────────────────────────────────────┘
                       ▲          ▲          ▲          ▲
                       │lease     │lease     │lease     │lease
                       │+heartbeat│          │          │
              ┌────────┴────┐  ┌──┴────┐  ┌──┴────┐  ┌──┴────┐
              │  scheduler   │  │worker1│  │worker2│  │workerN│
              │  (macmini)   │  │(linux)│  │(linux)│  │(linux)│
              │              │  │       │  │       │  │       │
              │ • PersistExec│  │•Persist│  │•Persist│  │•Persist│
              │ • AoneSched  │  │ Exec  │  │ Exec  │  │ Exec  │
              │ • DailySched │  │       │  │       │  │       │
              │ • PrWatch    │  │       │  │       │  │       │
              │ • AoneReply  │  │       │  │       │  │       │
              └──────────────┘  └───────┘  └───────┘  └───────┘
                 5 slots         3 slots    3 slots    3 slots
                                     total = 5 + N×3 slots
```

**Scheduler role** — the Task producer + all periodic sensor loops. Local state
(`.my-day/bridge/pr-watch.json`, `daily-scheduler.json`, event ledgers,
`claimed-snapshot.json`) is per-host and cannot be shared. **Exactly one host
in the fleet runs scheduler.** Currently: the macmini in Dallas.

**Worker role** — executor-only. `PersistenceExecutor.start()` then blocks on
`run_forever()`. Leases Tasks from control plane by its own `workerKey`,
spawns headless claude, reports back. Every added worker grows total lease
capacity by its local `JARVIS_DISPATCH_MAX`. No coordination needed between
workers — fenced `claimTask` on the control plane is the only interlock.

## Prerequisites (scheduler host, one-time)

1. **Claude binary on OSS** — Anthropic geo-blocks `claude.ai/install.sh` from
   Alibaba internal networks. Relay through OSS:

   ```bash
   # On scheduler host (with claude.ai reachable):
   VERSION=$(curl -fsSL https://downloads.claude.ai/claude-code-releases/latest)
   curl -fsSL -o /tmp/claude \
     "https://downloads.claude.ai/claude-code-releases/$VERSION/linux-x64/claude"
   sha256sum /tmp/claude  # note the hash
   ossutil cp /tmp/claude oss://<your-bucket>/claude-$VERSION-linux-x64
   ```

   Give workers the OSS URL + expected sha256.

2. **SSH from worker → scheduler** — `worker-install.sh` rsyncs credentials
   from the scheduler. Set up public-key SSH from each worker's `admin` user
   to the scheduler host's shell user.

3. **Fleet knowledge** — record which hosts run scheduler vs worker in
   `escalation/` or a team runbook, so an operator knows never to accidentally
   start `JARVIS_BRIDGE_ROLE=scheduler` on a worker host (would create dual
   Task producers).

## Install a worker (per-host, one-time)

Target: AliOS 7.2 (RHEL 7 lineage, glibc ≥ 2.17). Other RHEL-family should
work; Debian/Ubuntu need the yum lines swapped.

```bash
# On the worker host, as the ops user (e.g., admin):
git clone git@gitlab.alibaba-inc.com:terraflow/jarvis-preview.git ~/workspace/jarvis-preview
cd ~/workspace/jarvis-preview

# All env vars in one line. CLAUDE_OSS_URL + CLAUDE_SHA256 are required.
CLAUDE_OSS_URL="https://cc-packet.oss-cn-beijing-internal.aliyuncs.com/claude" \
CLAUDE_SHA256="c1efffaaf370aa187cb6a09dd93d4e511c646899b0078476f83791b664bde7fe" \
CREDS_SOURCE="gzzz@macmini-dallas" \
JARVIS_DISPATCH_MAX=3 \
  bash bootstrap/worker-install.sh
```

The script runs 10 steps in order:
1. Sanity: OS/arch/glibc (fails closed on unsupported host)
2. Install python3 (≥ 3.8) and git if missing
3. Fetch claude binary from OSS, verify sha256
4. Clone or update the jarvis repo
5. rsync credentials from scheduler (`~/.config/a1/`, `~/.claude/`,
   `bootstrap/.env`, `bridge/jarvis.env`)
6. Probe `a1id ready jarvis` — see [a1 fingerprint fallback](#a1-fingerprint-fallback)
7. Write `config/workspaces.local.json` for this host's paths
8. `bootstrap/preflight.sh --force` (fail closed)
9. Append `JARVIS_BRIDGE_ROLE=worker` + `JARVIS_DISPATCH_MAX=N` to
   `bridge/jarvis.env`; if DingTalk keys are present (inherited from
   scheduler), also set `JARVIS_NO_DINGTALK=1` so the worker skips stream
   client startup
10. Install `~/.config/systemd/user/jarvis-worker.service` +
    `loginctl enable-linger` so the worker survives operator logout

## Start / stop / observe

```bash
# Start (systemd user service, auto-restart on failure):
systemctl --user start jarvis-worker

# Enable auto-start on boot:
systemctl --user enable jarvis-worker

# Watch:
journalctl --user -u jarvis-worker -f
tail -f ~/workspace/jarvis-preview/.my-day/bridge/bot-worker.log

# Stop (graceful; drains in-flight Sessions as RESUMABLE without consuming Task
# retry budget, per the "preserve task sessions across bridge restarts" change):
systemctl --user stop jarvis-worker

# Or manually, bypassing systemd:
JARVIS_BRIDGE_ROLE=worker ~/workspace/jarvis-preview/bridge/run.sh start
JARVIS_BRIDGE_ROLE=worker ~/workspace/jarvis-preview/bridge/run.sh stop
JARVIS_BRIDGE_ROLE=worker ~/workspace/jarvis-preview/bridge/run.sh status
```

Worker pidfile: `.my-day/bridge/bot-worker.pid`  (scheduler uses `bot.pid`)
Worker log:     `.my-day/bridge/bot-worker.log`  (scheduler uses `bot.log`)

## Fleet visibility

From any jarvis checkout, ask the control plane which workers are alive:

```bash
bootstrap/control-plane-status.sh workers
```

Each worker registers with a stable `workerKey` (persisted in
`.my-day/bridge/worker-id` on the worker host, fcntl-locked to prevent
duplicates). A worker's `processUuid` changes on every restart but its
`workerKey` stays. A missed heartbeat past the lease safety margin transitions
the Worker to STALE/OFFLINE and its Task is subject to control-plane recovery
per the RESUME_ONLY / REPLAY_SAFE / MANUAL policy.

## Scaling

Add a worker: run `worker-install.sh` on a new host. It shows up in `/workers`
within one heartbeat interval (30s). No scheduler restart needed.

Remove a worker: `systemctl --user stop jarvis-worker`. Graceful stop
transitions in-flight Sessions to RESUMABLE; another worker (or the
scheduler's own executor) leases them from the queue on the next tick.

## Update claude binary

When Anthropic ships a new version:

```bash
# 1. On a host with claude.ai reachable (Dallas):
VERSION=$(curl -fsSL https://downloads.claude.ai/claude-code-releases/latest)
curl -fsSL -o /tmp/claude "https://downloads.claude.ai/claude-code-releases/$VERSION/linux-x64/claude"
SHA=$(sha256sum /tmp/claude | awk '{print $1}')
ossutil cp /tmp/claude oss://<your-bucket>/claude-$VERSION-linux-x64
echo "URL: https://<bucket>.oss-cn-beijing-internal.aliyuncs.com/claude-$VERSION-linux-x64"
echo "SHA: $SHA"

# 2. On each worker:
CLAUDE_OSS_URL="https://.../claude-$VERSION-linux-x64" \
CLAUDE_SHA256=$SHA \
  bash bootstrap/worker-install.sh
# Only step 3 does real work; other steps are idempotent no-ops.
systemctl --user restart jarvis-worker
```

## Update jarvis code

Same as before — merge to master on scheduler host. Each worker's
`worker-install.sh` re-run picks up the pull; or from cron/manually:

```bash
cd ~/workspace/jarvis-preview
git pull --ff-only
systemctl --user restart jarvis-worker
```

## Cross-host orphan detection

`bootstrap/coord.sh` uses `kill -0 <pid>` to test liveness, which is
same-host only. It is **not** the multi-host orphan-detection mechanism.
Instead, cross-host recovery goes through two channels already implemented
in the control plane (`docs/execution-architecture.md`):

1. **Fast channel**: bridge scheduler's tick watches `/workers` for
   STALE/OFFLINE Workers and remembers their assignments.
2. **Persistent channel**: `AoneScheduler` atomically snapshots the
   `jarvis-claimed` Aone inventory to `.my-day/bridge/claimed-snapshot.json`
   every tick and corroborates each entry with the control plane's Task
   timeline. Snapshot lives only on scheduler host — this is why exactly one
   scheduler runs.

Workers themselves don't do orphan detection; they just heartbeat and lease.

## Troubleshooting

### a1 fingerprint fallback

`~/.config/a1/` may not survive a host copy if the internal SSO system binds
tokens to a machine fingerprint. `worker-install.sh` detects this at step 6
(`bin/a1id ready jarvis` returns non-zero) and exits with code 3, printing
the interactive re-login command:

```bash
bin/a1id as jarvis -- auth login          # opens browser or provides SSO code
```

Complete the SSO flow, then re-run `worker-install.sh` (steps 1-5 are
idempotent). If the worker will handle Terraform Tasks, also:

```bash
bin/a1id as terraform-rd -- auth login
```

### glibc too old

If `worker-install.sh` step 1 fails with `glibc=X.Y < 2.17`, this host cannot
run the linux-x64 claude binary directly. Fallbacks (see Aone #84478029
comment for full analysis):
- **Path A**: run claude inside a podman container based on
  `docker.io/library/debian:12-slim` — needs `_headless_exec_command` wrapper
  work (Phase 2 or ad hoc).
- **Path B**: ask ops to upgrade the host OS to AliOS 3 (glibc ≥ 2.28) or
  Ubuntu 22.04. Preferred.
- **Path C**: musl variant — requires `/lib/ld-musl-x86_64.so.1`, not
  installed on RHEL-family. Not recommended.

### Control plane unreachable

`worker-install.sh` step 8 (`preflight`) probes the control plane URL. If it
fails:
- Confirm the worker host can reach `${JARVIS_CONTROL_PLANE_BASE_URL}` (from
  `bootstrap/.env` or `bridge/jarvis.env`)
- Confirm `JARVIS_CONTROL_PLANE_TOKEN` is populated and current
- Check `curl -I https://pre-agent.aliyun-inc.com/` returns 200/301

### Claude binary geo-blocked

`claude.ai` and `downloads.claude.ai` are blocked from Alibaba internal
networks — the installer will download an HTML "app unavailable in region"
page instead of the script. **Always** relay through OSS from a scheduler
host in an unrestricted region (Dallas macmini). `worker-install.sh` uses
your OSS URL, never `claude.ai/install.sh` directly.

### Systemd service won't start

Check `journalctl --user -u jarvis-worker -n 50`. Common issues:
- `bridge/run.sh` errors on missing env → confirm `bridge/jarvis.env` and
  `bootstrap/.env` are readable by the service user.
- `python3` not on PATH → systemd's PATH is minimal; the service exec is
  `bash bridge/run.sh daemon` which sources env and normalizes PATH.
- `loginctl enable-linger` not applied → service dies on user logout;
  re-run with `sudo` if step 10 warned about it.

## What Phase 1 does not do

- **Wake path cross-host**: SUSPENDED session wake still relies on the
  local `~/.claude/projects/<slug>/<sid>.jsonl` transcript, which is
  per-host. Phase 1 doesn't cross-host wake; the recovery scheduler picks a
  worker with workerKey affinity or falls back to REPLAY_SAFE (re-prompt
  from Aone context). Phase 2 addresses this.
- **Leader election**: if the scheduler host dies, no new Tasks are
  produced until it comes back up. In-flight Tasks continue on workers.
  Phase 2 or later may add hot standby.
- **Load balancing across workers**: control plane leases first-come; if
  one worker is faster it grabs more. That's usually fine; not tuned.

## Rollback

Stop workers first (they'll drain gracefully), then continue running the
scheduler as before:

```bash
# On each worker:
systemctl --user stop jarvis-worker
systemctl --user disable jarvis-worker

# Scheduler host: no change needed — it never depended on workers being up.
```
