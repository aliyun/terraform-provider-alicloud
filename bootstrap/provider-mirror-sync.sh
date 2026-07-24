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
# Freshness contract: each successful sync advances a synthetic heartbeat
# branch `jarvis-sync-heartbeat` (empty-tree commits chained parent-wise,
# committer date = sync time — fast-forward only, the server denies non-ff).
# Consumers measure mirror lag by the tip commit's age — they can't ask
# GitHub how far behind they are.
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
TMPD=$(mktemp -d)
trap 'rmdir "$LOCK" 2>/dev/null; rm -rf "$TMPD"' EXIT

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

# Remote ref inventory FIRST — the heartbeat chains onto its remote tip so
# every update stays fast-forward.
git ls-remote "$MIRROR" 'refs/heads/*' 'refs/tags/*' \
  | awk '$2 !~ /\^\{\}$/ {print $1" "$2}' | sort > "$TMPD/remote"

# Heartbeat: fresh synthetic commit each run (committer date = sync time),
# chained parent-wise so the branch only ever fast-forwards. Parent = remote
# tip when its object is in the cache (it always is — we authored it),
# else local tip, else root. If the cache was rebuilt and the chain broke,
# the ref lands in the divergent-skip list with a remedy hint below.
empty_tree=$(git -C "$CACHE" hash-object -t tree /dev/null)
hb_parent=$(awk -v r="$HEARTBEAT_REF" '$2==r{print $1}' "$TMPD/remote")
if [ -n "$hb_parent" ] && ! git -C "$CACHE" cat-file -e "$hb_parent" 2>/dev/null; then
  hb_parent=""
fi
[ -z "$hb_parent" ] && hb_parent=$(git -C "$CACHE" rev-parse -q --verify "$HEARTBEAT_REF" 2>/dev/null || true)
if [ -n "$hb_parent" ]; then
  hb_commit=$(GIT_DIR="$CACHE" git commit-tree "$empty_tree" -p "$hb_parent" -m "jarvis mirror sync heartbeat")
else
  hb_commit=$(GIT_DIR="$CACHE" git commit-tree "$empty_tree" -m "jarvis mirror sync heartbeat")
fi
git -C "$CACHE" update-ref "$HEARTBEAT_REF" "$hb_commit"

# Forced refspecs (+): heartbeat rewrites history every run, and upstream
# branch rewinds (rare) must also win. NO --prune: the target is a shared
# repo — deleting refs other teams created there is not ours to do; branches
# deleted upstream simply linger.
#
# The shared repo predates this pump and carries its OWN lineage for many
# branches plus re-created tags, and the server DENIES non-fast-forward
# updates (probe-verified: creates and fast-forwards land, any divergent
# update comes back "update-ref failed" — even singly). Policy: push only
# CREATES and FAST-FORWARDS, never force; divergent refs are skipped and
# stay frozen at the repo's historical state. Consumers track master
# (clean fast-forward) and new upstream branches/tags arrive as creates.
git -C "$CACHE" for-each-ref --format='%(objectname) %(refname)' refs/heads refs/tags \
  | sort > "$TMPD/local"
changed=$(comm -23 "$TMPD/local" "$TMPD/remote" | awk '{print $2}')

specs=(); skipped=0
while IFS= read -r ref; do
  [ -z "$ref" ] && continue
  local_sha=$(git -C "$CACHE" rev-parse -q --verify "$ref" 2>/dev/null || true)
  [ -z "$local_sha" ] && continue
  remote_sha=$(awk -v r="$ref" '$2==r{print $1}' "$TMPD/remote")
  # Absent remotely = create; remote tip an ancestor of ours = fast-forward.
  # Anything else (incl. annotated-tag rewrites, foreign history whose
  # objects we don't even have) classifies as divergent and is skipped.
  if [ -z "$remote_sha" ] \
     || git -C "$CACHE" merge-base --is-ancestor "$remote_sha" "$local_sha" 2>/dev/null; then
    specs+=("$ref:$ref")
  else
    skipped=$((skipped+1))
    echo "$ref" >> "$TMPD/skipped"
  fi
done <<< "$changed"

# Adaptive batch: try a slice; on rejection split in half and retry down to
# singles (the server's multi-ref tolerance is erratic — 418 refs and even
# 40-ref batches were rejected wholesale). Retry spew goes to a file so the
# cron log stays readable; a single-ref rejection aborts loudly.
pushes=0
push_refs() {
  [ "$#" -eq 0 ] && return 0
  if git -C "$CACHE" push -q "$MIRROR" "$@" 2>>"$TMPD/push.err"; then
    pushes=$((pushes+1))
    return 0
  fi
  if [ "$#" -le 1 ]; then
    echo "[mirror-sync] ERROR: single-ref push rejected: $*" >&2
    tail -5 "$TMPD/push.err" >&2
    return 1
  fi
  local mid=$(( $# / 2 ))
  push_refs "${@:1:$mid}"
  push_refs "${@:$((mid+1))}"
}

total=${#specs[@]}
i=0
while [ "$i" -lt "$total" ]; do
  push_refs "${specs[@]:$i:40}"
  i=$((i+40))
done

if [ "$skipped" -gt 0 ]; then
  echo "[mirror-sync] skipped $skipped divergent ref(s) (server denies non-ff; shared-repo history left untouched): $(tr '\n' ' ' < "$TMPD/skipped" | cut -c1-400)"
  if grep -q "jarvis-sync-heartbeat" "$TMPD/skipped" 2>/dev/null; then
    echo "[mirror-sync] WARNING: heartbeat diverged (cache rebuilt?) — delete the remote jarvis-sync-heartbeat branch once so the next run recreates it" >&2
  fi
fi
echo "[mirror-sync] $(date '+%F %T') master=$(git -C "$CACHE" rev-parse --short refs/heads/master) pushed $total ref(s) in $pushes push(es), skipped $skipped"
