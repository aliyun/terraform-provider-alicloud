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

## SchedulerEngine progressive cutover

The legacy periodic loops remain the default.  The new `SchedulerEngine` is a
separate, fenced control-plane path with `workerKey=bridge-scheduler`, a fresh
`processUuid`, and the same control-plane token used by the Task API. `hostId`
is recorded for observability only. The fixed Worker key, current process UUID,
status, and heartbeat ensure that only one Scheduler process owns admission at
a time. It is not a Task executor: Bridge advertises
`capabilities.dispatch.pull=false` and does not pull the public Task queue.

Use `bridge/scheduler/jobs/jobs.yaml` as the complete new-engine registration
registry. It contains only migrated jobs; each is registered and owned by the
new Engine while its legacy entry point is suppressed. Legacy jobs are absent
from this file and never enter the new control plane. The first supported
handover is `daily.probe`:

```bash
# bridge/scheduler/jobs/jobs.yaml
- key: daily.probe
  engine_runner: daily.probe
```

The YAML loader depends on `PyYAML`, installed into the isolated Bridge venv by
`bootstrap/worker-install.sh`. On a scheduler host that was not bootstrapped by
that installer, run `bash bootstrap/bridge-python.sh` once before
`bridge/run.sh start`; the run script automatically prefers `.venv/bridge`.

Before this starts the Engine, Bridge verifies the host, registers and checks
the ACTIVE Worker identity, registers the full job registry, recovers
interrupted slots, and then begins polling.  The legacy `DailyScheduler`
excludes only `daily.probe`; its existing `_ProbeJob` plus
`daily-scheduler.json` marker are reused during the handover.  A configured
new job without a runner mapping fails closed at startup. Each definition's
`enabled_env` remains an independent business switch; for example,
`JARVIS_PROBE_SCHED=0` registers a routed `daily.probe` as `DISABLED`.

A planned Scheduler restart closes admission, registers the Worker as
`DRAINING`, and waits for the admitted job to finish before starting a new
process. The Scheduler path never falls back to `SIGKILL`, and launchd restart
disables KeepAlive and avoids `kickstart -k`; a drain timeout leaves the old
process in place and refuses successor startup. An unexpected process crash is
recovered by replaying the same scheduled slot from the beginning. Generic
checkpoint/resume is deferred to job-specific implementations.

## Push-only credential distribution

Constraints on the worker side (AliOS 7.2):
- **No browser** — SSO interactive login impossible.
- **SSH inconvenient** — worker doesn't reach out to scheduler; scheduler
  doesn't push over ssh either.
- **Only OSS is reachable** — same channel that carries the claude binary.

The install therefore uses two scripts:

1. `bootstrap/worker-credentials-package.sh` — runs ONCE on the scheduler
   (macmini). Packages `~/.config/a1/`, `~/.claude/`, and repo env files;
   encrypts with a fresh AES-256 passphrase; uploads to OSS; prints the URL
   + sha256 + passphrase.

2. `bootstrap/worker-install.sh` — runs on each worker. Pulls the encrypted
   bundle from OSS, sha256-verifies, decrypts, extracts to the right paths,
   rewrites the packager's `$HOME` prefix to this host's, then continues the
   normal install (systemd unit, preflight, etc.).

Neither script requires ssh in either direction.

## Prerequisites (scheduler host, one-time-per-refresh)

Two OSS objects to prepare — one for the claude binary (long-lived), one for
the credential bundle (short-lived, revoke after fleet install).

### A. Claude binary → OSS

Anthropic geo-blocks `claude.ai/install.sh` from Alibaba internal networks.
Relay through OSS from a scheduler host in an unrestricted region (Dallas).

```bash
VERSION=$(curl -fsSL https://downloads.claude.ai/claude-code-releases/latest)
curl -fsSL -o /tmp/claude \
  "https://downloads.claude.ai/claude-code-releases/$VERSION/linux-x64/claude"
sha256sum /tmp/claude
ossutil cp /tmp/claude oss://<your-bucket>/claude-$VERSION-linux-x64
```

Give workers the OSS URL + expected sha256.

### B. Credential bundle → OSS (repeat before each worker install run)

The bundle carries: `~/.config/a1/` (a1 CLI login state), `~/.claude/*.json`
(model-lane settings archives), `bootstrap/.env` + `bridge/jarvis.env`
(GitHub PAT, control-plane token, gateway routes), plus `~/.git-credentials`
+ `~/.gitconfig` so workers can `git pull` jarvis-preview over HTTPS with a
GitLab Deploy Token — no ssh setup on the worker.

One-time setup: create a Deploy Token in GitLab web UI (scheduler-side, once):

1. Open https://code.alibaba-inc.com/terraflow/jarvis-preview
2. Settings → Repository → **Deploy Tokens** → Add token
3. Name: `jarvis-fleet-2026-07`  Scopes: **`read_repository` ONLY**
4. Save. GitLab prints a username (e.g., `gitlab+deploy-token-12345`) and a
   secret — **copy both immediately, they're shown only once**.

Package:

```bash
GIT_TOKEN_USER="gitlab+deploy-token-12345" \
GIT_TOKEN="<the secret GitLab printed>" \
OSS_BUCKET=cc-packet \
  bash bootstrap/worker-credentials-package.sh
```

Bundle now contains `~/.git-credentials`
(`https://<user>:<token>@code.alibaba-inc.com`) and `~/.gitconfig`
(`credential.helper = store`). Revoke by deleting the token in GitLab UI —
workers can no longer pull; rotate by creating a new token + re-packaging.

Output prints three values:

- `CREDS_OSS_URL`     — the OSS URL of the encrypted bundle
- `CREDS_SHA256`      — sha256 of the ciphertext (integrity check on worker)
- `CREDS_PASSPHRASE`  — 32-hex random string (256-bit) needed to decrypt

Distribute these three values to each worker via any secure channel (DingTalk
DM to yourself, in-person, ops runbook page). The passphrase alone is useless
without the OSS object; the OSS object alone is useless without the passphrase.

**After all workers install**, revoke by deleting the object:

```bash
ossutil rm oss://<bucket>/jarvis-worker-creds-<ts>.tar.gz.enc
```

Re-package before adding a NEW worker later (fresh passphrase every time).
To rotate the Deploy Token: revoke the old one in GitLab UI → create a
new one → repackage with the new `GIT_TOKEN_USER`/`GIT_TOKEN` → re-run
`worker-install.sh` on each worker (idempotent — updates the credential
file in place).

### C. Fleet knowledge

Record which hosts run scheduler vs worker in a team runbook or deployment inventory,
so an operator knows never to accidentally start `JARVIS_BRIDGE_ROLE=scheduler`
on a worker host (would create dual Task producers).

## Install a worker (per-host, one-time)

Target: AliOS 7.2 (RHEL 7 lineage, glibc ≥ 2.17). Other RHEL-family should
work; Debian/Ubuntu need the yum lines swapped.

### One-shot from bare metal (recommended)

`bootstrap/worker-bootstrap.sh` is self-contained (no repo needed on the
host). A copy is published to OSS; paste ONE command on the new worker:

```bash
GIT_TOKEN_USER=open_jarvis GIT_TOKEN=<code.alibaba-inc.com token> \
CLAUDE_OSS_URL="https://cc-packet.oss-cn-beijing.aliyuncs.com/claude" \
CLAUDE_SHA256="c1efffaaf370aa187cb6a09dd93d4e511c646899b0078476f83791b664bde7fe" \
CREDS_OSS_URL="<bundle URL from worker-credentials-package.sh>" \
CREDS_SHA256="<sha printed by packager>" \
CREDS_PASSPHRASE="<32-hex printed by packager>" \
JARVIS_DISPATCH_MAX=3 \
  bash -c "$(curl -fsSL https://cc-packet.oss-cn-beijing.aliyuncs.com/worker-bootstrap.sh)"
```

It installs git/curl if missing, clones the repo via HTTPS+token (then scrubs
the token from `.git/config` — later pulls use the credential store the
bundle installs), and hands off to `worker-install.sh` below. The OSS copy is
a mirror; after changing the script in-repo, re-publish:
`ossutil cp bootstrap/worker-bootstrap.sh oss://cc-packet/worker-bootstrap.sh -f`.

### If the repo is already on the host

```bash
cd ~/workspace/jarvis-preview

# All five env vars are required (feed from Prereq A + B outputs):
CLAUDE_OSS_URL="https://cc-packet.oss-cn-beijing-internal.aliyuncs.com/claude" \
CLAUDE_SHA256="c1efffaaf370aa187cb6a09dd93d4e511c646899b0078476f83791b664bde7fe" \
CREDS_OSS_URL="https://cc-packet.oss-cn-beijing-internal.aliyuncs.com/jarvis-worker-creds-<ts>.tar.gz.enc" \
CREDS_SHA256="<sha printed by packager>" \
CREDS_PASSPHRASE="<32-hex printed by packager>" \
JARVIS_DISPATCH_MAX=3 \
  bash bootstrap/worker-install.sh
```

The script runs these steps in order:
1. Sanity: OS/arch/glibc (fails closed on unsupported host)
2. Install python3 (≥ 3.8) and git if missing — yum first, then the
   canary-proven fallback: uv + python-build-standalone 3.12 (glibc 2.17
   baseline) with `/usr/bin/python3` symlinked to it (cron/systemd PATH
   lacks `~/.local/bin`)
3. Fetch claude binary from OSS, sha256 verify
4. Ensure jarvis repo present (`git clone` only if missing; `git pull`
   deferred to step 5b)
5. Pull encrypted credential bundle from OSS, sha256 verify, decrypt (openssl
   aes-256-cbc + pbkdf2 100k iters, passphrase read from stdin — never on
   argv), extract classified into `$HOME` (a1/claude/git-credentials) and
   `$JARVIS_ROOT` (env files), then rewrite the packager's `$HOME` prefix →
   this host's `$HOME` in all extracted json/yaml/env/conf files. Scratch
   dir zeroed + wiped on exit.
5b. Rewrite `origin` remote to HTTPS URL from MANIFEST + attempt `git pull`
    (non-fatal) via the newly installed credential helper + Deploy Token.
6. Probe `a1id ready jarvis` — see [a1 portability fallback](#a1-portability)
7. Write `config/workspaces.local.json` for this host's paths (if the bundle
   didn't ship one)
7.5. Build cloudspec from source if missing (Linux has no upstream artifacts;
   `bootstrap/cloudspec-build-linux.sh`, non-fatal — preflight WARNs on skip;
   `JARVIS_SKIP_CLOUDSPEC_BUILD=1` to opt out)
8. `bootstrap/preflight.sh --force` (fail closed)
9. Append `JARVIS_BRIDGE_ROLE=worker` + `JARVIS_DISPATCH_MAX=N` to
   `bridge/jarvis.env`; if DingTalk keys are present (inherited from
   scheduler), also set `JARVIS_NO_DINGTALK=1` so the worker skips stream
   client startup (run.sh + the bot enforce this by role regardless)
10. Install `~/.config/systemd/user/jarvis-worker.service` +
    `loginctl enable-linger`; when sudo/linger is blocked by command control,
    fall back to a user crontab (`@reboot` autostart + `*/10` watchdog via
    `run.sh start`, installing cronie + starting crond as needed). Without
    linger, prefer `run.sh start` over `systemctl --user` — the systemd user
    unit dies with your last ssh session, run.sh survives it

No openssh-clients is installed on the worker — HTTPS + Deploy Token is
the only supported git auth. If you need ssh for other reasons on the
worker, install it separately.

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
# 1. On a host with claude.ai reachable (Dallas macmini):
VERSION=$(curl -fsSL https://downloads.claude.ai/claude-code-releases/latest)
curl -fsSL -o /tmp/claude "https://downloads.claude.ai/claude-code-releases/$VERSION/linux-x64/claude"
SHA=$(sha256sum /tmp/claude | awk '{print $1}')
ossutil cp /tmp/claude oss://<your-bucket>/claude-$VERSION-linux-x64

# 2. On each worker (creds can be re-used if valid; if a1 token TTL has
#    lapsed, re-package creds first with worker-credentials-package.sh):
cd ~/workspace/jarvis-preview
git pull --ff-only
CLAUDE_OSS_URL="https://<bucket>.oss-cn-beijing-internal.aliyuncs.com/claude-$VERSION-linux-x64" \
CLAUDE_SHA256=$SHA \
CREDS_OSS_URL=<same as install> \
CREDS_SHA256=<same> \
CREDS_PASSPHRASE=<same> \
  bash bootstrap/worker-install.sh
# Only step 3 does real work (new binary); creds re-download is idempotent.
systemctl --user restart jarvis-worker
```

## Refresh credentials

Do this when:
- a1 token approaching expiry (before workers start failing bookend)
- Anthropic gateway token rotated
- GitHub PAT rotated
- Adding a new worker after original bundle was revoked

```bash
# On scheduler (fresh a1 login first if TTL is short):
bin/a1id as jarvis -- auth login              # optional; if fingerprint-portable
OSS_BUCKET=cc-packet bash bootstrap/worker-credentials-package.sh
# → prints new CREDS_OSS_URL / CREDS_SHA256 / CREDS_PASSPHRASE

# On each existing worker:
CREDS_OSS_URL=<new> CREDS_SHA256=<new> CREDS_PASSPHRASE=<new> \
CLAUDE_OSS_URL=<current> CLAUDE_SHA256=<current> \
  bash bootstrap/worker-install.sh
systemctl --user restart jarvis-worker

# After all workers updated:
ossutil rm oss://<bucket>/jarvis-worker-creds-<old-ts>.tar.gz.enc
```

## Update jarvis code

Same as before — merge to master on scheduler host. Each worker's
`worker-install.sh` re-run picks up the pull; or from cron/manually:

```bash
cd ~/workspace/jarvis-preview
git pull --ff-only
systemctl --user restart jarvis-worker
```

## Cross-host interrupted-session detection

Worker/Task/Session state in the control plane is the only coordination truth;
there is no local PID-based coord fallback. Cross-host recovery goes through
two channels already implemented in the control plane (`docs/execution-architecture.md`):

1. **Fast channel**: bridge scheduler's tick watches `/workers` for
   STALE/OFFLINE Workers and remembers their assignments.
2. **Persistent channel**: `AoneScheduler` atomically snapshots the
   `jarvis-claimed` Aone inventory to `.my-day/bridge/claimed-snapshot.json`
   every tick and corroborates each entry with the control plane's Task
   timeline. Snapshot lives only on scheduler host — this is why exactly one
   scheduler runs.

Workers themselves don't do orphan detection; they just heartbeat and lease.

## Troubleshooting

### a1 portability

`~/.config/a1/` may not survive a host copy if the internal SSO system binds
tokens to a machine fingerprint (hostname / MAC / TPM / RAM role). Workers
have no browser and can't easily be ssh'd to, so the standard interactive
fallback (`bin/a1id as jarvis -- auth login`) is impractical.

If step 6 fails, options in decreasing preference:

- **A. Ask the a1 team** for a machine-portable or long-lived service token
  that survives cross-host copy. Drop it in the packager's
  `~/.config/a1/jarvis/token.json` before re-running the packager.
- **B. Fresh-token window** — some SSO systems mint tokens that are portable
  for their initial TTL but not after refresh. Package immediately after a
  fresh `bin/a1id as jarvis -- auth login` on the scheduler, then install
  all workers within that window (typically hours).
- **C. Escalate** — without working a1 on the worker, `claim.sh`, `wrap.sh`,
  and `aone-get.sh` fail. The worker can't participate in the fleet.

`worker-install.sh` exits with code 2 on step 6 failure. Fix upstream, then
re-run (steps 1-5 are idempotent).

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
