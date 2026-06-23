#!/usr/bin/env bash
# test/workspace_test.sh – registry lookups for bootstrap/workspace.sh
# Uses a tmp JARVIS_WORKSPACES_FILE. Covers resolve, ops, ensure(missing), list.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

PASS=0
FAIL=0
assert_pass() { echo "PASS: $1"; PASS=$((PASS + 1)); }
assert_fail() { echo "FAIL: $1"; FAIL=$((FAIL + 1)); }

# tmp config: one repo with a real path, one mcp_server with a pool, no ops
mkdir -p "$tmp/exists"
cat > "$tmp/workspaces.json" << JSON
{
  "workspaces": {
    "terraform_provider": {
      "repo": "terraform-provider-alicloud", "path": "$tmp/exists",
      "pools": ["tf_provider","tf_customer"],
      "ops": {"build":"go build ./...","test":"go test ./alicloud -run <Name>","vet":"go vet ./...","fmt":"gofmt -l alicloud/"}
    },
    "mcp_server": { "repo":"aliyun-automation-agent","pool":"mcp_server","app":283346,"ops":{} }
  }
}
JSON
export JARVIS_WORKSPACES_FILE="$tmp/workspaces.json"
WS="$proj_root/bootstrap/workspace.sh"

# resolve by repo, key, pool
[ "$(bash "$WS" resolve terraform-provider-alicloud)" = "$tmp/exists" ] \
  && assert_pass "resolve by repo" || assert_fail "resolve by repo"
[ "$(bash "$WS" resolve tf_customer)" = "$tmp/exists" ] \
  && assert_pass "resolve by pool" || assert_fail "resolve by pool"

# ops
[ "$(bash "$WS" ops terraform-provider-alicloud build)" = "go build ./..." ] \
  && assert_pass "ops build" || assert_fail "ops build"
[ "$(bash "$WS" ops terraform-provider-alicloud vet)" = "go vet ./..." ] \
  && assert_pass "ops vet" || assert_fail "ops vet"

# ensure: present path → 0, missing path → 3
bash "$WS" ensure terraform-provider-alicloud; [ $? -eq 0 ] \
  && assert_pass "ensure present exit 0" || assert_fail "ensure present"
bash "$WS" ensure aliyun-automation-agent >/dev/null 2>&1; [ $? -eq 3 ] \
  && assert_pass "ensure missing exit 3" || assert_fail "ensure missing exit 3"

# list: both repos appear
out=$(bash "$WS" list)
echo "$out" | grep -q "terraform-provider-alicloud" && echo "$out" | grep -q "aliyun-automation-agent" \
  && assert_pass "list both repos" || assert_fail "list both repos"

echo ""; echo "PASS: $PASS  FAIL: $FAIL"
[ "$FAIL" -gt 0 ] && { echo "TESTS FAILED"; exit 1; }
echo "All tests passed"; exit 0
