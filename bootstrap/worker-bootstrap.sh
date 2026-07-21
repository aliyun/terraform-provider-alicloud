#!/usr/bin/env bash
# bootstrap/worker-bootstrap.sh — bare-metal one-shot worker bootstrap.
#
# Solves the chicken-and-egg: worker-install.sh lives IN the repo, so a fresh
# host has nothing to run. This script is fully self-contained (no repo needed)
# — publish a copy to OSS and paste ONE command on each new worker:
#
#   GIT_TOKEN_USER=open_jarvis GIT_TOKEN=<code.alibaba-inc.com token> \
#   CLAUDE_OSS_URL="https://cc-packet.oss-cn-beijing.aliyuncs.com/claude" \
#   CLAUDE_SHA256="<sha256 of claude binary>" \
#   CREDS_OSS_URL="<bundle URL from worker-credentials-package.sh>" \
#   CREDS_SHA256="<bundle sha256>" \
#   CREDS_PASSPHRASE="<bundle passphrase>" \
#   JARVIS_DISPATCH_MAX=3 \
#     bash -c "$(curl -fsSL https://cc-packet.oss-cn-beijing.aliyuncs.com/worker-bootstrap.sh)"
#
# What it does:
#   1. yum-installs git/curl if missing (sudo yum is whitelisted on AliOS)
#   2. clones jarvis-preview into $JARVIS_ROOT via HTTPS+token (idempotent —
#      an existing checkout is left alone), then STRIPS the token from
#      .git/config (later pulls use ~/.git-credentials from the cred bundle)
#   3. hands off to bootstrap/worker-install.sh (env passed through), which
#      does everything else: python/uv, claude binary, credential bundle,
#      cloudspec source build, preflight, role env, systemd + cron fallback
#
# The OSS copy is a BOOTSTRAP MIRROR of this file — the repo is the source of
# truth. Re-publish after changing it:
#   ossutil cp bootstrap/worker-bootstrap.sh oss://cc-packet/worker-bootstrap.sh -f
#
# No secrets live in this file; they all arrive via the pasted env vars.
set -euo pipefail

die()  { printf '[ERR] %s\n' "$*" >&2; exit 1; }
info() { printf '[..]  %s\n' "$*"; }
ok()   { printf '[OK]  %s\n' "$*"; }
step() { printf '\n===== %s =====\n' "$*"; }

# Fail fast on every required input BEFORE any slow work — a mispasted command
# should die in the first second, not after a clone.
: "${GIT_TOKEN_USER:?export GIT_TOKEN_USER=<code.alibaba-inc.com token user, e.g. open_jarvis>}"
: "${GIT_TOKEN:?export GIT_TOKEN=<code.alibaba-inc.com private/deploy token>}"
: "${CLAUDE_OSS_URL:?export CLAUDE_OSS_URL=<OSS URL of claude linux-x64 binary>}"
: "${CLAUDE_SHA256:?export CLAUDE_SHA256=<sha256 of claude binary>}"
: "${CREDS_OSS_URL:?export CREDS_OSS_URL=<OSS URL of encrypted credential bundle>}"
: "${CREDS_SHA256:?export CREDS_SHA256=<sha256 of encrypted bundle>}"
: "${CREDS_PASSPHRASE:?export CREDS_PASSPHRASE=<passphrase from worker-credentials-package.sh>}"

GIT_HOST="${GIT_HOST:-code.alibaba-inc.com}"
GIT_REPO_PATH="${GIT_REPO_PATH:-terraflow/jarvis-preview}"
JARVIS_ROOT="${JARVIS_ROOT:-$HOME/workspace/jarvis-preview}"
CLEAN_URL="https://$GIT_HOST/$GIT_REPO_PATH.git"

step "0. Base tools (git / curl)"
command -v curl >/dev/null 2>&1 || sudo yum install -y curl
command -v git  >/dev/null 2>&1 || sudo yum install -y git
ok "git = $(git --version | awk '{print $3}')  curl = $(curl --version | awk 'NR==1{print $2}')"

step "1. Clone $GIT_REPO_PATH → $JARVIS_ROOT"
if [ -d "$JARVIS_ROOT/.git" ]; then
  ok "repo already present; leaving as-is (worker-install pulls latest after creds land)"
else
  mkdir -p "$(dirname "$JARVIS_ROOT")"
  # Token goes into the clone URL only transiently; scrubbed from .git/config
  # right after. Subsequent pulls authenticate via the credential-store file
  # the bundle installs (~/.git-credentials).
  git clone "https://$GIT_TOKEN_USER:$GIT_TOKEN@$GIT_HOST/$GIT_REPO_PATH.git" "$JARVIS_ROOT" \
    || die "clone failed — check GIT_TOKEN/GIT_TOKEN_USER and repo read access"
  (cd "$JARVIS_ROOT" && git remote set-url origin "$CLEAN_URL")
  ok "cloned; origin scrubbed to $CLEAN_URL"
fi

step "2. Hand off to worker-install.sh"
[ -f "$JARVIS_ROOT/bootstrap/worker-install.sh" ] \
  || die "checkout has no bootstrap/worker-install.sh — wrong branch/repo?"
export JARVIS_ROOT
exec bash "$JARVIS_ROOT/bootstrap/worker-install.sh"
