#!/usr/bin/env bash
# test/agents_md_tracked_test.sh — AGENTS.md 已入库跟踪 + pre-commit 自动同步不变量。
#
# 背景:AGENTS.md 原为 hook 生成的 untracked 文件(P6.a),漂移风险靠「手工别改」维持。
# 现改为**入库跟踪 + pre-commit 钩子自动重生成同步**,根治漂移。本测试锁两件事:
#   A. 静态:AGENTS.md 已 tracked、非 gitignored、内容 == mirror_sed_claude_to_codex(CLAUDE.md)、
#      且 round-trip 恒等(保证 mirror.sh check 双向配对能过)。
#   B. 动态:pre-commit 钩子在 staged 含 CLAUDE.md 时自动重生成 AGENTS.md 并 add 进本次提交。

set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
# shellcheck source=../bootstrap/skills-mirror-lib.sh
source "$repo_root/bootstrap/skills-mirror-lib.sh"

fail() { echo "agents_md_tracked_test: FAIL — $1" >&2; exit 1; }

# ── A. 静态不变量 ──────────────────────────────────────────────
[ -f "$repo_root/AGENTS.md" ] || fail "AGENTS.md 不在 disk 上"

git -C "$repo_root" ls-files --error-unmatch AGENTS.md >/dev/null 2>&1 \
  || fail "AGENTS.md 未被 git 跟踪(应已入库)"

if git -C "$repo_root" check-ignore -q AGENTS.md; then
  fail "AGENTS.md 仍被 .gitignore 忽略(应已移除该条目)"
fi

diff -u <(mirror_sed_claude_to_codex < "$repo_root/CLAUDE.md") "$repo_root/AGENTS.md" \
  || fail "AGENTS.md != mirror_sed_claude_to_codex(CLAUDE.md)"

# round-trip 恒等:保证 mirror.sh check 的 AGENTS.md→CLAUDE.md 反向配对也过
diff -u <(mirror_sed_codex_to_claude < "$repo_root/AGENTS.md") "$repo_root/CLAUDE.md" \
  || fail "round-trip 破坏:sed_codex_to_claude(AGENTS.md) != CLAUDE.md"

# ── B. pre-commit 钩子端到端:改 CLAUDE.md 只 stage 它,AGENTS.md 应被自动重生成+入提交 ──
sandbox="$(mktemp -d)"
trap 'rm -rf "$sandbox"' EXIT
mkdir -p "$sandbox/bootstrap/git-hooks"
cp "$repo_root/bootstrap/mirror.sh"            "$sandbox/bootstrap/mirror.sh"
cp "$repo_root/bootstrap/skills-mirror-lib.sh" "$sandbox/bootstrap/skills-mirror-lib.sh"
cp "$repo_root/bootstrap/git-hooks/pre-commit" "$sandbox/bootstrap/git-hooks/pre-commit"
chmod +x "$sandbox/bootstrap/git-hooks/pre-commit"
cp "$repo_root/CLAUDE.md"                      "$sandbox/CLAUDE.md"

git -C "$sandbox" init -q
git -C "$sandbox" config user.email test@example.com
git -C "$sandbox" config user.name test
git -C "$sandbox" config commit.gpgsign false
git -C "$sandbox" config core.hooksPath "$sandbox/bootstrap/git-hooks"

# 初始提交(CLAUDE.md staged → 钩子生成 AGENTS.md 并 add)
git -C "$sandbox" add bootstrap CLAUDE.md
git -C "$sandbox" commit -q -m "init"
git -C "$sandbox" ls-files --error-unmatch AGENTS.md >/dev/null 2>&1 \
  || fail "钩子未在初始提交中自动加入 AGENTS.md"

# 更新场景:只改 + stage CLAUDE.md,不碰 AGENTS.md;钩子须重生成并入提交
printf '\n<!-- sentinel-%s -->\n' "$$" >> "$sandbox/CLAUDE.md"
git -C "$sandbox" add CLAUDE.md   # 故意只 add CLAUDE.md
git -C "$sandbox" commit -q -m "edit CLAUDE.md only"

committed_agents="$(git -C "$sandbox" show HEAD:AGENTS.md)"
expected_agents="$(mirror_sed_claude_to_codex < "$sandbox/CLAUDE.md")"
[ "$committed_agents" = "$expected_agents" ] \
  || fail "更新提交里 AGENTS.md 未自动跟随 CLAUDE.md 重生成(钩子 add 逻辑失效)"

# 反证:sentinel 确实进了提交里的 AGENTS.md(证明是新内容,不是旧快照)
git -C "$sandbox" show HEAD:AGENTS.md | grep -q "sentinel-$$" \
  || fail "提交的 AGENTS.md 未含本轮 CLAUDE.md 的新增内容"

echo "agents_md_tracked_test: PASS"
