#!/usr/bin/env bash
# bootstrap/worker-credentials-package.sh
#
# Run on the SCHEDULER host (the machine with working a1/claude/GitHub creds).
# Packages the credential set that a Linux worker needs, encrypts it, uploads to
# OSS, and prints the three env vars a worker needs to consume it.
#
# Push-only distribution: no ssh, no rsync from workers. Workers can't reach the
# scheduler and have no browser for interactive SSO. This script + worker-install.sh
# form the credential delivery loop.
#
# What's packaged:
#   ~/.config/a1/                              a1 CLI login state (all identities)
#   ~/.claude/                                 claude settings + gateway/token files
#   $JARVIS_ROOT/bootstrap/.env                GitHub PAT, control-plane token, etc.
#   $JARVIS_ROOT/bridge/jarvis.env             DingTalk, model lanes, JARVIS_* vars
#   $JARVIS_ROOT/config/workspaces.local.json  (if present) — sanitized/re-derived
#                                              on worker; this ships as reference
#
# Crypto: openssl aes-256-cbc + pbkdf2 + 100k iters. Works on OpenSSL 1.0.2+
# (AliOS 7.2 default) and 3.x (macOS Homebrew).
#
# Required env vars:
#   OSS_BUCKET       e.g. cc-packet
#
# Optional env vars:
#   OSS_ENDPOINT     https://oss-cn-beijing-internal.aliyuncs.com  (worker-facing)
#   OSS_KEY          object key; default: jarvis-worker-creds-<ts>.tar.gz.enc
#   JARVIS_ROOT      default: this checkout's repo root
#   OSSUTIL          default: `ossutil` on PATH
#
# Usage:
#   OSS_BUCKET=cc-packet bash bootstrap/worker-credentials-package.sh
#
# Output: URL / SHA256 / PASSPHRASE (feed to worker-install.sh on each worker).

set -euo pipefail

die()  { printf '[ERR]  %s\n' "$*" >&2; exit 1; }
info() { printf '[..]   %s\n' "$*"; }
ok()   { printf '[OK]   %s\n' "$*"; }
warn() { printf '[WARN] %s\n' "$*" >&2; }
step() { printf '\n===== %s =====\n' "$*"; }

: "${OSS_BUCKET:?export OSS_BUCKET=<your OSS bucket name>}"

OSS_ENDPOINT="${OSS_ENDPOINT:-https://oss-cn-beijing-internal.aliyuncs.com}"
OSS_KEY="${OSS_KEY:-jarvis-worker-creds-$(date +%Y%m%d-%H%M%S).tar.gz.enc}"
OSSUTIL="${OSSUTIL:-ossutil}"

# Resolve JARVIS_ROOT = the repo containing this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]:-$0}")" && pwd)"
JARVIS_ROOT="${JARVIS_ROOT:-$(cd "$SCRIPT_DIR/.." && pwd)}"
[ -f "$JARVIS_ROOT/bridge/run.sh" ] \
  || die "JARVIS_ROOT=$JARVIS_ROOT doesn't look like a jarvis checkout"

command -v openssl >/dev/null 2>&1 || die "openssl not found on PATH"
command -v tar >/dev/null 2>&1 || die "tar not found on PATH"

# ossutil is OPTIONAL: if installed, we upload automatically. If missing, we
# still do the tar+encrypt work, save the ciphertext to a durable location,
# and print instructions so the operator can upload via web console /
# ossbrowser / any other tool. This makes the packager usable on hosts where
# installing ossutil is inconvenient (e.g., macOS without brew tap set up).
OSSUTIL_AVAILABLE=0
if command -v "$OSSUTIL" >/dev/null 2>&1; then
  OSSUTIL_AVAILABLE=1
else
  warn "ossutil not found — will package + encrypt locally and print manual upload instructions."
  warn "  To auto-upload next time: install ossutil (https://help.aliyun.com/document_detail/120075.html)"
fi

# ---------------------------------------------------------------------------
step "1. Sanity check credential sources"
missing=()
[ -d "$HOME/.config/a1" ] || missing+=("~/.config/a1/")
[ -d "$HOME/.claude" ]    || missing+=("~/.claude/")
[ -f "$JARVIS_ROOT/bootstrap/.env" ] \
  || warn "$JARVIS_ROOT/bootstrap/.env missing (workers won't have GitHub/control-plane tokens)"
[ -f "$JARVIS_ROOT/bridge/jarvis.env" ] \
  || warn "$JARVIS_ROOT/bridge/jarvis.env missing (workers won't have model-lane / DingTalk config)"
if [ "${#missing[@]}" -gt 0 ]; then
  die "required credential dirs missing: ${missing[*]}"
fi
ok "sources present"

# ---------------------------------------------------------------------------
step "2. Warn on Mac-specific absolute paths in claude settings"
# claude settings JSONs sometimes reference $HOME with the packager's actual path.
# On a Linux worker with a different $HOME (e.g. /home/admin vs /Users/gzzz),
# hardcoded paths become dead links. Warn now so operator can fix before ship.
if grep -RlEo '/Users/[^"[:space:]]+' "$HOME/.claude" 2>/dev/null | head -1 >/dev/null; then
  warn "Mac-specific /Users/* paths found in ~/.claude/ (may need worker-side sanitize):"
  grep -REno '/Users/[^"[:space:]]+' "$HOME/.claude" 2>/dev/null | head -5 | sed 's/^/       /' >&2
  warn "worker-install.sh step 5 rewrites /Users/<macuser> → \$HOME after extract, but verify manually if paths reference NON-HOME locations (e.g. /Users/gzzz/workspace/…)"
else
  ok "no /Users/* absolute paths detected in ~/.claude/"
fi

if [ -f "$JARVIS_ROOT/bootstrap/.env" ]; then
  if grep -qEo '/Users/[^"[:space:]]+' "$JARVIS_ROOT/bootstrap/.env" 2>/dev/null; then
    warn "bootstrap/.env contains /Users/* paths; may need JARVIS_ROOT / JARVIS_TATA_ROOT overrides on worker"
    grep -Eno '/Users/[^"[:space:]]+' "$JARVIS_ROOT/bootstrap/.env" | head -3 | sed 's/^/       /' >&2
  fi
fi
if [ -f "$JARVIS_ROOT/bridge/jarvis.env" ]; then
  if grep -qEo '/Users/[^"[:space:]]+' "$JARVIS_ROOT/bridge/jarvis.env" 2>/dev/null; then
    warn "bridge/jarvis.env contains /Users/* paths (JARVIS_SETTINGS, CLAUDE_BIN, JARVIS_TATA_SETTINGS, etc)"
    grep -Eno '/Users/[^"[:space:]]+' "$JARVIS_ROOT/bridge/jarvis.env" | head -5 | sed 's/^/       /' >&2
    warn "worker-install.sh rewrites /Users/<macuser> → \$HOME, but paths under /Users/<macuser>/workspace or /Users/<macuser>/.claude/idea_settings.json will remap OK; anything ELSE needs manual review"
  fi
fi

# ---------------------------------------------------------------------------
step "3. Stage credential tree in a temp dir"
STAGE=$(mktemp -d)
trap 'rm -rf "$STAGE"' EXIT

mkdir -p "$STAGE/home/.config" "$STAGE/home/.claude" \
         "$STAGE/repo/bootstrap" "$STAGE/repo/bridge" "$STAGE/repo/config"

# Copy — preserve mode/owner so 0600 secrets stay 0600 on extract.
# ~/.config/a1: whole dir (login state per identity, all small JSON tokens).
cp -Rp "$HOME/.config/a1"  "$STAGE/home/.config/a1"

# ~/.claude: SELECTIVE. Ship only the top-level *.json settings archives that
# claude --settings points to (glm5.2.json, idea_settings.json, etc). Skip:
#   tasks/, projects/, todos/, statsig/, shell-snapshots/, ide/, __store.db*
# — all per-machine transient state (session transcripts, task/hook state,
# telemetry, IDE snapshots) that workers must NOT inherit from the packager.
# Shipping them bloats the bundle by tens/hundreds of MB and leaks past
# conversations to workers unnecessarily.
find "$HOME/.claude" -maxdepth 1 -type f \
  \( -name '*.json' -o -name 'CLAUDE.md' -o -name '.credentials.json' \) \
  -exec cp -p {} "$STAGE/home/.claude/" \; 2>/dev/null || true

n_json=$(find "$STAGE/home/.claude" -maxdepth 1 -type f -name '*.json' | wc -l | tr -d ' ')
[ "$n_json" -gt 0 ] || warn "no ~/.claude/*.json files found — worker may lack gateway/settings config"
[ -f "$JARVIS_ROOT/bootstrap/.env" ] \
  && cp -p "$JARVIS_ROOT/bootstrap/.env" "$STAGE/repo/bootstrap/.env"
[ -f "$JARVIS_ROOT/bridge/jarvis.env" ] \
  && cp -p "$JARVIS_ROOT/bridge/jarvis.env" "$STAGE/repo/bridge/jarvis.env"
[ -f "$JARVIS_ROOT/config/workspaces.local.json" ] \
  && cp -p "$JARVIS_ROOT/config/workspaces.local.json" "$STAGE/repo/config/workspaces.local.json"

# Record the packager's $HOME so worker-install can rewrite /Users/<mac-user>/…
# → /home/<linux-user>/… where safely possible (paths under $HOME map cleanly).
printf 'PACKAGER_HOME=%s\n' "$HOME" > "$STAGE/MANIFEST"
printf 'PACKAGER_USER=%s\n' "$(id -un)" >> "$STAGE/MANIFEST"
printf 'PACKAGER_HOST=%s\n' "$(hostname -s 2>/dev/null || hostname)" >> "$STAGE/MANIFEST"
printf 'PACKAGED_AT=%s\n' "$(date -u +%FT%TZ)" >> "$STAGE/MANIFEST"

info "staged tree:"
(cd "$STAGE" && find . -maxdepth 4 -type d | head -15 | sed 's/^/  /')

# ---------------------------------------------------------------------------
step "4. Tar + encrypt (AES-256-CBC + PBKDF2 100k iters)"
PLAINTAR="$STAGE/creds.tar.gz"
CIPHER="$STAGE/creds.tar.gz.enc"

tar czf "$PLAINTAR" -C "$STAGE" home repo MANIFEST
ok "tarball: $(stat -f '%z' "$PLAINTAR" 2>/dev/null || stat -c '%s' "$PLAINTAR") bytes"

PASSPHRASE=$(openssl rand -hex 32)  # 256 bits of entropy
openssl enc -aes-256-cbc -salt -pbkdf2 -iter 100000 \
  -in "$PLAINTAR" -out "$CIPHER" \
  -pass "pass:$PASSPHRASE"

SHA=$(sha256sum "$CIPHER" 2>/dev/null | awk '{print $1}')
[ -n "$SHA" ] || SHA=$(shasum -a 256 "$CIPHER" | awk '{print $1}')
ok "encrypted: $(stat -f '%z' "$CIPHER" 2>/dev/null || stat -c '%s' "$CIPHER") bytes  sha256=$SHA"

# Overwrite plaintext with zeros before removal (best-effort — /tmp on macOS
# is usually memory-backed, but be paranoid).
dd if=/dev/zero of="$PLAINTAR" bs=1M count=1 conv=notrunc 2>/dev/null || true
rm -f "$PLAINTAR"

# ---------------------------------------------------------------------------
# Compute the expected worker-facing URL — works whether we upload here or
# operator uploads via web console / ossbrowser to the same bucket+key.
if [[ "$OSS_ENDPOINT" == *"://$OSS_BUCKET."* ]]; then
  URL="$OSS_ENDPOINT/$OSS_KEY"
else
  URL="${OSS_ENDPOINT%/}"
  URL="${URL/#https:\/\//https:\/\/$OSS_BUCKET.}"
  URL="${URL/#http:\/\//http:\/\/$OSS_BUCKET.}"
  URL="$URL/$OSS_KEY"
fi

# Persist the ciphertext outside the temp dir so it survives after the script
# exits (the temp dir is trap-cleaned). Operator can upload it manually if
# ossutil wasn't available.
OUT_FILE="${OUT_DIR:-$HOME}/${OSS_KEY##*/}"
cp -p "$CIPHER" "$OUT_FILE"
chmod 600 "$OUT_FILE"

if [ "$OSSUTIL_AVAILABLE" = 1 ]; then
  step "5. Upload to OSS via ossutil"
  "$OSSUTIL" cp -f "$OUT_FILE" "oss://$OSS_BUCKET/$OSS_KEY" 2>&1 | tail -3
  UPLOAD_STATE=uploaded
else
  step "5. SKIP auto-upload (ossutil not installed)"
  info "ciphertext saved to $OUT_FILE (mode 600)"
  info "upload via web console / ossbrowser as object key: $OSS_KEY"
  UPLOAD_STATE=manual
fi

# ---------------------------------------------------------------------------
step "6. Distribute these three values to each worker (secure channel)"
if [ "$UPLOAD_STATE" = uploaded ]; then
  cat <<EOF

═══════════════════════════════════════════════════════════════════════════════
CREDENTIAL BUNDLE UPLOADED — DISTRIBUTE THESE 3 VALUES

  CREDS_OSS_URL="$URL"
  CREDS_SHA256="$SHA"
  CREDS_PASSPHRASE="$PASSPHRASE"

On each worker:

  CLAUDE_OSS_URL="https://$OSS_BUCKET.oss-cn-beijing-internal.aliyuncs.com/claude" \\
  CLAUDE_SHA256="c1efffaaf370aa187cb6a09dd93d4e511c646899b0078476f83791b664bde7fe" \\
  CREDS_OSS_URL="$URL" \\
  CREDS_SHA256="$SHA" \\
  CREDS_PASSPHRASE="$PASSPHRASE" \\
  JARVIS_DISPATCH_MAX=3 \\
    bash bootstrap/worker-install.sh

After all workers install, REVOKE by deleting the OSS object:
  $OSSUTIL rm oss://$OSS_BUCKET/$OSS_KEY

The passphrase is 256-bit random; leaked without the OSS object it grants
nothing. Deleting the object cuts the last shred of usefulness.
═══════════════════════════════════════════════════════════════════════════════
EOF
else
  cat <<EOF

═══════════════════════════════════════════════════════════════════════════════
CREDENTIAL BUNDLE READY — MANUAL UPLOAD REQUIRED

Local ciphertext:
  $OUT_FILE
  ($(wc -c <"$OUT_FILE" | tr -d ' ') bytes, sha256=$SHA, mode 600)

Upload it to your OSS bucket via web console / ossbrowser / any tool:
  bucket:     $OSS_BUCKET
  object key: $OSS_KEY

After upload, distribute THESE 3 VALUES to each worker:

  CREDS_OSS_URL="$URL"
  CREDS_SHA256="$SHA"
  CREDS_PASSPHRASE="$PASSPHRASE"

On each worker:

  CLAUDE_OSS_URL="https://$OSS_BUCKET.oss-cn-beijing-internal.aliyuncs.com/claude" \\
  CLAUDE_SHA256="c1efffaaf370aa187cb6a09dd93d4e511c646899b0078476f83791b664bde7fe" \\
  CREDS_OSS_URL="$URL" \\
  CREDS_SHA256="$SHA" \\
  CREDS_PASSPHRASE="$PASSPHRASE" \\
  JARVIS_DISPATCH_MAX=3 \\
    bash bootstrap/worker-install.sh

After all workers install, REVOKE by deleting the OSS object (via console /
ossbrowser / ossutil) — the passphrase alone grants nothing.

Also delete the local ciphertext to leave no trace:
  shred -u "$OUT_FILE"   # linux
  rm -P "$OUT_FILE"      # macOS
═══════════════════════════════════════════════════════════════════════════════
EOF
fi
