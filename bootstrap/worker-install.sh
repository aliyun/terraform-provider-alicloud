#!/usr/bin/env bash
# bootstrap/worker-install.sh — one-time bootstrap for a Linux Jarvis worker.
#
# Target: AliOS 7.2 (RHEL 7 lineage, glibc >= 2.17). Other RHEL-family should work.
#
# Prereq on the operator side (macmini or any scheduler host with working creds):
#   1) Upload the Linux x86_64 claude binary to an OSS bucket reachable from AliOS
#      (internal endpoint preferred for speed). Note its sha256.
#   2) Ensure SSH from this AliOS worker → operator host is set up (public-key auth)
#      OR be ready to manually copy creds if CREDS_SOURCE is left empty.
#
# Required env vars (fail fast if unset):
#   CLAUDE_OSS_URL   e.g. https://cc-packet.oss-cn-beijing-internal.aliyuncs.com/claude
#   CLAUDE_SHA256    expected sha256 of the binary (from Anthropic manifest)
#
# Optional env vars (defaults shown):
#   JARVIS_REPO_URL     git@gitlab.alibaba-inc.com:terraflow/jarvis-preview.git
#   JARVIS_ROOT         $HOME/workspace/jarvis-preview
#   CLAUDE_BIN          $HOME/.local/bin/claude
#   CREDS_SOURCE        (empty; else "user@host" for rsync-pull)
#   JARVIS_DISPATCH_MAX 3        (per-host slot count)
#
# Usage:
#   CLAUDE_OSS_URL=... CLAUDE_SHA256=... CREDS_SOURCE=gzzz@macmini.dallas.local \
#     bash bootstrap/worker-install.sh
#
# Exit codes:
#   0  = success, worker ready to start
#   1  = fatal (sanity/download/verify failed)
#   2  = creds missing, manually populate then re-run
#   3  = a1 login fingerprint bound to source; interactive re-login required

set -euo pipefail

JARVIS_REPO_URL="${JARVIS_REPO_URL:-git@gitlab.alibaba-inc.com:terraflow/jarvis-preview.git}"
JARVIS_ROOT="${JARVIS_ROOT:-$HOME/workspace/jarvis-preview}"
CLAUDE_BIN="${CLAUDE_BIN:-$HOME/.local/bin/claude}"
CREDS_SOURCE="${CREDS_SOURCE:-}"
DISPATCH_MAX="${JARVIS_DISPATCH_MAX:-3}"

die()  { printf '[ERR] %s\n' "$*" >&2; exit "${2:-1}"; }
info() { printf '[..]  %s\n' "$*"; }
ok()   { printf '[OK]  %s\n' "$*"; }
warn() { printf '[WARN] %s\n' "$*" >&2; }
step() { printf '\n===== %s =====\n' "$*"; }

: "${CLAUDE_OSS_URL:?export CLAUDE_OSS_URL=<OSS URL of claude linux-x64 binary>}"
: "${CLAUDE_SHA256:?export CLAUDE_SHA256=<expected sha256 of claude binary>}"

# ---------------------------------------------------------------------------
step "1. Sanity · OS / arch / glibc"
arch=$(uname -m)
[ "$arch" = x86_64 ] || die "arch=$arch not supported (need x86_64)"
glibc=$(ldd --version 2>&1 | head -1 | awk '{print $NF}')
# glibc >= 2.17 (claude linux-x64 baseline)
lowest=$(printf '%s\n2.17\n' "$glibc" | sort -V | head -1)
[ "$lowest" = "2.17" ] || die "glibc=$glibc < 2.17; claude binary won't load"
os_name=$(grep -E '^PRETTY_NAME=' /etc/os-release 2>/dev/null | cut -d= -f2 | tr -d '"')
ok "OS=$os_name  arch=$arch  glibc=$glibc  kernel=$(uname -r)"

# ---------------------------------------------------------------------------
step "2. Install python3 (bridge needs 3.8+) and git if missing"
have_py3() {
  command -v python3 >/dev/null 2>&1 \
    && python3 -c 'import sys; sys.exit(0 if sys.version_info >= (3,8) else 1)' 2>/dev/null
}
if have_py3; then
  ok "python3 = $(python3 --version 2>&1)"
else
  info "python3 (>=3.8) missing; attempting yum install"
  for pkg in python38 python3 rh-python38-python; do
    if sudo yum install -y "$pkg" 2>/dev/null; then
      info "installed $pkg"
      break
    fi
  done
  have_py3 || die "could not install python3.8+; try SCL or conda manually then re-run"
  ok "python3 = $(python3 --version 2>&1)"
fi
command -v git >/dev/null 2>&1 || sudo yum install -y git
git_ver=$(git --version | awk '{print $3}')
lowest_git=$(printf '%s\n2.5\n' "$git_ver" | sort -V | head -1)
[ "$lowest_git" = "2.5" ] || warn "git=$git_ver < 2.5; worktree flows may misbehave"

# ---------------------------------------------------------------------------
step "3. Fetch claude binary from OSS + sha256 verify"
mkdir -p "$(dirname "$CLAUDE_BIN")"
if [ -x "$CLAUDE_BIN" ] \
   && [ "$(sha256sum "$CLAUDE_BIN" | awk '{print $1}')" = "$CLAUDE_SHA256" ]; then
  ok "claude already installed with matching sha256 ($("$CLAUDE_BIN" --version 2>&1 | head -1))"
else
  info "downloading from $CLAUDE_OSS_URL (~265MB)"
  curl -fL --progress-bar -o "$CLAUDE_BIN" "$CLAUDE_OSS_URL"
  chmod +x "$CLAUDE_BIN"
  actual=$(sha256sum "$CLAUDE_BIN" | awk '{print $1}')
  [ "$actual" = "$CLAUDE_SHA256" ] \
    || die "sha256 mismatch: got $actual expected $CLAUDE_SHA256"
  ok "claude installed → $CLAUDE_BIN ($("$CLAUDE_BIN" --version 2>&1 | head -1))"
fi

# ---------------------------------------------------------------------------
step "4. Clone / update jarvis repo"
if [ -d "$JARVIS_ROOT/.git" ]; then
  ok "repo present at $JARVIS_ROOT; updating"
  (cd "$JARVIS_ROOT" && git pull --ff-only 2>&1 | tail -3)
else
  mkdir -p "$(dirname "$JARVIS_ROOT")"
  git clone "$JARVIS_REPO_URL" "$JARVIS_ROOT"
  ok "cloned to $JARVIS_ROOT"
fi

# ---------------------------------------------------------------------------
step "5. Sync shared credentials from scheduler host"
if [ -z "$CREDS_SOURCE" ]; then
  warn "CREDS_SOURCE unset — cannot auto-pull creds. Populate manually:"
  cat >&2 <<MANUAL

  Required paths on this host (copy from a working scheduler):

    ~/.config/a1/                      # a1 CLI login state (may not survive rehost;
                                       # if so see Step 6 fallback)
    ~/.claude/                         # claude settings archives (idea_settings.json,
                                       # glm5.2.json, etc.)
    $JARVIS_ROOT/bootstrap/.env        # GitHub PAT, control-plane token, etc.
    $JARVIS_ROOT/bridge/jarvis.env     # DingTalk keys (optional for worker),
                                       # settings archive paths, model lanes

  After copying, re-run this script.

MANUAL
  exit 2
fi
info "rsync from $CREDS_SOURCE"
for path in ".config/a1/" ".claude/"; do
  mkdir -p "$HOME/$path"
  rsync -a --delete "$CREDS_SOURCE:~/$path" "$HOME/$path" \
    || die "rsync ~/$path failed"
done
rsync -a "$CREDS_SOURCE:$JARVIS_ROOT/bootstrap/.env" "$JARVIS_ROOT/bootstrap/.env" \
  || warn "bootstrap/.env not present on source; ensure it's populated"
rsync -a "$CREDS_SOURCE:$JARVIS_ROOT/bridge/jarvis.env" "$JARVIS_ROOT/bridge/jarvis.env" \
  || warn "bridge/jarvis.env not present on source"
chmod 600 "$JARVIS_ROOT/bootstrap/.env" "$JARVIS_ROOT/bridge/jarvis.env" 2>/dev/null || true
ok "creds synced"

# ---------------------------------------------------------------------------
step "6. a1 login probe (fingerprint may re-require login on new host)"
cd "$JARVIS_ROOT"
if bin/a1id ready jarvis; then
  ok "a1 jarvis login OK on this host"
else
  warn "a1 jarvis login FAILED — token likely fingerprint-bound to source host."
  warn "Fallback: log in interactively on THIS machine (needs browser/SSO):"
  warn "    $JARVIS_ROOT/bin/a1id as jarvis -- auth login"
  warn "Then re-run this script (steps 1-5 are idempotent)."
  exit 3
fi
if bin/a1id ready terraform-rd 2>/dev/null; then
  ok "a1 terraform-rd login OK"
else
  warn "a1 terraform-rd not logged in (needed for Terraform Task types)"
  warn "  bin/a1id as terraform-rd -- auth login  (if this worker will handle Terraform)"
fi

# ---------------------------------------------------------------------------
step "7. Write config/workspaces.local.json (this host's paths)"
local_ws="$JARVIS_ROOT/config/workspaces.local.json"
if [ ! -f "$local_ws" ]; then
  cat > "$local_ws" <<EOF
{
  "workspace_root": "$HOME/workspace"
}
EOF
  ok "created $local_ws (workspace_root=$HOME/workspace)"
else
  ok "$local_ws exists; leaving alone (verify manually if paths need adjustment)"
fi

# ---------------------------------------------------------------------------
step "8. bootstrap/preflight.sh --force"
bash bootstrap/preflight.sh --force || die "preflight failed; fix errors above and re-run"

# ---------------------------------------------------------------------------
step "9. Configure worker role in bridge/jarvis.env"
env_file="$JARVIS_ROOT/bridge/jarvis.env"
touch "$env_file"
if grep -qE '^JARVIS_BRIDGE_ROLE=' "$env_file"; then
  role=$(grep '^JARVIS_BRIDGE_ROLE=' "$env_file" | tail -1 | cut -d= -f2)
  [ "$role" = worker ] \
    || die "$env_file has JARVIS_BRIDGE_ROLE=$role; expected 'worker' (fix manually then re-run)"
  ok "JARVIS_BRIDGE_ROLE=worker already set"
else
  printf '\n# multi-worker phase 1: this host is a worker\nJARVIS_BRIDGE_ROLE=worker\n' >> "$env_file"
  ok "appended JARVIS_BRIDGE_ROLE=worker → $env_file"
fi
if ! grep -qE '^JARVIS_DISPATCH_MAX=' "$env_file"; then
  printf 'JARVIS_DISPATCH_MAX=%s\n' "$DISPATCH_MAX" >> "$env_file"
  ok "appended JARVIS_DISPATCH_MAX=$DISPATCH_MAX → $env_file"
fi
# Workers don't need DingTalk output — clear keys to force degraded mode (skip stream)
if grep -qE '^DINGTALK_APP_KEY=' "$env_file"; then
  warn "DINGTALK_APP_KEY present on worker; setting JARVIS_NO_DINGTALK=1 to skip stream client"
  grep -qE '^JARVIS_NO_DINGTALK=' "$env_file" \
    || printf 'JARVIS_NO_DINGTALK=1\n' >> "$env_file"
fi

# ---------------------------------------------------------------------------
step "10. Install systemd user unit"
mkdir -p "$HOME/.config/systemd/user"
unit="$HOME/.config/systemd/user/jarvis-worker.service"
cat > "$unit" <<EOF
[Unit]
Description=Jarvis bridge (worker role) — leases Tasks from control plane
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
WorkingDirectory=$JARVIS_ROOT
ExecStart=/bin/bash $JARVIS_ROOT/bridge/run.sh daemon
Restart=on-failure
RestartSec=10
Environment=JARVIS_BRIDGE_ROLE=worker
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=default.target
EOF
ok "wrote $unit"

# enable-linger so user service survives logout
if command -v loginctl >/dev/null 2>&1; then
  sudo loginctl enable-linger "$USER" 2>/dev/null \
    && ok "loginctl enable-linger $USER" \
    || warn "enable-linger failed (may need root); service will only run while $USER is logged in"
fi
systemctl --user daemon-reload
ok "systemctl --user daemon-reload done"

# ---------------------------------------------------------------------------
step "Done · worker install complete"
cat <<SUMMARY

Next actions:
  1) Start worker:
       systemctl --user start jarvis-worker
     Or manually (foreground-ish):
       JARVIS_BRIDGE_ROLE=worker $JARVIS_ROOT/bridge/run.sh start

  2) Enable auto-start on boot:
       systemctl --user enable jarvis-worker

  3) Watch logs:
       journalctl --user -u jarvis-worker -f
     Or the bridge's own log:
       tail -f $JARVIS_ROOT/.my-day/bridge/bot-worker.log

  4) Verify registration on control plane (run this from any jarvis checkout):
       $JARVIS_ROOT/bootstrap/control-plane-status.sh workers

  5) Stop:
       systemctl --user stop jarvis-worker
     Or:
       JARVIS_BRIDGE_ROLE=worker $JARVIS_ROOT/bridge/run.sh stop

Slot capacity on this host: $DISPATCH_MAX (from JARVIS_DISPATCH_MAX).
This host will show up as a separate Worker in the control plane's /workers
endpoint. The scheduler host continues to be the sole Task producer.
SUMMARY
