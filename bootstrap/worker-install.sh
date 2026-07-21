#!/usr/bin/env bash
# bootstrap/worker-install.sh — one-time bootstrap for a Linux Jarvis worker.
#
# Target: AliOS 7.2 (RHEL 7 lineage, glibc >= 2.17). Other RHEL-family should work.
#
# Push-only credential distribution: worker never ssh's out to scheduler and has
# no browser for interactive SSO. The operator first runs
# `bootstrap/worker-credentials-package.sh` on the scheduler host — that packages
# ~/.config/a1, ~/.claude, and repo env files into an OSS object encrypted with a
# one-shot passphrase. Each worker then curls the object, sha256-verifies it,
# decrypts, and extracts. Deleting the OSS object revokes access.
#
# Required env vars (fail fast if unset):
#   CLAUDE_OSS_URL      OSS URL of the claude linux-x64 binary
#   CLAUDE_SHA256       expected sha256 (from Anthropic manifest)
#   CREDS_OSS_URL       OSS URL of the encrypted credential bundle
#   CREDS_SHA256        expected sha256 of the ciphertext
#   CREDS_PASSPHRASE    32-hex passphrase printed by worker-credentials-package.sh
#
# Optional env vars (defaults shown):
#   JARVIS_REPO_URL     git@gitlab.alibaba-inc.com:terraflow/jarvis-preview.git
#   JARVIS_ROOT         $HOME/workspace/jarvis-preview
#   CLAUDE_BIN          $HOME/.local/bin/claude
#   JARVIS_DISPATCH_MAX 3        (per-host slot count)
#
# Usage — see worker-credentials-package.sh output for the full one-liner. Shape:
#
#   CLAUDE_OSS_URL=... CLAUDE_SHA256=... \
#   CREDS_OSS_URL=... CREDS_SHA256=... CREDS_PASSPHRASE=... \
#   JARVIS_DISPATCH_MAX=3 \
#     bash bootstrap/worker-install.sh
#
# Exit codes:
#   0  = success, worker ready to start
#   1  = fatal (sanity/download/verify/decrypt failed)
#   2  = a1 login not portable across hosts — see step 6 output for options

set -euo pipefail

# Export BEFORE step 8 preflight runs so any role-aware tooling in the
# preflight chain sees worker mode; step 9 then persists it to bridge/jarvis.env
# for the runtime. Note verify.sh runs the SAME full tool+cred check set for
# every role — workers need gh/aliyun creds for the same triage/PR work the
# scheduler host does (gh + aliyun install sudo-free from tarballs, deps.lock).
export JARVIS_BRIDGE_ROLE="${JARVIS_BRIDGE_ROLE:-worker}"

JARVIS_REPO_URL="${JARVIS_REPO_URL:-https://code.alibaba-inc.com/terraflow/jarvis-preview.git}"
JARVIS_ROOT="${JARVIS_ROOT:-$HOME/workspace/jarvis-preview}"
CLAUDE_BIN="${CLAUDE_BIN:-$HOME/.local/bin/claude}"
DISPATCH_MAX="${JARVIS_DISPATCH_MAX:-3}"

die()  { printf '[ERR] %s\n' "$*" >&2; exit "${2:-1}"; }
info() { printf '[..]  %s\n' "$*"; }
ok()   { printf '[OK]  %s\n' "$*"; }
warn() { printf '[WARN] %s\n' "$*" >&2; }
step() { printf '\n===== %s =====\n' "$*"; }

: "${CLAUDE_OSS_URL:?export CLAUDE_OSS_URL=<OSS URL of claude linux-x64 binary>}"
: "${CLAUDE_SHA256:?export CLAUDE_SHA256=<expected sha256 of claude binary>}"
: "${CREDS_OSS_URL:?export CREDS_OSS_URL=<OSS URL of encrypted credential bundle>}"
: "${CREDS_SHA256:?export CREDS_SHA256=<sha256 of encrypted bundle>}"
: "${CREDS_PASSPHRASE:?export CREDS_PASSPHRASE=<passphrase from worker-credentials-package.sh>}"

# ---------------------------------------------------------------------------
step "1. Sanity · OS / arch / glibc"
arch=$(uname -m)
[ "$arch" = x86_64 ] || die "arch=$arch not supported (need x86_64)"
glibc=$(ldd --version 2>&1 | awk 'NR==1{print $NF}')
# glibc >= 2.17 (claude linux-x64 baseline)
lowest=$(printf '%s\n2.17\n' "$glibc" | sort -V | awk 'NR==1')
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
lowest_git=$(printf '%s\n2.5\n' "$git_ver" | sort -V | awk 'NR==1')
[ "$lowest_git" = "2.5" ] || warn "git=$git_ver < 2.5; worktree flows may misbehave"
# openssl: needed to decrypt the credential bundle in step 5. AliOS minimal
# install often omits it. Any version >= 1.0.2 works (we use -md sha256, not
# -pbkdf2 which requires 1.1.1+).
command -v openssl >/dev/null 2>&1 || sudo yum install -y openssl
ok "openssl = $(openssl version 2>&1)"
# jq: bootstrap/claim.sh + bootstrap/verify.sh parse Aone JSON via jq — without
# it every claim call and every jq-check in preflight silently fails. AliOS
# minimal install skips it.
command -v jq >/dev/null 2>&1 || sudo yum install -y jq
ok "jq = $(jq --version 2>&1)"
# unzip: cloudspec CLI ships as a zip (deps.lock) — needed at preflight install.
command -v unzip >/dev/null 2>&1 || sudo yum install -y unzip
ok "unzip = $(unzip -v 2>&1 | head -1)"

# ---------------------------------------------------------------------------
step "3. Fetch claude binary from OSS + sha256 verify"
mkdir -p "$(dirname "$CLAUDE_BIN")"
if [ -x "$CLAUDE_BIN" ] \
   && [ "$(sha256sum "$CLAUDE_BIN" | awk '{print $1}')" = "$CLAUDE_SHA256" ]; then
  ok "claude already installed with matching sha256 ($("$CLAUDE_BIN" --version 2>&1 | awk 'NR==1'))"
else
  info "downloading from $CLAUDE_OSS_URL (~265MB)"
  curl -fL --progress-bar -o "$CLAUDE_BIN" "$CLAUDE_OSS_URL"
  chmod +x "$CLAUDE_BIN"
  actual=$(sha256sum "$CLAUDE_BIN" | awk '{print $1}')
  [ "$actual" = "$CLAUDE_SHA256" ] \
    || die "sha256 mismatch: got $actual expected $CLAUDE_SHA256"
  ok "claude installed → $CLAUDE_BIN ($("$CLAUDE_BIN" --version 2>&1 | awk 'NR==1'))"
fi

# ---------------------------------------------------------------------------
step "4. Ensure jarvis repo present (git pull deferred to step 5b)"
if [ -d "$JARVIS_ROOT/.git" ]; then
  ok "repo present at $JARVIS_ROOT"
else
  # Fresh clone — needs auth. If operator pre-set GIT_HTTPS_URL (with token
  # embedded, `https://user:token@host/repo.git`), use that; otherwise fall
  # back to $JARVIS_REPO_URL (default ssh). For fully headless bootstrap the
  # operator can also tarball-relay the repo before running this script.
  mkdir -p "$(dirname "$JARVIS_ROOT")"
  clone_url="${GIT_HTTPS_URL:-$JARVIS_REPO_URL}"
  git clone "$clone_url" "$JARVIS_ROOT" \
    || die "git clone $clone_url failed; either tarball-relay the repo first or provide GIT_HTTPS_URL with token"
  ok "cloned to $JARVIS_ROOT"
fi

# ---------------------------------------------------------------------------
step "5. Fetch + decrypt credential bundle from OSS"
STAGE=$(mktemp -d -t jarvis-creds.XXXXXX)
# shellcheck disable=SC2064
trap "rm -rf '$STAGE'" EXIT
CIPHER="$STAGE/creds.tar.gz.enc"
PLAIN="$STAGE/creds.tar.gz"

info "downloading from $CREDS_OSS_URL"
curl -fL --progress-bar -o "$CIPHER" "$CREDS_OSS_URL" \
  || die "credential bundle download failed"

actual=$(sha256sum "$CIPHER" | awk '{print $1}')
[ "$actual" = "$CREDS_SHA256" ] \
  || die "creds sha256 mismatch: got $actual expected $CREDS_SHA256"
ok "bundle sha256 verified"

# Decrypt via stdin to keep the passphrase off argv.
if ! printf '%s' "$CREDS_PASSPHRASE" \
     | openssl enc -aes-256-cbc -d -md sha256 -pass stdin \
                   -in "$CIPHER" -out "$PLAIN"; then
  die "decrypt failed — CREDS_PASSPHRASE wrong, or bundle from a mismatched openssl KDF (see packager script)"
fi

# Extract to staging, then dispatch to correct roots by top-level dir.
EXTRACT="$STAGE/extract"
mkdir -p "$EXTRACT"
tar xzf "$PLAIN" -C "$EXTRACT" || die "tar extract failed"

# 5a. HOME-scoped: ~/.config/a1 + ~/.claude — cp -R with mode preservation.
if [ -d "$EXTRACT/home/.config/a1" ]; then
  mkdir -p "$HOME/.config"
  rm -rf "$HOME/.config/a1"
  cp -Rp "$EXTRACT/home/.config/a1" "$HOME/.config/a1"
fi
if [ -d "$EXTRACT/home/.claude" ]; then
  # Preserve existing worker-local files not in bundle (rare — bundle should be complete),
  # but overwrite anything from bundle.
  mkdir -p "$HOME/.claude"
  cp -Rp "$EXTRACT/home/.claude/." "$HOME/.claude/"
fi
# ~/.aliyun: aliyun CLI credentials (config.json, plaintext AK). verify.sh's
# `aliyun sts GetCallerIdentity` cred check and triage OpenAPI 查证 need it.
if [ -d "$EXTRACT/home/.aliyun" ]; then
  mkdir -p "$HOME/.aliyun"
  cp -Rp "$EXTRACT/home/.aliyun/." "$HOME/.aliyun/"
  chmod 600 "$HOME/.aliyun/config.json" 2>/dev/null || true
fi
# ~/.git-credentials + ~/.gitconfig: git's `store` credential helper needs
# both. cp -p preserves 0600 on ~/.git-credentials (mandatory — token is
# sensitive). ~/.gitconfig is 0644 (references helper name, not the token).
for f in .git-credentials .gitconfig; do
  if [ -f "$EXTRACT/home/$f" ]; then
    cp -p "$EXTRACT/home/$f" "$HOME/$f"
  fi
done
[ -f "$HOME/.git-credentials" ] && chmod 600 "$HOME/.git-credentials"

# 5b. Rewrite Mac-specific $HOME paths → this host's $HOME, best-effort.
# Only rewrites when packager's $HOME was recorded and differs from ours.
PACKAGER_HOME=$(awk -F= '/^PACKAGER_HOME=/{print $2}' "$EXTRACT/MANIFEST" 2>/dev/null || true)
if [ -n "$PACKAGER_HOME" ] && [ "$PACKAGER_HOME" != "$HOME" ]; then
  info "rewriting $PACKAGER_HOME → $HOME in extracted config/env files"
  # Use printf-safe delimiter | so slashes in paths don't collide
  find "$HOME/.claude" "$HOME/.config/a1" -type f \
       \( -name '*.json' -o -name '*.yaml' -o -name '*.yml' -o -name '*.env' -o -name '*.conf' \) \
       -exec sed -i.bak "s|$PACKAGER_HOME|$HOME|g" {} \; 2>/dev/null || true
  find "$HOME/.claude" "$HOME/.config/a1" -name '*.bak' -delete 2>/dev/null || true
fi

# 5c. Repo-scoped: overlay env + workspaces.local.json into $JARVIS_ROOT.
# if-blocks, NOT `[ -f ] && cp`: under set -e a bare `[ -f missing ] && cp`
# returns 1 and kills the whole installer when an optional file is absent
# from the bundle (workspaces.local.json is optional on the packager).
if [ -f "$EXTRACT/repo/bootstrap/.env" ]; then
  cp -p "$EXTRACT/repo/bootstrap/.env" "$JARVIS_ROOT/bootstrap/.env"
fi
if [ -f "$EXTRACT/repo/bridge/jarvis.env" ]; then
  cp -p "$EXTRACT/repo/bridge/jarvis.env" "$JARVIS_ROOT/bridge/jarvis.env"
fi
if [ -f "$EXTRACT/repo/config/workspaces.local.json" ]; then
  cp -p "$EXTRACT/repo/config/workspaces.local.json" "$JARVIS_ROOT/config/workspaces.local.json"
fi

# Rewrite Mac paths in env files too (JARVIS_ROOT / JARVIS_TATA_ROOT / CLAUDE_BIN etc)
if [ -n "$PACKAGER_HOME" ] && [ "$PACKAGER_HOME" != "$HOME" ]; then
  for f in "$JARVIS_ROOT/bootstrap/.env" "$JARVIS_ROOT/bridge/jarvis.env"; do
    [ -f "$f" ] && sed -i.bak "s|$PACKAGER_HOME|$HOME|g" "$f" && rm -f "$f.bak"
  done
fi

chmod 600 "$JARVIS_ROOT/bootstrap/.env" "$JARVIS_ROOT/bridge/jarvis.env" 2>/dev/null || true

# Zero-then-remove the plaintext tarball — decrypt output shouldn't linger.
dd if=/dev/zero of="$PLAIN" bs=1M count=1 conv=notrunc 2>/dev/null || true
rm -f "$PLAIN" "$CIPHER"
ok "creds installed; scratch dir wiped"

# Show packaging provenance for the operator's log.
if [ -f "$EXTRACT/MANIFEST" ]; then
  info "bundle provenance:"
  sed 's/^/       /' "$EXTRACT/MANIFEST"
fi

# ---------------------------------------------------------------------------
step "5b. Rewrite origin to HTTPS + (non-fatal) git pull for latest master"
# The packager records a full HTTPS URL in MANIFEST; worker uses it as origin
# so git pull goes through the credential helper (~/.git-credentials +
# ~/.gitconfig extracted in step 5a) — no ssh needed.
GIT_HTTPS_URL_MF=$(awk -F= '/^GIT_HTTPS_URL=/{print $2}' "$EXTRACT/MANIFEST" 2>/dev/null || true)
if [ -n "$GIT_HTTPS_URL_MF" ] && [ -d "$JARVIS_ROOT/.git" ]; then
  cur=$(cd "$JARVIS_ROOT" && git remote get-url origin 2>/dev/null || echo)
  if [ "$cur" != "$GIT_HTTPS_URL_MF" ]; then
    info "resetting origin: $cur → $GIT_HTTPS_URL_MF"
    (cd "$JARVIS_ROOT" && git remote set-url origin "$GIT_HTTPS_URL_MF")
  fi
fi

if [ -d "$JARVIS_ROOT/.git" ]; then
  info "attempting git pull for latest master (non-fatal)"
  if (cd "$JARVIS_ROOT" && git pull --ff-only 2>&1 | tail -3); then
    ok "repo up-to-date with origin/master"
  else
    warn "git pull failed — repo stays at whatever the tarball/clone had; check remote URL + Deploy Token"
    warn "  cur origin: $(cd "$JARVIS_ROOT" && git remote get-url origin 2>/dev/null)"
  fi
fi

# ---------------------------------------------------------------------------
step "6. a1 login portability probe (bundled tokens must survive rehost)"
cd "$JARVIS_ROOT"
if bin/a1id ready jarvis; then
  ok "a1 jarvis login OK on this host (tokens portable)"
else
  cat >&2 <<'DIAG'

[ERR] a1 jarvis auth is not usable on this host.

This is likely because a1's login token was bound to the scheduler host's
fingerprint (hostname / MAC / TPM / RAM role). Since this worker has no
browser and can't easily be ssh'd to, the standard interactive fallback
(`bin/a1id as jarvis -- auth login`) is impractical.

Options (in decreasing preference):
  A) Ask a1 team for a machine-portable or long-lived service token that
     survives cross-host copy. Feed it via ~/.config/a1/jarvis/token.json.
  B) Package tokens right AFTER a fresh SSO login on the scheduler and
     use them immediately (some tokens are portable for their initial TTL
     but not after refresh). Re-run this install within that window.
  C) Escalate: this worker can't participate in the fleet without a1
     writes, since claim.sh / wrap.sh / aone-get.sh all call a1.

For now: FAILED. Fix upstream, then re-run this script (steps 1-5 are
idempotent).

DIAG
  exit 2
fi
if bin/a1id ready terraform-rd 2>/dev/null; then
  ok "a1 terraform-rd login OK"
else
  warn "a1 terraform-rd not logged in (needed for Terraform Task types)"
  warn "  Package again after logging into terraform-rd on the scheduler, or"
  warn "  this worker will fail-closed on Terraform Tasks."
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

# enable-linger so user service survives logout. On hosts where sudo is
# intercepted by command control (sat_app whitelist), fall back to a user
# crontab: @reboot for boot autostart + a */10 watchdog. run.sh start is
# pidfile-guarded, so the watchdog is a no-op while the bridge is alive.
# (RHEL7-lineage systemd doesn't kill user processes on logout by default,
# so linger/cron only matter for boot autostart, not for ssh disconnect.)
linger_ok=0
if command -v loginctl >/dev/null 2>&1; then
  if sudo loginctl enable-linger "$USER" 2>/dev/null; then
    ok "loginctl enable-linger $USER"
    linger_ok=1
  else
    warn "enable-linger failed (sudo blocked / needs root); installing crontab fallback"
  fi
fi
if [ "$linger_ok" = 0 ]; then
  cron_cmd="JARVIS_BRIDGE_ROLE=worker $JARVIS_ROOT/bridge/run.sh start"
  cron_tmp=$(mktemp)
  # Idempotent: strip any previous lines that invoke this repo's run.sh, then
  # append the current pair. Survives re-runs and JARVIS_ROOT relocation.
  crontab -l 2>/dev/null | grep -vF "$JARVIS_ROOT/bridge/run.sh" > "$cron_tmp" || true
  {
    printf '@reboot sleep 30 && %s\n' "$cron_cmd"
    printf '*/10 * * * * %s >/dev/null 2>&1\n' "$cron_cmd"
  } >> "$cron_tmp"
  if crontab "$cron_tmp"; then
    ok "crontab fallback installed (@reboot + */10 watchdog)"
  else
    warn "crontab install failed — start the worker manually after each reboot"
  fi
  rm -f "$cron_tmp"
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
