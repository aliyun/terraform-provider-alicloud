#!/usr/bin/env bash
# test/sync_provider_dedup_test.sh — 校验 aone-triage 与 terraform-pr-review 两 skill 内
# scripts/sync-provider.sh 内容 byte-identical(P3.c 收敛后应始终一致)。
#
# 未来若再需要真单点维护(避免拷贝),把该脚本迁到 bootstrap/sync-provider.sh 并统一引用。

set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

AONE="$repo_root/.claude/skills/aone-triage/scripts/sync-provider.sh"
PRR="$repo_root/.claude/skills/terraform-pr-review/scripts/sync-provider.sh"

if diff -u "$AONE" "$PRR"; then
    echo "sync_provider_dedup_test: PASS(两 skill 版 byte-identical)"
    exit 0
else
    echo "sync_provider_dedup_test: FAIL"
    echo ""
    echo "aone-triage 与 terraform-pr-review 的 sync-provider.sh 出现差异,违反 P3.c 单点约定。"
    echo "修复:选一份为准,cp 覆盖另一份;或迁到 bootstrap/sync-provider.sh 做真单点维护。"
    exit 1
fi
