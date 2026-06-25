#!/usr/bin/env bash
# Idempotent am bind from env vars. Source or run before any `am` send.
# Creds read from env only — nothing hardcoded.
#   DINGTALK_APP_KEY     appKey
#   DINGTALK_APP_SECRET  appSecret
#   DINGTALK_ROBOT_CODE  robotCode (defaults to appKey)
#   DINGTALK_STAFF_ID    operator/account staffId
set -euo pipefail

export PATH="$HOME/.local/bin:$PATH"

if ! command -v am >/dev/null 2>&1; then
  echo "ensure-bind: am CLI not found; install: curl -fsSL https://am.io.alibaba-inc.com/install.sh | bash" >&2
  exit 1
fi

# already bound → done
if am bind list 2>/dev/null | grep -q "bot"; then
  exit 0
fi

key="${DINGTALK_APP_KEY:-}"; secret="${DINGTALK_APP_SECRET:-}"
robot="${DINGTALK_ROBOT_CODE:-$key}"; staff="${DINGTALK_STAFF_ID:-}"
if [ -z "$key" ] || [ -z "$secret" ] || [ -z "$staff" ]; then
  echo "ensure-bind: set DINGTALK_APP_KEY / DINGTALK_APP_SECRET / DINGTALK_STAFF_ID" >&2
  exit 1
fi

am bind --type=bot --access-key-id="$key" --access-key-secret="$secret" \
  --account-id="$staff" --robot-code="$robot"
