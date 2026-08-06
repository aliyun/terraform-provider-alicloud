#!/usr/bin/env bash
# test/aone_assign_test.sh — policy matrix for bootstrap/aone-assign.sh.
#
# Regression for the tf_customer assignee overwrite loop (2026-07/08):
#   #84486902  辰羿 took it over → jarvis reassigned to 过载 → 辰羿 took it back
#   #84955165  jarvis 辰羿→新山 → 辰羿 took it back
#   #85115148  jarvis 过载→新山→过载 with no human involved at all
# "幂等同步 assignee" was implemented as unconditional overwrite, so a human
# takeover was reverted on every dispatch and the urgent/non-urgent classifier
# could even flip two humans back and forth by itself.
#
# Contract under test:
#   rc=0  written, or already the target (idempotent no-op)
#   rc=1  usage error / a1 write failure
#   rc=3  refused — active API-team human owns it, or ownership unverifiable
#
# A fake a1 supplies the workitem detail and records writes; no network needed.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
script="$repo_root/bootstrap/aone-assign.sh"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

pass=0; fail=0
ok(){ echo "PASS: $1"; pass=$((pass + 1)); }
no(){ echo "FAIL: $1" >&2; fail=$((fail + 1)); }

cat > "$tmp/a1" <<'EOF'
#!/usr/bin/env bash
args="$*"
if [[ "$args" == *"workitem get"* ]]; then
    [ -f "$STUB_DETAIL" ] && cat "$STUB_DETAIL"
    exit "${STUB_GET_RC:-0}"
fi
if [[ "$args" == *"workitem update"* ]]; then
    echo "$args" >> "$STUB_LOG"
    exit "${STUB_UPDATE_RC:-0}"
fi
exit 0
EOF
chmod +x "$tmp/a1"

# Real roster: the policy must agree with the shipped contacts.json, not a fixture.
contacts="$repo_root/config/contacts.json"

set_owner() {
    printf '{"id":"1","fields":[{"identifier":"assignedTo","value":"%s","displayValue":"%s"}]}' \
        "$1" "$2" > "$tmp/detail.json"
}

# run <desc> <expected-rc> <expected-write: yes|no> <workitem> <target> [env...]
run() {
    local desc="$1" want_rc="$2" want_write="$3" wi="$4" target="$5"; shift 5
    : > "$tmp/log.txt"
    local rc out
    out="$(env "$@" \
        STUB_DETAIL="$tmp/detail.json" STUB_LOG="$tmp/log.txt" \
        JARVIS_A1="$tmp/a1" JARVIS_CACHE_DIR="$tmp/cache" \
        JARVIS_CONTACTS="${OVERRIDE_CONTACTS:-$contacts}" \
        bash "$script" "$wi" "$target" 2>&1)"
    rc=$?
    local wrote="no"; [ -s "$tmp/log.txt" ] && wrote="yes"
    if [ "$rc" = "$want_rc" ] && [ "$wrote" = "$want_write" ]; then
        ok "$desc (rc=$rc write=$wrote)"
    else
        no "$desc — want rc=$want_rc write=$want_write, got rc=$rc write=$wrote :: $(echo "$out" | head -1)"
    fi
}

# ── protected: active API-team humans ───────────────────────────────────────
set_owner 320687 辰羿
run "human owner (辰羿) blocks reassign to 过载"        3 no 85020657 484483
set_owner 484483 过载
run "human owner (过载) blocks reassign to 新山"        3 no 85115148 521957
set_owner 521957 新山
run "human owner (新山) blocks reassign back to 过载"   3 no 85115148 484483

# ── idempotent convergence: the route phase re-runs every dispatch ──────────
set_owner 484483 过载
run "target == current is a no-op success"              0 no 84486902 484483

# ── reassignable: agents and people outside the roster ──────────────────────
set_owner WORKER_1782379562571 open-jarvis
run "digital worker owner is reassignable"              0 yes 84486902 484483
set_owner 416996 鲁斯
run "non-roster reporter is reassignable"               0 yes 85020657 484483
set_owner 479782 谜拟
run "legacy_inbound_only owner is reassignable"         0 yes 84486902 320687

# ── fail-closed ─────────────────────────────────────────────────────────────
set_owner 320687 辰羿
run "a1 read failure refuses (fail-closed)"             3 no 85020657 484483 STUB_GET_RC=1
printf '{"id":"1","fields":[]}' > "$tmp/detail.json"
run "missing assignedTo refuses (fail-closed)"          3 no 85020657 484483
set_owner 320687 辰羿
OVERRIDE_CONTACTS="$tmp/no-such-roster.json" \
    run "unreadable roster refuses (fail-closed)"       3 no 85020657 484483
unset OVERRIDE_CONTACTS

# ── escape hatch (repo-owner authorization only) ────────────────────────────
set_owner 320687 辰羿
run "JARVIS_ASSIGN_OK=1 bypasses protection"            0 yes 85020657 484483 JARVIS_ASSIGN_OK=1

# ── write failure surfaces, roster damage does not silently allow ───────────
set_owner WORKER_1782379562571 open-jarvis
run "a1 write failure returns rc=1"                     1 yes 84486902 484483 STUB_UPDATE_RC=1
set_owner 320687 辰羿
printf 'not json' > "$tmp/bad-roster.json"
OVERRIDE_CONTACTS="$tmp/bad-roster.json" \
    run "corrupt roster refuses (fail-closed)"          3 no 85020657 484483
unset OVERRIDE_CONTACTS

# ── argument validation ─────────────────────────────────────────────────────
rc=0; bash "$script" >/dev/null 2>&1 || rc=$?
[ "$rc" = 1 ] && ok "missing args → rc=1" || no "missing args → rc=$rc"
rc=0; bash "$script" abc 484483 >/dev/null 2>&1 || rc=$?
[ "$rc" = 1 ] && ok "non-numeric workitem id → rc=1" || no "non-numeric id → rc=$rc"
rc=0; bash "$script" 85020657 'foo;rm' >/dev/null 2>&1 || rc=$?
[ "$rc" = 1 ] && ok "unsafe staff id → rc=1" || no "unsafe staff id → rc=$rc"

echo "aone_assign_test: $pass passed, $fail failed"
[ "$fail" -eq 0 ]
