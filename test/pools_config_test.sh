#!/usr/bin/env bash
set -uo pipefail

POOLS_JSON="$(dirname "$0")/../config/pools.json"
PASS=0
FAIL=0

ok()  { echo "PASS: $1"; PASS=$((PASS + 1)); }
bad() { echo "FAIL: $1"; FAIL=$((FAIL + 1)); }

# S1: parses
jq . "$POOLS_JSON" >/dev/null 2>&1 && ok "jq parses config/pools.json" || bad "jq parses config/pools.json"

# S2: 5 pools across 2 lines
[ "$(jq '.pools|length' "$POOLS_JSON")" = "5" ] && ok "5 pools" || bad "expected 5 pools"
[ "$(jq '.lines|length' "$POOLS_JSON")" = "2" ] && ok "2 lines" || bad "expected 2 lines"

# S3: all 5 project ids present (verbatim)
for id in 1086837 528766 2124589 1120451 2100304; do
  jq -e --argjson p "$id" '[.pools[].project]|index($p)' "$POOLS_JSON" >/dev/null \
    && ok "project $id present" || bad "project $id missing"
done

# S4: every pool names its line
[ "$(jq '[.pools[]|select(.line==null)]|length' "$POOLS_JSON")" = "0" ] \
  && ok "every pool has a line" || bad "a pool lacks line"

# S5: mcp_server keeps delivery facts (schema evolved single .app → .apps[]; assert on apps[])
[ "$(jq -r '.pools.mcp_server.apps[0].app' "$POOLS_JSON")" = "283346" ] \
  && ok "mcp_server.apps[0].app preserved" || bad "mcp_server.apps[0].app lost"

# S6: claim.tag
[ "$(jq -r '.claim.tag' "$POOLS_JSON")" = "jarvis-claimed" ] \
  && ok "claim.tag jarvis-claimed" || bad "claim.tag wrong"

# S7: claim.done_tag
[ "$(jq -r '.claim.done_tag' "$POOLS_JSON")" = "jarvis-done" ] \
  && ok "claim.done_tag jarvis-done" || bad "claim.done_tag wrong"

# S8: claim.ttl_min
[ "$(jq '.claim.ttl_min' "$POOLS_JSON")" = "45" ] \
  && ok "claim.ttl_min 45" || bad "claim.ttl_min wrong"

# S9: api_toolkit.done_status 必须是项目 2100304 支持的枚举(为「已发布」)。
# 该项目状态枚举无「已完成」——配错会导致 claim.sh finish 时 status 更新被拒、不流转。
[ "$(jq -r '.pools.api_toolkit.done_status' "$POOLS_JSON")" = "已发布" ] \
  && ok "api_toolkit.done_status 已发布(项目有效枚举)" || bad "api_toolkit.done_status 应为「已发布」(项目 2100304 无「已完成」枚举)"

echo ""; echo "Results: $PASS passed, $FAIL failed"
[ "$FAIL" -gt 0 ] && exit 1; exit 0
