#!/usr/bin/env bash
# bootstrap/provider-mirror-sync.sh — github→internal provider mirror pump.
#
# Runs ONLY on the scheduler host (macmini) — the single machine with GitHub
# reach (via its local proxy). Every other fleet machine consumes the internal
# mirror through sync-provider.sh; IDC workers cannot reach github.com at all
# (probe 2026-07-23: worker→github hangs, worker→macmini:7899 firewalled).
#
# What it syncs: refs/heads/* + refs/tags/* only. refs/pull/* is deliberately
# excluded (thousands of PR refs; GitLab would balloon and reviewers never
# need them via the mirror — PR/issue DATA is API-layer and stays on the
# scheduler-side gh flows).
#
# Freshness contract: each successful sync force-pushes a synthetic
# single-commit branch `jarvis-sync-heartbeat` (empty tree, committer date =
# sync time). Consumers measure mirror lag by that commit's age — they can't
# ask GitHub how far behind they are.
#
# Install (macmini crontab):
#   */10 * * * * $HOME/workspace/jarvis-preview/bootstrap/provider-mirror-sync.sh >> $HOME/.jarvis/mirror-sync.log 2>&1
#
# Env:
#   JARVIS_PROVIDER_UPSTREAM      default https://github.com/aliyun/terraform-provider-alicloud.git
#   JARVIS_PROVIDER_MIRROR_URL    default https://code.alibaba-inc.com/opensource-tools/terraform-provider-alicloud.git
#   JARVIS_PROVIDER_MIRROR_CACHE  default ~/.jarvis/mirrors/terraform-provider-alicloud.git
#   JARVIS_GITHUB_PROXY           e.g. http://127.0.0.1:7899 — applied to github.com ONLY
#
# The mirror target is the SHARED opensource-tools repo (reused per owner's
# call — no new repo). Sharing rules: we push exact upstream refs (any other
# faithful sync job is ref-compatible with ours), we NEVER --prune the remote
# (must not delete branches other people created there), and our only foreign
# ref is the jarvis-sync-heartbeat branch.
set -euo pipefail

UPSTREAM="${JARVIS_PROVIDER_UPSTREAM:-https://github.com/aliyun/terraform-provider-alicloud.git}"
MIRROR="${JARVIS_PROVIDER_MIRROR_URL:-https://code.alibaba-inc.com/opensource-tools/terraform-provider-alicloud.git}"
CACHE="${JARVIS_PROVIDER_MIRROR_CACHE:-$HOME/.jarvis/mirrors/terraform-provider-alicloud.git}"
GH_PROXY="${JARVIS_GITHUB_PROXY:-}"
HEARTBEAT_REF="refs/heads/jarvis-sync-heartbeat"

mkdir -p "$(dirname "$CACHE")"

# mkdir-based lock: macOS has no flock(1); overlapping cron runs must not
# interleave fetch/push on the same cache.
LOCK="$CACHE.lock"
if ! mkdir "$LOCK" 2>/dev/null; then
  echo "[mirror-sync] $(date '+%F %T') another run holds $LOCK — skipping"
  exit 0
fi
trap 'rmdir "$LOCK" 2>/dev/null || true' EXIT

# Proxy scoped to github.com only: the push to code.alibaba-inc.com must NOT
# go through the proxy (internal host, proxy would refuse/leak).
GIT=(git)
[ -n "$GH_PROXY" ] && GIT=(git -c "http.https://github.com/.proxy=$GH_PROXY")

if [ ! -d "$CACHE" ]; then
  echo "[mirror-sync] $(date '+%F %T') init bare cache $CACHE"
  git init --bare -q "$CACHE"
  git -C "$CACHE" remote add origin "$UPSTREAM"
  git -C "$CACHE" config remote.origin.fetch '+refs/heads/*:refs/heads/*'
  git -C "$CACHE" config --add remote.origin.fetch '+refs/tags/*:refs/tags/*'
fi

"${GIT[@]}" -C "$CACHE" fetch -q --prune origin

# Heartbeat: fresh synthetic commit each run so committer date = sync time.
empty_tree=$(git -C "$CACHE" hash-object -t tree /dev/null)
hb_commit=$(GIT_DIR="$CACHE" git commit-tree "$empty_tree" -m "jarvis mirror sync heartbeat")
git -C "$CACHE" update-ref "$HEARTBEAT_REF" "$hb_commit"

# Forced refspecs (+): heartbeat rewrites history every run, and upstream
# branch rewinds (rare) must also win. NO --prune: the target is a shared
# repo — deleting refs other teams created there is not ours to do; branches
# deleted upstream simply linger.
git -C "$CACHE" push -q "$MIRROR" \
  '+refs/heads/*:refs/heads/*' '+refs/tags/*:refs/tags/*'

echo "[mirror-sync] $(date '+%F %T') master=$(git -C "$CACHE" rev-parse --short refs/heads/master) heartbeat pushed"
