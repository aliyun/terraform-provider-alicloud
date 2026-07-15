#!/usr/bin/env bash
# cloudspec-core.sh — 管理仓库内的 cloudspec-core 受控技能快照。
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"
lock="$repo_root/config/cloudspec-core.lock.json"

usage() {
  cat <<'EOF'
usage: bootstrap/cloudspec-core.sh <check|doctor|sync> [plugin-checkout]

  check   校验锁文件、Claude/Codex 双端技能快照及 hooks 排除策略
  doctor  在 check 基础上检查 amp、aliyun cspec 等运行依赖
  sync    从锁定 commit 的 cloudspec-plugin checkout 刷新受控技能快照
EOF
}

die() {
  echo "cloudspec-core: $*" >&2
  exit 1
}

lock_value() {
  jq -er "$1" "$lock"
}

required_skills() {
  jq -er '.skills[]' "$lock"
}

check_snapshot() {
  test -f "$lock" || die "missing lock: $lock"
  command -v jq >/dev/null 2>&1 || die "jq is required"
  test "$(lock_value '.distribution')" = "vendored-skill-snapshot" \
    || die "unsupported distribution"
  test "$(lock_value '.hooksIncluded')" = "false" \
    || die "Marketplace hooks must stay excluded"

  local skill
  local -a mirror_files=()
  while IFS= read -r skill; do
    test -f "$repo_root/.claude/skills/$skill/SKILL.md" \
      || die "missing Claude skill: $skill"
    test -f "$repo_root/.agents/skills/$skill/SKILL.md" \
      || die "missing Codex skill: $skill"
    mirror_files+=("$repo_root/.claude/skills/$skill/SKILL.md")
  done < <(required_skills)

  test ! -e "$repo_root/.claude/plugins/cloudspec-core/hooks" \
    || die "telemetry hooks must not be vendored"
  test ! -e "$repo_root/.agents/plugins/cloudspec-core/hooks" \
    || die "telemetry hooks must not be vendored"

  # 快速健康检查只验各 skill 入口；全量 references 镜像由 pre-commit/测试覆盖。
  bash "$repo_root/bootstrap/mirror.sh" check "${mirror_files[@]}" >/dev/null \
    || die "skill mirror drift detected"
  echo "cloudspec-core: snapshot OK ($(lock_value '.version') @ $(lock_value '.source.commit'))"
}

doctor() {
  check_snapshot

  local failed=0
  local command
  for command in git amp aliyun python3; do
    if command -v "$command" >/dev/null 2>&1; then
      echo "PASS $command"
    else
      echo "FAIL $command"
      failed=1
    fi
  done

  if command -v aliyun >/dev/null 2>&1 && aliyun cspec --help >/dev/null 2>&1; then
    echo "PASS aliyun-cspec"
  else
    echo "FAIL aliyun-cspec"
    failed=1
  fi

  if command -v amp >/dev/null 2>&1 && amp --version >/dev/null 2>&1; then
    echo "PASS amp-cli"
  else
    echo "FAIL amp-cli"
    failed=1
  fi

  return "$failed"
}

sync_snapshot() {
  local checkout="${1:-${CLOUDSPEC_PLUGIN_CHECKOUT:-}}"
  test -n "$checkout" || die "sync requires plugin-checkout or CLOUDSPEC_PLUGIN_CHECKOUT"
  test -d "$checkout/.git" || die "not a git checkout: $checkout"

  local expected_commit expected_tree actual_commit actual_tree source_root skill
  expected_commit="$(lock_value '.source.commit')"
  expected_tree="$(lock_value '.source.skillsTree')"
  actual_commit="$(git -C "$checkout" rev-parse HEAD)"
  actual_tree="$(git -C "$checkout" rev-parse 'HEAD:plugins/cloudspec-core/skills')"
  test "$actual_commit" = "$expected_commit" \
    || die "checkout commit $actual_commit != locked $expected_commit"
  test "$actual_tree" = "$expected_tree" \
    || die "skills tree $actual_tree != locked $expected_tree"

  source_root="$checkout/plugins/cloudspec-core/skills"
  while IFS= read -r skill; do
    test -f "$source_root/$skill/SKILL.md" || die "source skill missing: $skill"
    rm -rf "$repo_root/.claude/skills/$skill"
    cp -R "$source_root/$skill" "$repo_root/.claude/skills/$skill"
  done < <(required_skills)

  bash "$repo_root/bootstrap/mirror.sh" to-codex --all >/dev/null
  check_snapshot
}

case "${1:-}" in
  check) check_snapshot ;;
  doctor) doctor ;;
  sync) shift; sync_snapshot "${1:-}" ;;
  -h|--help|help|'') usage ;;
  *) usage >&2; exit 2 ;;
esac
