#!/usr/bin/env bash
# test/has_unmerged_mr_test.sh — unit tests for claim.sh::_has_unmerged_mr
#
# Regression for ticket 85020657 ping-pong: `gh search prs <aone-id>` returned
# [] because PR bodies are sanitized (CLAUDE.md 工作纪律 #5) and never carry
# the workitem id; finish then fired on an OPEN PR.
#
# Coverage:
#   Check 2  — concrete PR links from JARVIS_MR_CR_LINKS ∪ the ticket's own
#              Aone comments, each resolved via `gh pr view --json state`.
#   Check 3  — head-branch match against open api-tool-agent PRs, used only when
#              no concrete link was found anywhere.
#
# Fake `gh` and `a1` are used so no network/token is required.

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
JARVIS_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Extract the helpers without sourcing claim.sh — claim.sh ends with
# `*) exit 1` so a dummy command would terminate the sourcing shell.
eval "$(awk '/^_extract_pr_links\(\)/,/^}/' "$JARVIS_ROOT/bootstrap/claim.sh")"
eval "$(awk '/^_has_unmerged_mr\(\)/{f=1} f{print} f&&/^}/{exit}' \
    "$JARVIS_ROOT/bootstrap/claim.sh")"

# Point jarvis_root at a temp dir so Check 1 (git worktree list) finds nothing.
jarvis_root="$(mktemp -d)"
GH_DIR="$(mktemp -d)"
trap 'rm -rf "$jarvis_root" "$GH_DIR"' EXIT

cat > "$GH_DIR/gh" <<'EOF'
#!/usr/bin/env bash
if [ "${1:-}" = "pr" ] && [ "${2:-}" = "list" ]; then
    printf '%s\n' ${FAKE_OPEN_BRANCHES:-}
    exit 0
fi
if [ "${1:-}" = "pr" ] && [ "${2:-}" = "view" ]; then
    case "${3:-}" in
        10117|10119) echo "MERGED"; exit 0 ;;
        10118) echo "OPEN"; exit 0 ;;
        10120) exit 1 ;;   # gh failure (PR removed / no token / network)
        *) exit 1 ;;
    esac
fi
exit 1
EOF
chmod +x "$GH_DIR/gh"

# Fake a1: only `project workitem comment list <id> -f json` matters here.
cat > "$GH_DIR/fake_a1" <<'EOF'
#!/usr/bin/env bash
[ -n "${FAKE_COMMENTS:-}" ] && [ -f "$FAKE_COMMENTS" ] && cat "$FAKE_COMMENTS"
exit 0
EOF
chmod +x "$GH_DIR/fake_a1"

export PATH="$GH_DIR:$PATH"
export JARVIS_GITHUB_TOKEN="fake-token"
export FAKE_OPEN_BRANCHES=""
export FAKE_COMMENTS=""

pass=0; fail=0
check() {
    local desc="$1" expected="$2"; shift 2
    _has_unmerged_mr "$@"
    local rc=$?
    if [ "$rc" = "$expected" ]; then
        echo "  PASS: $desc (rc=$rc)"; pass=$((pass+1))
    else
        echo "  FAIL: $desc — expected rc=$expected got rc=$rc"; fail=$((fail+1))
    fi
}

PR=https://github.com/aliyun/terraform-provider-alicloud/pull

# --- Check 2: links from JARVIS_MR_CR_LINKS ---------------------------------
JARVIS_MR_CR_LINKS="$PR/10118"
check "OPEN PR blocks finish" 0 999test999

JARVIS_MR_CR_LINKS="$PR/10117"
check "MERGED PR does not block" 1 999test999

JARVIS_MR_CR_LINKS="$PR/10120"
check "gh failure blocks finish (unverified)" 0 999test999

JARVIS_MR_CR_LINKS="https://code.example/cr/999"
check "non-github link skipped, no block" 1 999test999

JARVIS_MR_CR_LINKS="$PR/10117 $PR/10118"
check "MERGED+OPEN mix blocks finish" 0 999test999

JARVIS_MR_CR_LINKS="$PR/10119/files"
check "/files suffix num parsed, MERGED no block" 1 999test999

# --- Check 2: links harvested from the ticket's own Aone comments -----------
# The RD aggregate reply carries the PR link in markdown form, so the ticket is
# a durable record of its own PRs even when the caller passes no env links.
A1="$GH_DIR/fake_a1"
JARVIS_MR_CR_LINKS=""
cat > "$GH_DIR/comments-open.json" <<EOF
[{"content":"验收通过。链接：PR [$PR/10118]($PR/10118)"}]
EOF
FAKE_COMMENTS="$GH_DIR/comments-open.json"
check "OPEN PR harvested from Aone comments blocks finish" 0 999test999

cat > "$GH_DIR/comments-merged.json" <<EOF
[{"content":"已合并。链接：PR [$PR/10117]($PR/10117)"}]
EOF
FAKE_COMMENTS="$GH_DIR/comments-merged.json"
check "MERGED PR harvested from comments does not block" 1 999test999

# Under-reported env links: executor passed only the merged one, the ticket also
# records an open one. The union must still block.
JARVIS_MR_CR_LINKS="$PR/10117"
FAKE_COMMENTS="$GH_DIR/comments-open.json"
check "under-reported env link + open PR in comments blocks" 0 999test999

# Same link from both sources must be deduped, not double-queried.
JARVIS_MR_CR_LINKS="$PR/10117"
FAKE_COMMENTS="$GH_DIR/comments-merged.json"
check "duplicate link across sources still merges clean" 1 999test999

A1=""
FAKE_COMMENTS=""
JARVIS_MR_CR_LINKS=""

# --- Check 3: head-branch fallback when no concrete link exists -------------
FAKE_OPEN_BRANCHES="worktree-999test999-gpdb-status feat/other-thing"
check "open PR whose branch carries the id blocks finish" 0 999test999

FAKE_OPEN_BRANCHES="feat/unrelated docs/vpc-router-id-description"
check "no branch match and no links → finish allowed" 1 999test999

# A concrete link outranks the branch heuristic: link says MERGED, so even a
# stale same-id branch on an unrelated open PR must not block.
JARVIS_MR_CR_LINKS="$PR/10117"
FAKE_OPEN_BRANCHES="worktree-999test999-stale"
check "concrete MERGED link wins over stale branch match" 1 999test999
JARVIS_MR_CR_LINKS=""
FAKE_OPEN_BRANCHES=""

echo "PASS=$pass FAIL=$fail"
[ "$fail" = 0 ]
