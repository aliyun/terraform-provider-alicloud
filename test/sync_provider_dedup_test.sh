#!/usr/bin/env bash
# test/sync_provider_dedup_test.sh — sync-provider 单点约定:
#   真源 = bootstrap/sync-provider.sh(唯一含同步逻辑之处);
#   两个 skill 的 scripts/sync-provider.sh 是 byte-identical 薄 wrapper,只 exec 真源。
# 违反任一条 = 单点约定破坏(历史上双拷贝曾真实分叉过,勿回退)。

set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

CANON="$repo_root/bootstrap/sync-provider.sh"
AONE="$repo_root/.claude/skills/aone-triage/scripts/sync-provider.sh"
PRR="$repo_root/.claude/skills/terraform-pr-review/scripts/sync-provider.sh"

fail=0

if [ ! -f "$CANON" ]; then
    echo "sync_provider_dedup_test: 缺真源 bootstrap/sync-provider.sh" >&2
    fail=1
elif ! grep -q "reset --hard FETCH_HEAD" "$CANON"; then
    echo "sync_provider_dedup_test: 真源缺同步逻辑(fetch + reset --hard FETCH_HEAD)" >&2
    fail=1
fi

if ! diff -u "$AONE" "$PRR"; then
    echo "sync_provider_dedup_test: 两 skill wrapper 出现差异,违反单点约定;以任一份为准 cp 覆盖另一份" >&2
    fail=1
fi

for w in "$AONE" "$PRR"; do
    if ! grep -q "bootstrap/sync-provider.sh" "$w"; then
        echo "sync_provider_dedup_test: $w 未指向真源 bootstrap/sync-provider.sh" >&2
        fail=1
    fi
    if grep -q "reset --hard" "$w"; then
        echo "sync_provider_dedup_test: $w 含真实同步逻辑,应为薄 wrapper(逻辑只进 bootstrap/)" >&2
        fail=1
    fi
done

if [ "$fail" -eq 0 ]; then
    echo "sync_provider_dedup_test: PASS(bootstrap 真源 + 双薄 wrapper 一致)"
    exit 0
fi
exit 1
