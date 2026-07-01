#!/usr/bin/env bash
# PreToolUse guard for Edit/Write/MultiEdit/NotebookEdit.
# Blocks writes when current git branch is master/main, unless:
#   1) target path is inside a repo that does NOT contain bootstrap/master-allowlist
#      (opt-in per repo — repos without the marker are unguarded)
#   2) target path matches a line in bootstrap/master-allowlist (glob, prefix ok)
#   3) env JARVIS_MASTER_OK=1 is set (per-turn escape hatch, only when the
#      repo owner has explicitly said "go on master")
#
# Contract:
#   stdin  = full Claude Code tool call JSON ({tool_name, tool_input:{file_path,...}})
#   exit 0 = allow
#   exit 2 = block with reason on stderr (Claude sees it and stops)
#
# Wired via .claude/settings.json hooks.PreToolUse[matcher=Edit|Write|MultiEdit|NotebookEdit].
set -euo pipefail

# Env-level per-turn override. Repo owner sets this before delegating direct-to-master work.
if [ "${JARVIS_MASTER_OK:-}" = "1" ]; then
    exit 0
fi

# Extract file_path from tool JSON on stdin. Fall through to allow on parse error
# so an incidental hook glitch never blocks all edits — the intended failure mode
# is a specific "on master" block, not "everything is broken".
input="$(cat)"
file_path=$(printf '%s' "$input" | jq -r '.tool_input.file_path // .tool_input.notebook_path // empty' 2>/dev/null || true)
[ -n "$file_path" ] || exit 0

# Resolve absolute path; if the target dir doesn't exist yet, resolve the parent.
# `pwd -P` normalizes /var vs /private/var on macOS so string-strip vs repo_root works.
if [ -e "$file_path" ]; then
    abs="$(cd "$(dirname "$file_path")" 2>/dev/null && pwd -P)/$(basename "$file_path")"
else
    parent="$(dirname "$file_path")"
    [ -d "$parent" ] || exit 0   # dangling path, let the tool itself decide
    abs="$(cd "$parent" && pwd -P)/$(basename "$file_path")"
fi

# Not in a git repo? out of scope. Normalize repo_root the same way (git may return /private/var).
repo_raw="$(git -C "$(dirname "$abs")" rev-parse --show-toplevel 2>/dev/null || true)"
[ -n "$repo_raw" ] || exit 0
repo_root="$(cd "$repo_raw" && pwd -P)"

# Guard is opt-in per repo: the repo must carry bootstrap/master-allowlist (even if empty).
allowlist="$repo_root/bootstrap/master-allowlist"
[ -f "$allowlist" ] || exit 0

branch="$(git -C "$repo_root" rev-parse --abbrev-ref HEAD 2>/dev/null || echo '?')"
case "$branch" in
    master|main) ;;   # continue to allowlist check
    *) exit 0 ;;      # any non-primary branch: allow
esac

# Path relative to repo root, for allowlist matching.
rel="${abs#"$repo_root"/}"

# Allowlist: one entry per line, blank / #-comment lines skipped. Supports:
#   - exact path                (e.g.  README.md)
#   - directory prefix ending / (e.g.  runs/)
#   - glob via fnmatch          (e.g.  .my-day/**  or  **/*.log)
while IFS= read -r pattern || [ -n "$pattern" ]; do
    pattern="${pattern%%#*}"
    pattern="${pattern#"${pattern%%[![:space:]]*}"}"   # ltrim
    pattern="${pattern%"${pattern##*[![:space:]]}"}"   # rtrim
    [ -z "$pattern" ] && continue

    # Directory prefix
    case "$pattern" in
        */) [ "${rel#"$pattern"}" != "$rel" ] && exit 0 ;;
    esac

    # Glob / exact — bash extglob-free, use case
    # shellcheck disable=SC2254
    case "$rel" in
        $pattern) exit 0 ;;
    esac
done < "$allowlist"

# Blocked. Reason goes to stderr; Claude Code surfaces it to the model + user.
cat >&2 <<EOF
worktree-guard: refusing Edit/Write on branch '$branch' for path:
  $rel

CLAUDE.md 工作纪律 #1: 改文件先开 worktree,严禁直接合入主干。

解锁方式(择一):
  1) 开 worktree 后重跑本次修改:
       git worktree add ../jarvis-<slug> -b feat/<slug> $branch
       cd ../jarvis-<slug> && <重跑改动> && git push -u origin feat/<slug>
     由主 Agent 通过 gh pr create / a1 CR 走人工评审合入。

  2) 仓库主人当面授权本轮直改 master(必须 "go on master" / "直接改" 明说,
     任务级"帮我改 X"不含此授权):
       JARVIS_MASTER_OK=1 <重跑本次工具调用>
     env 只在当前 shell 生效,下条 Edit/Write 仍需再次显式授权。

  3) 若该路径本应长期免管(如 runs/*、.my-day/*、bootstrap 自身补丁),
     加行到 $allowlist,提 CR 评审合入后生效。

repo=$repo_root  branch=$branch
EOF
exit 2
