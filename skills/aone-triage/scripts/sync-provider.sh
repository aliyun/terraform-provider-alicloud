#!/usr/bin/env bash
# Ensure ~/terraflow/providers/alicloud holds the latest terraform-provider-alicloud master.
# No repo -> clone; existing repo -> fetch + hard-reset to origin default branch.
set -euo pipefail

REPO_DIR="${TERRAFLOW_ALICLOUD:-$HOME/terraflow/providers/alicloud}"
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
  git -C "$REPO_DIR" reset --hard "origin/$DEF"
fi
echo "[sync-provider] ready: $REPO_DIR @ $(git -C "$REPO_DIR" rev-parse --short HEAD)"
