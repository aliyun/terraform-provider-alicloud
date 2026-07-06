#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

# .claude ↔ .agents 是 token-transform 镜像(见 bootstrap/skills-mirror-lib.sh):
# `.claude/skills/` ↔ `.Codex/skills/` 等按方向替换。裸 diff 会把这些**正确的**
# 镜像替换误报成 drift,因此必须先跑 sed transform 再比对(mirror-aware),
# 与 html_report_preview_skill_sync_test.sh 口径一致。真源同步性由
# bootstrap/mirror.sh check 保证;本测试锚定 provider-resource-dev 这份镜像不裂。
# shellcheck source=../bootstrap/skills-mirror-lib.sh
source "$repo_root/bootstrap/skills-mirror-lib.sh"

for rel in \
  "SKILL.md" \
  "references/zhenyuan-verification.md"; do
  expected="$tmpdir/${rel//\//__}"
  mirror_sed_codex_to_claude < "$repo_root/.agents/skills/provider-resource-dev/$rel" > "$expected"
  diff -u \
    "$expected" \
    "$repo_root/.claude/skills/provider-resource-dev/$rel"
done

echo "provider_resource_dev_skill_sync_test: PASS"
