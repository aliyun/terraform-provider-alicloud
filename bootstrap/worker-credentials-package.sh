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
#   ~/.aliyun/config.json                      aliyun CLI AK (OpenAPI 查证 / cred check)
#   $JARVIS_ROOT/bootstrap/.env                GitHub PAT, worker control-plane token, etc.
#   $JARVIS_ROOT/bridge/jarvis.env             DingTalk, model lanes, JARVIS_* vars
#   The operator-only JARVIS_CONTROL_PLANE_ADMIN_TOKEN is stripped from both
#   env files before packaging and must never be installed on worker hosts.
#   $JARVIS_ROOT/config/workspaces.local.json  (if present) — sanitized/re-derived
#                                              on worker; this ships as reference
#
# Crypto: openssl aes-256-cbc + SHA-256 KDF. Works on OpenSSL 1.0.2 (AliOS 7.2
# default, no `-pbkdf2` flag) through 3.x (macOS Homebrew). We rely on the
# 256-bit random passphrase for security rather than KDF strength; a stronger
# KDF (PBKDF2 / scrypt / argon2) matters when passphrases have low entropy,
# not when they're `openssl rand -hex 32`.
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

# Git auth for workers to `git clone` / `git pull` jarvis-preview via HTTPS +
# GitLab Deploy Token. Create the token via web UI:
#   jarvis-preview → Settings → Repository → Deploy Tokens → Add token
#   Scope: read_repository ONLY
#   Save both the username (e.g. 'gitlab+deploy-token-12345') and secret
#   — GitLab shows the secret only once.
: "${GIT_TOKEN:?export GIT_TOKEN=<GitLab Deploy Token secret; see script header for setup>}"
: "${GIT_TOKEN_USER:?export GIT_TOKEN_USER=<username GitLab printed, e.g. gitlab+deploy-token-12345>}"
GIT_HOST="${GIT_HOST:-code.alibaba-inc.com}"
GIT_REPO_PATH="${GIT_REPO_PATH:-terraflow/jarvis-preview}"

# Resolve JARVIS_ROOT = the repo containing this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]:-$0}")" && pwd)"
JARVIS_ROOT="${JARVIS_ROOT:-$(cd "$SCRIPT_DIR/.." && pwd)}"
[ -f "$JARVIS_ROOT/bridge/run.sh" ] \
  || die "JARVIS_ROOT=$JARVIS_ROOT doesn't look like a jarvis checkout"

command -v openssl >/dev/null 2>&1 || die "openssl not found on PATH"
command -v tar >/dev/null 2>&1 || die "tar not found on PATH"

# ossutil is OPTIONAL: if installed AND runnable, we upload automatically.
# If missing (or the file exists but isn't a working ossutil — e.g., an HTML
# error page saved by a botched install), we still do the tar+encrypt work,
# save the ciphertext to a durable location, and print instructions so the
# operator can upload via web console / ossbrowser / any other tool.
OSSUTIL_AVAILABLE=0
if command -v "$OSSUTIL" >/dev/null 2>&1; then
  # Verify it's actually a working ossutil, not a mis-downloaded HTML/XML page
  # or an unrelated binary at the same name. --version is the cheapest probe.
  if ossutil_version=$("$OSSUTIL" --version 2>&1) && [ -n "$ossutil_version" ]; then
    OSSUTIL_AVAILABLE=1
    ok "ossutil detected: $(printf '%s' "$ossutil_version" | head -1)"
  else
    warn "ossutil found at $(command -v "$OSSUTIL") but --version failed:"
    printf '%s\n' "$ossutil_version" | head -3 | sed 's/^/       /' >&2
    warn "treating as unavailable; will save encrypted bundle locally for manual upload"
    warn "  (likely a broken install — remove and re-install from https://help.aliyun.com/document_detail/120075.html if you want auto-upload)"
  fi
else
  warn "ossutil not found — will package + encrypt locally and print manual upload instructions."
  warn "  To auto-upload next time: install ossutil (https://help.aliyun.com/document_detail/120075.html)"
fi

# ---------------------------------------------------------------------------
step "1. Sanity check credential sources"
missing=()
[ -d "$HOME/.config/a1" ] || missing+=("~/.config/a1/")
[ -d "$HOME/.claude" ]    || missing+=("~/.claude/")
# verify.sh hard-checks `aliyun sts GetCallerIdentity` on every role — a bundle
# without ~/.aliyun/config.json produces workers that fail preflight daily.
[ -f "$HOME/.aliyun/config.json" ] || missing+=("~/.aliyun/config.json")
[ -f "$JARVIS_ROOT/bootstrap/.env" ] \
  || warn "$JARVIS_ROOT/bootstrap/.env missing (workers won't have GitHub/control-plane tokens)"
[ -f "$JARVIS_ROOT/bridge/jarvis.env" ] \
  || warn "$JARVIS_ROOT/bridge/jarvis.env missing (workers won't have model-lane / DingTalk config)"
if [ "${#missing[@]}" -gt 0 ]; then
  die "required credential sources missing: ${missing[*]}"
fi
ok "sources present"

# ---------------------------------------------------------------------------
step "1.5 Validate credentials are ALIVE (never ship dead creds)"
# Root cause of the first canary's WARN jarvis-a1-session: the bundle was cut
# while the packager's a1 session was expired, so every worker inherited a dead
# session. Liveness is cheap to probe here and expensive to discover on N
# workers — hard-fail the packager instead.

# a1: whoami must resolve to the jarvis digital-worker account.
a1_expect="WORKER_1782379562571"
a1_acct=$(JARVIS_A1_IDENTITY=jarvis "$JARVIS_ROOT/bin/a1id" -- auth whoami 2>/dev/null \
            | awk '/Account:/{print $2}' || true)
[ "$a1_acct" = "$a1_expect" ] \
  || die "a1 jarvis session dead/expired (whoami Account='${a1_acct:-<empty>}' != $a1_expect).
       Fix: browser-login BUC as open_jarvis, run \`bin/a1id login jarvis\`, then re-run this packager."
ok "a1 session alive (Account=$a1_acct)"

# GitHub PAT (bootstrap/.env → JARVIS_GITHUB_TOKEN): workers use it for the
# whole GitHub path (github-identity.sh check/push, PR flows). Probe with the
# token AS READ FROM .env — unset first so a token lingering in the calling
# shell's environment can't mask an .env that ships without one, and require
# non-empty so gh can't fall back to ambient stored auth (禁回退, 纪律 #6).
if command -v gh >/dev/null 2>&1 && [ -f "$JARVIS_ROOT/bootstrap/.env" ]; then
  gh_token=$( (unset JARVIS_GITHUB_TOKEN
               source "$JARVIS_ROOT/bootstrap/.env" >/dev/null 2>&1 || true
               printf '%s' "${JARVIS_GITHUB_TOKEN:-}") )
  [ -n "$gh_token" ] \
    || die "bootstrap/.env carries no JARVIS_GITHUB_TOKEN — workers would ship without the GitHub PAT.
       Fix: add JARVIS_GITHUB_TOKEN=<api-tool-agent PAT> to bootstrap/.env, then re-run this packager."
  gh_login=$(GH_TOKEN="$gh_token" gh api user --jq .login 2>/dev/null || true)
  [ "$gh_login" = "api-tool-agent" ] \
    || die "JARVIS_GITHUB_TOKEN in bootstrap/.env invalid (gh api user → '${gh_login:-<empty>}', want api-tool-agent).
       Fix: refresh the api-tool-agent token, update bootstrap/.env, then re-run this packager."
  ok "GitHub token alive (login=$gh_login)"
else
  warn "gh CLI or bootstrap/.env missing on packager — skipping GitHub token liveness probe"
fi

# aliyun AK (~/.aliyun/config.json): workers use it for OpenAPI 查证.
if command -v aliyun >/dev/null 2>&1; then
  aliyun sts GetCallerIdentity >/dev/null 2>&1 \
    || die "aliyun credential invalid on packager (sts GetCallerIdentity failed).
       Fix: repair ~/.aliyun/config.json (aliyun configure), then re-run this packager."
  ok "aliyun credential alive"
else
  warn "aliyun CLI missing on packager PATH — skipping aliyun cred liveness probe"
fi

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

# ~/.aliyun: config.json only (profiles + AK). plugins/ is per-arch binary
# cache the worker regenerates on its own — shipping it bloats the bundle.
mkdir -p "$STAGE/home/.aliyun"
cp -p "$HOME/.aliyun/config.json" "$STAGE/home/.aliyun/config.json"
chmod 600 "$STAGE/home/.aliyun/config.json"

# Git auth: HTTPS + Deploy Token via git's `store` credential helper.
# Package ~/.git-credentials + ~/.gitconfig so worker can `git pull` without
# openssh-clients (AliOS minimal install omits it). Deploy Token is repo-
# scoped read-only + revocable in GitLab web UI.
printf 'https://%s:%s@%s\n' "$GIT_TOKEN_USER" "$GIT_TOKEN" "$GIT_HOST" \
  > "$STAGE/home/.git-credentials"
chmod 600 "$STAGE/home/.git-credentials"
cat > "$STAGE/home/.gitconfig" <<GITCFG
[credential]
    helper = store
[credential "https://$GIT_HOST"]
    helper = store
# Workers have no openssh-clients (by design — HTTPS + token only), but
# config/workspaces.json registers internal repos with ssh URLs (works on the
# scheduler mac via ssh keys). Rewrite them transparently so lazy workspace
# clones just work through the credential store above.
[url "https://$GIT_HOST/"]
    insteadOf = git@gitlab.alibaba-inc.com:
GITCFG
chmod 644 "$STAGE/home/.gitconfig"
# Record the HTTPS URL worker-install.sh should set as `origin`, overriding
# any ssh URL the tarballed repo may have inherited from the packager's clone.
printf 'GIT_HTTPS_URL=https://%s/%s.git\n' "$GIT_HOST" "$GIT_REPO_PATH" \
  >> "$STAGE/MANIFEST"
ok "packaged git-credentials for Deploy Token ${GIT_TOKEN_USER}@${GIT_HOST}/${GIT_REPO_PATH}"
if [ -f "$JARVIS_ROOT/bootstrap/.env" ]; then
  bash "$SCRIPT_DIR/worker-env-sanitize.sh" \
    "$JARVIS_ROOT/bootstrap/.env" "$STAGE/repo/bootstrap/.env"
fi
if [ -f "$JARVIS_ROOT/bridge/jarvis.env" ]; then
  bash "$SCRIPT_DIR/worker-env-sanitize.sh" \
    "$JARVIS_ROOT/bridge/jarvis.env" "$STAGE/repo/bridge/jarvis.env"
fi
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
openssl enc -aes-256-cbc -salt -md sha256 \
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
