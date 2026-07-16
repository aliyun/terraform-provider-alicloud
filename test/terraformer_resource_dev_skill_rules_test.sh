#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

# shellcheck source=../bootstrap/skills-mirror-lib.sh
source "$repo_root/bootstrap/skills-mirror-lib.sh"

for rel in \
  "SKILL.md" \
  "references/alicloud-resource-development.md"; do
  claude_file="$repo_root/.claude/skills/terraformer-resource-dev/$rel"
  codex_file="$repo_root/.agents/skills/terraformer-resource-dev/$rel"
  test -f "$claude_file"
  test -f "$codex_file"

  expected="$tmpdir/${rel//\//__}"
  mirror_sed_codex_to_claude < "$codex_file" > "$expected"
  diff -u "$expected" "$claude_file"
done

expected_layout=$'SKILL.md\nreferences/alicloud-resource-development.md'
for skill_root in \
  "$repo_root/.claude/skills/terraformer-resource-dev" \
  "$repo_root/.agents/skills/terraformer-resource-dev"; do
  actual_layout="$(find "$skill_root" -type f | sed "s#^$skill_root/##" | LC_ALL=C sort)"
  if [[ "$actual_layout" != "$expected_layout" ]]; then
    echo "terraformer_resource_dev_skill_rules_test: unexpected layout in $skill_root" >&2
    diff -u <(printf '%s\n' "$expected_layout") <(printf '%s\n' "$actual_layout") >&2 || true
    exit 1
  fi
done

for skill in \
  "$repo_root/.claude/skills/terraformer-resource-dev/SKILL.md" \
  "$repo_root/.agents/skills/terraformer-resource-dev/SKILL.md"; do
  for term in \
    "description: 用于开发、诊断或修复 Terraformer 中的阿里云资源" \
    "# Terraformer 资源开发" \
    "## 核心模型" \
    "## 每次任务的起始动作" \
    "## 证据优先级" \
    "## 选择一种资源发现模式" \
    "## 只修改适用文件" \
    "## 验证门禁" \
    "## 交付" \
    "bootstrap/workspace.sh dir terraformer" \
    "停止并按 missing_capability 升级" \
    "aone-triage" \
    "loops/adhoc-intake.md" \
    "bootstrap/claim.sh claim" \
    "bootstrap/wrap.sh sync <id>" \
    "bootstrap/wrap.sh done" \
    "bootstrap/claim.sh release" \
    "references/alicloud-resource-development.md" \
    "terraform-rd" \
    "terraform-qa" \
    "InitResources" \
    "禁止生产或推导资源关联关系"; do
    if ! grep -Fq -- "$term" "$skill"; then
      echo "terraformer_resource_dev_skill_rules_test: missing '$term' in $skill" >&2
      exit 1
    fi
  done
  for old_term in \
    "description: Use when developing" \
    "# Terraformer Resource Development" \
    "## Core model" \
    "## Start every task" \
    "## Evidence hierarchy" \
    "## Evidence order" \
    "## Select one discovery pattern" \
    "## Choose one discovery pattern" \
    "## Change only applicable files" \
    "## Change only the applicable files" \
    "## Validation gates" \
    "## Delivery" \
    "prior state" \
    "fallback" \
    "schema flatten" \
    "state/HCL" \
    "Terraformer checkout" \
    "tracked files" \
    "scope/filter" \
    "child List API" \
    "client/service" \
    "endpoint"; do
    if grep -Fq -- "$old_term" "$skill"; then
      echo "terraformer_resource_dev_skill_rules_test: unexpected English prose '$old_term' in $skill" >&2
      exit 1
    fi
  done
done

for reference in \
  "$repo_root/.claude/skills/terraformer-resource-dev/references/alicloud-resource-development.md" \
  "$repo_root/.agents/skills/terraformer-resource-dev/references/alicloud-resource-development.md"; do
  for term in \
    "# Alicloud Terraformer 资源开发" \
    "## 目录" \
    "## 1. 运行时架构" \
    "## 2. 证据真源检查清单" \
    "## 3. InitResources 资源发现模式" \
    "A. 直接全量 List + 单字段 Import ID" \
    "B. 单次 List 返回多段 Import ID 的全部片段" \
    "C. 父子遍历" \
    "D. 无法完整枚举" \
    "## 4. 多段式 Import ID" \
    "## 5. 分页与错误处理" \
    "## 6. 文件选择" \
    "## 7. 测试与验证" \
    "## 8. 常见错误" \
    'd.SetId(...)' \
    'ParseResourceId(...)' \
    "多段式 Import ID 本身并不意味着必须遍历父资源" \
    "Data Source 可以要求父资源 ID" \
    "每个父资源都必须重置分页状态" \
    "使用 token 分页时，只要返回的 next token 为空就终止，不受当前页数量影响" \
    "使用页码分页时" \
    "禁止从 Provider schema、Data Source 参数或 API 字段名生产或推导关联关系" \
    "不阻塞核心的资源发现与 Import ID 支持" \
    "go test ./providers/alicloud" \
    "go test ./..." \
    "/tmp/terraformer"; do
    if ! grep -Fq -- "$term" "$reference"; then
      echo "terraformer_resource_dev_skill_rules_test: missing '$term' in $reference" >&2
      exit 1
    fi
  done
  for old_term in \
    "# Alicloud Terraformer resource development" \
    "## Contents" \
    "## 1. Runtime architecture" \
    "## 2. Source-of-truth checklist" \
    "## 3. InitResources discovery patterns" \
    "### A. Direct full List with a single-field Import ID" \
    "### B. One List returns every multipart-ID segment" \
    "### C. Parent-child traversal" \
    "### D. Complete enumeration is unavailable" \
    "## 4. Multipart Import IDs" \
    "## 5. Pagination and errors" \
    "## 6. File selection" \
    "## 7. Tests and validation" \
    "## 8. Common mistakes" \
    "prior state" \
    "fallback" \
    "schema flatten" \
    "import round trip" \
    "child List API" \
    "attachment" \
    "page size" \
    "page number" \
    "action" \
    "endpoint" \
    "decode" \
    "retry helper" \
    "consumer" \
    "producer" \
    "connection map" \
    "drift"; do
    if grep -Fq -- "$old_term" "$reference"; then
      echo "terraformer_resource_dev_skill_rules_test: unexpected English prose '$old_term' in $reference" >&2
      exit 1
    fi
  done
done

evaluation_report="$repo_root/docs/superpowers/reports/2026-07-16-terraformer-resource-dev-forward-evaluation.md"
test -f "$evaluation_report"
for term in \
  "Scenario A — PASS" \
  "Scenario B — PASS" \
  "Scenario C — PASS"; do
  if ! grep -Fq -- "$term" "$evaluation_report"; then
    echo "terraformer_resource_dev_skill_rules_test: missing '$term' in $evaluation_report" >&2
    exit 1
  fi
done

echo "terraformer_resource_dev_skill_rules_test: PASS"
