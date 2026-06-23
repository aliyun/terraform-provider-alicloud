#!/usr/bin/env bash
# Ensure the alicloud provider repo is synced. Path is read from config/workspaces.json
# (.workspaces.terraform_provider.path) — edit there to relocate. No repo -> clone; existing repo -> fetch only.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG="$SCRIPT_DIR/../../../../config/workspaces.json"
REPO_DIR="$(jq -r '.workspaces.terraform_provider.path' "$CONFIG")"
REPO_DIR="${REPO_DIR/#\~/$HOME}"
REMOTE="https://github.com/aliyun/terraform-provider-alicloud.git"

if [ ! -d "$REPO_DIR/.git" ]; then
  echo "[sync-provider] cloning into $REPO_DIR ..."
  mkdir -p "$(dirname "$REPO_DIR")"
  git clone --depth 1 "$REMOTE" "$REPO_DIR"
else
  echo "[sync-provider] updating $REPO_DIR ..."
  DEF=$(git -C "$REPO_DIR" remote show origin | sed -n 's/.*HEAD branch: //p')
  DEF="${DEF:-master}"
  git -C "$REPO_DIR" fetch --depth 1 origin "$DEF"
fi
echo "[sync-provider] ready: $REPO_DIR @ $(git -C "$REPO_DIR" rev-parse --short HEAD)"
