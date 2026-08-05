#!/usr/bin/env bash
# PreToolUse guard for Bash: deterministic red lines that must hold even when
# the session runs with bypassPermissions (headless workers) and no auto-mode
# classifier is in the loop. Complements worktree-guard (edit-on-master) and
# run-interactive-worker-hook (cr quit family).
#
# Blocked patterns (autonomy.md 永久停止项 / redline):
#   1. git push whose refspec targets master/main — jarvis 仓与上游 master 是
#      release_prod 人工硬门 (pushing a FEATURE branch is normal and passes).
#      Exemption: repos whose origin slug is listed in bootstrap/push-master-allowlist
#      (data repos with a direct-push-to-master contract, e.g. terraflow/tf_playground
#      的 KNOWLEDGE 蒸馏仓). The repo dir must be resolvable from `git -C <dir>` or a
#      leading `cd <dir> &&`; anything unresolvable stays blocked (fail-closed).
#   2. Any git push mentioning the upstream slug aliyun/terraform-provider-alicloud
#      (only the api-tool-agent fork may be pushed, via github-identity.sh).
#   3. Catastrophic rm: `rm -rf` (any flag order) aimed at /, ~, $HOME, or a
#      top-level system dir. Repo-scoped rm -rf keeps working.
#   4. Any actual curl/wget invocation of Acube createBuildTaskV2. Terraform
#      work continues on the source Aone Task; this API must never create a
#      downstream 528766 workitem again.
#   5. Jarvis delivery execution is allowlist-only: config/pools.json daily and
#      prestage IDs may run; every production, new, unknown, opaque rerun, and
#      opaque task-action path stops for a human. Pipeline renumbering therefore
#      fails closed instead of bypassing a fixed production-ID denylist.
#
# Contract (same as worktree-guard):
#   stdin  = tool call JSON {tool_name, tool_input:{command,...}}
#   exit 0 = allow;  exit 2 + stderr = block
# Parse failures allow — the intended failure mode is a specific block, never
# "all Bash is broken".
set -uo pipefail

input="$(cat)"
cmd=$(printf '%s' "$input" | jq -r '.tool_input.command // empty' 2>/dev/null || true)
[ -n "$cmd" ] || exit 0

deny() {
    echo "redline-guard: $1" >&2
    exit 2
}

# --- Acube downstream workitem red line ------------------------------------
# Delegate shell tokenization to the same parser used by the Codex/interactive
# Worker fence. It distinguishes executed curl/wget (including bash -lc and
# same-command variable wrappers) from rg/grep/printf text inspection.
guard_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
guard_py="${JARVIS_A1_COMMAND_GUARD:-$guard_dir/a1_command_guard.py}"
python_bin="/usr/bin/python3"
[ -x "$python_bin" ] || python_bin="$(command -v python3 2>/dev/null || true)"
if [ -n "$python_bin" ] && [ -f "$guard_py" ]; then
    "$python_bin" -I "$guard_py" --check-pretool-command "$cmd" \
        >/dev/null 2>"${TMPDIR:-/tmp}/jarvis-redline-guard.$$"
    guard_rc=$?
    if [ "$guard_rc" -eq 2 ]; then
        detail="$(cat "${TMPDIR:-/tmp}/jarvis-redline-guard.$$" 2>/dev/null || true)"
        rm -f "${TMPDIR:-/tmp}/jarvis-redline-guard.$$"
        deny "${detail#a1 safety: }"
    fi
    rm -f "${TMPDIR:-/tmp}/jarvis-redline-guard.$$"
    if [ "$guard_rc" -ne 0 ]; then
        if printf '%s' "$cmd" | grep -qE '(^|[;&|[:space:]])([^;&|[:space:]]*/)?(a1|a1id)([[:space:]]|$)'; then
            deny "a1 safety classifier failed(rc=$guard_rc); blocking a1/a1id command fail-closed"
        fi
    fi
else
    # The specific red line must fail closed if its parser is unavailable.
    if printf '%s' "$cmd" | grep -Fq 'createBuildTaskV2' \
        && printf '%s' "$cmd" | grep -qE '(^|[;&|[:space:]])(curl|wget)([[:space:]]|$)'; then
        deny "Acube createBuildTaskV2 is permanently disabled for Jarvis"
    fi
    if printf '%s' "$cmd" | grep -qE '(^|[;&|[:space:]])([^;&|[:space:]]*/)?(a1|a1id)([[:space:]]|$)'; then
        deny "a1 safety classifier is unavailable; blocking a1/a1id command fail-closed"
    fi
fi

# Per-turn escape hatch applies to repository/master maintenance red lines,
# never to the permanent Acube downstream-task prohibition above.
if [ "${JARVIS_MASTER_OK:-}" = "1" ]; then
    exit 0
fi

# Master-push exemption for contract-approved data repos (fail-closed):
# resolve the pushed repo's dir from the command text, read the pushed remote's
# URL, and allow only when its slug is listed in bootstrap/push-master-allowlist
# with a path boundary (":slug[.git]" / "/slug[.git]" — no substring over-match).
push_master_allowlisted() {
    local allowlist
    allowlist="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/push-master-allowlist"
    [ -f "$allowlist" ] || return 1

    # Repo dir: `git -C <dir> push …` (quoted or bare) or a leading `cd <dir> && git push …`.
    local dir=""
    dir=$(printf '%s' "$cmd" | sed -nE 's/.*git[[:space:]]+-C[[:space:]]+"([^"]+)"[[:space:]]+push.*/\1/p' | head -1)
    [ -n "$dir" ] || dir=$(printf '%s' "$cmd" | sed -nE 's/.*git[[:space:]]+-C[[:space:]]+([^"[:space:]]+)[[:space:]]+push.*/\1/p' | head -1)
    [ -n "$dir" ] || dir=$(printf '%s' "$cmd" | sed -nE 's/^[[:space:]]*cd[[:space:]]+"([^"]+)"[[:space:]]*(&&|;)[[:space:]]*git[[:space:]]+push.*/\1/p' | head -1)
    [ -n "$dir" ] || dir=$(printf '%s' "$cmd" | sed -nE 's/^[[:space:]]*cd[[:space:]]+([^"[:space:]&;|]+)[[:space:]]*(&&|;)[[:space:]]*git[[:space:]]+push.*/\1/p' | head -1)
    [ -n "$dir" ] && [ -d "$dir" ] || return 1

    # Remote = first non-option token after `push`; default origin.
    local tail remote="" tok
    tail=$(printf '%s' "$cmd" | sed -nE 's/.*git[[:space:]]+(-C[[:space:]]+("[^"]*"|[^[:space:]]+)[[:space:]]+)?push[[:space:]]+(.*)$/\3/p' | head -1)
    for tok in $tail; do
        case "$tok" in
            -*) continue ;;
            *) remote="$tok"; break ;;
        esac
    done
    [ -n "$remote" ] || remote="origin"

    local url slug
    url=$(git -C "$dir" remote get-url "$remote" 2>/dev/null) || return 1
    while IFS= read -r slug; do
        case "$slug" in ''|\#*) continue ;; esac
        case "$url" in
            *":$slug.git"|*"/$slug.git"|*":$slug"|*"/$slug") return 0 ;;
        esac
    done < "$allowlist"
    return 1
}

# --- git push red lines -----------------------------------------------------
# Examine each `git push` occurrence's argument tail (handles chained commands).
# `git -C <dir> push` must match too — a bare `git[[:space:]]+push` pattern lets
# any -C form bypass the red line entirely.
git_push_re='git([[:space:]]+-C[[:space:]]+("[^"]*"|[^[:space:]]+))?[[:space:]]+push'
if printf '%s' "$cmd" | grep -qE "(^|[;&|[:space:]])$git_push_re"; then
    if printf '%s' "$cmd" | grep -qE 'aliyun/terraform-provider-alicloud'; then
        deny "push touching upstream aliyun/terraform-provider-alicloud is release_prod (human gate); use bootstrap/github-identity.sh push to the api-tool-agent fork instead"
    fi
    # refspec targeting master/main: "origin master", "master:master",
    # "HEAD:master", "+master", "origin main" …
    if printf '%s' "$cmd" | grep -qE "$git_push_re"'[^;&|]*[[:space:]:+]((refs/heads/)?(master|main))([[:space:]]|$|;|\||&)'; then
        push_master_allowlisted || deny "push targeting master/main is release_prod (human gate, autonomy.md 永久停止项); push a feature branch and open an MR instead (data repos with a direct-push contract go in bootstrap/push-master-allowlist)"
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
