#!/usr/bin/env bash
# test/has_unmerged_mr_test.sh — unit tests for claim.sh::_has_unmerged_mr
# Check 2a (JARVIS_MR_CR_LINKS → gh pr view --json state direct query).
#
# Regression for ticket 85020657 ping-pong: `gh search prs <aone-id>` returned
# [] because PR bodies are sanitized (CLAUDE.md 工作纪律 #5) and never carry
# the workitem id; finish then fired on an OPEN PR. Check 2a parses each link's
# owner/repo/num and asks `gh pr view` directly.
#
# A fake `gh` is placed on PATH so no network/token is required.

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
JARVIS_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Extract _has_unmerged_mr() without sourcing claim.sh — claim.sh ends with
# `*) exit 1` so a dummy command would terminate the sourcing shell.
eval "$(awk '/^_has_unmerged_mr\(\)/{f=1} f{print} f&&/^}/{exit}' \
    "$JARVIS_ROOT/bootstrap/claim.sh")"

# Point jarvis_root at a temp dir so Check 1 (git worktree list) finds nothing
# and we fall through to Check 2a.
jarvis_root="$(mktemp -d)"
GH_DIR="$(mktemp -d)"
trap 'rm -rf "$jarvis_root" "$GH_DIR"' EXIT

cat > "$GH_DIR/gh" <<'EOF'
#!/usr/bin/env bash
num=""
while [ $# -gt 0 ]; do
    case "$1" in
        --repo|-q) shift ;;
        pr|view) shift ;;
        *) [ -z "${num:-}" ] && num="$1"; shift ;;
    esac
done
case "$num" in
    10117|10119) echo "MERGED"; exit 0 ;;
    10118) echo "OPEN"; exit 0 ;;
    10120) exit 1 ;;   # gh failure (PR removed / no token / network)
    *) exit 1 ;;
esac
EOF
chmod +x "$GH_DIR/gh"
export PATH="$GH_DIR:$PATH"
export JARVIS_GITHUB_TOKEN="fake-token"

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

# 1. OPEN PR → block finish (return 0)
JARVIS_MR_CR_LINKS="https://github.com/aliyun/terraform-provider-alicloud/pull/10118"
check "OPEN PR blocks finish" 0 999test999

# 2. MERGED PR → no block (return 1)
JARVIS_MR_CR_LINKS="https://github.com/aliyun/terraform-provider-alicloud/pull/10117"
check "MERGED PR does not block" 1 999test999

# 3. gh failure → conservatively block (unverified, return 0)
JARVIS_MR_CR_LINKS="https://github.com/aliyun/terraform-provider-alicloud/pull/10120"
check "gh failure blocks finish (unverified)" 0 999test999

# 4. non-github link → skipped, no block
JARVIS_MR_CR_LINKS="https://code.example/cr/999"
check "non-github link skipped, no block" 1 999test999

# 5. MERGED + OPEN mix → block (any non-MERGED blocks)
JARVIS_MR_CR_LINKS="https://github.com/aliyun/terraform-provider-alicloud/pull/10117 https://github.com/aliyun/terraform-provider-alicloud/pull/10118"
check "MERGED+OPEN mix blocks finish" 0 999test999

# 6. link with /files suffix → num still parsed (10119 MERGED → no block)
JARVIS_MR_CR_LINKS="https://github.com/aliyun/terraform-provider-alicloud/pull/10119/files"
check "/files suffix num parsed, MERGED no block" 1 999test999

echo "PASS=$pass FAIL=$fail"
[ "$fail" = 0 ]
