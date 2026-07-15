#!/usr/bin/env bash
# thin wrapper — 真源在 bootstrap/sync-provider.sh(单点维护,勿在此加逻辑);
# 本文件仅保住既有调用路径(.Codex/skills/*/scripts/sync-provider.sh 及 .agents 镜像)。
exec bash "$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd)/bootstrap/sync-provider.sh" "$@"
