#!/usr/bin/env bash
# PreToolUse guard for Bash: deterministic red lines that must hold even when
# the session runs with bypassPermissions (headless workers) and no auto-mode
# classifier is in the loop. Complements worktree-guard (edit-on-master) and
# run-interactive-worker-hook (cr quit family).
#
# Blocked patterns (autonomy.md 永久停止项 / redline):
#   1. git push whose refspec targets master/main — jarvis 仓与上游 master 是
#      release_prod 人工硬门 (pushing a FEATURE branch is normal and passes).
#   2. Any git push mentioning the upstream slug aliyun/terraform-provider-alicloud
#      (only the api-tool-agent fork may be pushed, via github-identity.sh).
#   3. Catastrophic rm: `rm -rf` (any flag order) aimed at /, ~, $HOME, or a
#      top-level system dir. Repo-scoped rm -rf keeps working.
#
# Contract (same as worktree-guard):
#   stdin  = tool call JSON {tool_name, tool_input:{command,...}}
#   exit 0 = allow;  exit 2 + stderr = block
# Parse failures allow — the intended failure mode is a specific block, never
# "all Bash is broken".
set -uo pipefail

# Per-turn escape hatch, same semantics as worktree-guard: repo owner only.
if [ "${JARVIS_MASTER_OK:-}" = "1" ]; then
    exit 0
fi

input="$(cat)"
cmd=$(printf '%s' "$input" | jq -r '.tool_input.command // empty' 2>/dev/null || true)
[ -n "$cmd" ] || exit 0

deny() {
    echo "redline-guard: $1" >&2
    exit 2
}

# --- git push red lines -----------------------------------------------------
# Examine each `git push` occurrence's argument tail (handles chained commands).
if printf '%s' "$cmd" | grep -qE '(^|[;&|[:space:]])git[[:space:]]+push'; then
    if printf '%s' "$cmd" | grep -qE 'aliyun/terraform-provider-alicloud'; then
        deny "push touching upstream aliyun/terraform-provider-alicloud is release_prod (human gate); use bootstrap/github-identity.sh push to the api-tool-agent fork instead"
    fi
    # refspec targeting master/main: "origin master", "master:master",
    # "HEAD:master", "+master", "origin main" …
    if printf '%s' "$cmd" | grep -qE 'git[[:space:]]+push[^;&|]*[[:space:]:+]((refs/heads/)?(master|main))([[:space:]]|$|;|\||&)'; then
        deny "push targeting master/main is release_prod (human gate, autonomy.md 永久停止项); push a feature branch and open an MR instead"
    fi
fi

# --- catastrophic rm --------------------------------------------------------
# rm with both -r and -f (any order/combined) aimed at filesystem roots.
if printf '%s' "$cmd" | grep -qE '(^|[;&|[:space:]])rm[[:space:]]+(-[a-zA-Z]*r[a-zA-Z]*f[a-zA-Z]*|-[a-zA-Z]*f[a-zA-Z]*r[a-zA-Z]*)[[:space:]]'; then
    if printf '%s' "$cmd" | grep -qE 'rm[[:space:]]+-[a-zA-Z]+[[:space:]]+("?\$HOME"?/?|~/?|/|/Users/?|/home/?|/etc/?|/usr/?|/var/?)([[:space:]]|$|;|\||&)'; then
        deny "rm -rf aimed at a filesystem root (/, ~, \$HOME, /Users, /home, /etc, /usr, /var) is blocked; target a specific project path"
    fi
fi

exit 0
