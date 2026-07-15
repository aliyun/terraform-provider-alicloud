#!/bin/bash
# test/codex_agents_sync_test.sh — .codex/agents/*.toml 必须等于由 .claude/agents/*.md
# 生成的内容(md 为唯一真源;生成规则见 bootstrap/agents-mirror-lib.sh)。
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

# shellcheck source=../bootstrap/agents-mirror-lib.sh
source "$repo_root/bootstrap/agents-mirror-lib.sh"

fail=0
found=0
for md in "$repo_root"/.claude/agents/*.md; do
    [ -f "$md" ] || continue
    found=$((found + 1))
    name="$(basename "$md" .md)"
    toml="$repo_root/.codex/agents/$name.toml"
    if [ ! -f "$toml" ]; then
        echo "codex_agents_sync_test: missing $toml (expected mirror of $md)" >&2
        fail=1
        continue
    fi
    agents_md_to_toml "$md" > "$tmpdir/$name.toml"
    grep -q '^name = "' "$tmpdir/$name.toml" || { echo "codex_agents_sync_test: gen 缺 name 字段: $md" >&2; fail=1; }
    if ! diff -u "$tmpdir/$name.toml" "$toml"; then
        echo "codex_agents_sync_test: drift — $toml ≠ gen($md);修复: bash bootstrap/mirror.sh to-codex $md" >&2
        fail=1
    fi
done
if [ "$found" -lt 3 ]; then
    echo "codex_agents_sync_test: 应至少有 3 个 agent md,实得 $found" >&2
    fail=1
fi

if [ "$fail" -eq 0 ]; then
    echo "codex_agents_sync_test: PASS"
    exit 0
fi
exit 1
